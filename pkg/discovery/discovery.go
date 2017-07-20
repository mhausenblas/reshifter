package discovery

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"

	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
)

// ProbeEtcd probes an endpoint at path /version to figure
// which version of etcd it is and in which mode (secure or insecure)
// it is used.
func ProbeEtcd(endpoint string) (version, apiversion string, secure bool, err error) {
	u, err := url.Parse(endpoint + "/version")
	if err != nil {
		return "", "", false, fmt.Errorf("Can't parse endpoint %s: %s", endpoint, err)
	}
	switch u.Scheme {
	case "https": // secure etcd
		secure = true
		clientcert, clientkey, cerr := util.ClientCertAndKeyFromEnv()
		if cerr != nil {
			err = cerr
			return "", "", false, err
		}
		version, err = getVersionSecure(u.String(), clientcert, clientkey)
		if err != nil {
			return "", "", false, err
		}
	case "http":
		secure = false
		version, err = getVersion(u.String())
		if err != nil {
			return "", "", false, err
		}
	default:
		return "", "", false, fmt.Errorf("Can't determine what scheme is use")
	}
	// try to figure out if v2 API is in use:
	apiversion = "v3"
	kprefix, err := checkv2(endpoint, secure)
	if err != nil {
		return "", "", false, err
	}
	if kprefix != "" { // a v2 API in an etcd3
		apiversion = "v2"
	}
	return version, apiversion, secure, nil
}

// ProbeKubernetesDistro probes an etcd cluster for which Kubernetes
// distribution is present by scanning the available keys.
func ProbeKubernetesDistro(endpoint string) (distro types.KubernetesDistro, err error) {
	distro = types.NotADistro
	version, apiversion, secure, err := ProbeEtcd(endpoint)
	if err != nil {
		return distro, err
	}
	switch {
	case strings.HasPrefix(version, "3"):
		switch apiversion {
		case "v2":
			distro, err = getKubeDistrov2(endpoint, secure)
			if err != nil {
				return distro, err
			}
		case "v3":
			distro, err = getKubeDistrov3(endpoint, secure)
			if err != nil {
				return distro, err
			}
		}
	case strings.HasPrefix(version, "2"):
		distro, err = getKubeDistrov2(endpoint, secure)
		if err != nil {
			return distro, err
		}
	}
	return distro, nil
}

// CountKeysFor iterates over well-known keys of a given Kubernetes distro
// and returns the number of keys and their values total size, in the
// respective key range/subtree.
func CountKeysFor(endpoint string, startkey, endkey string) (numkeys, totalsize int, err error) {
	version, apiversion, secure, err := ProbeEtcd(endpoint)
	if err != nil {
		return 0, 0, err
	}
	switch {
	case strings.HasPrefix(version, "3"):
		switch apiversion {
		case "v2":
			numkeys, totalsize, err = statsv2(endpoint, secure, startkey)
			if err != nil {
				return 0, 0, err
			}
		case "v3":
			numkeys, totalsize, err = statsv3(endpoint, secure, startkey, endkey)
			if err != nil {
				return 0, 0, err
			}
		}
	case strings.HasPrefix(version, "2"):
		numkeys, totalsize, err = statsv2(endpoint, secure, startkey)
		if err != nil {
			return 0, 0, err
		}
	}
	return numkeys, totalsize, nil
}

func getVersion(endpoint string) (etcdversion string, err error) {
	var etcdr types.EtcdResponse
	res, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	err = json.NewDecoder(res.Body).Decode(&etcdr)
	if err != nil {
		return "", fmt.Errorf("Can't decode response from etcd: %s", err)
	}
	_ = res.Body.Close()
	etcdversion = etcdr.EtcdServerVersion
	return etcdversion, nil
}

func getVersionSecure(endpoint, clientcert, clientkey string) (string, error) {
	var etcdr types.EtcdResponse
	cert, err := tls.LoadX509KeyPair(clientcert, clientkey)
	if err != nil {
		return "", err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	err = json.NewDecoder(res.Body).Decode(&etcdr)
	if err != nil {
		return "", fmt.Errorf("Can't decode response from etcd: %s", err)
	}
	_ = res.Body.Close()
	return etcdr.EtcdServerVersion, nil
}

func checkv2(endpoint string, secure bool) (string, error) {
	c2, err := util.NewClient2(endpoint, secure)
	if err != nil {
		return "", fmt.Errorf("Can't connect to etcd v2 API: %s", err)
	}
	log.WithFields(log.Fields{"func": "discovery.checkv2"}).Debug(fmt.Sprintf("Got etcd cluster with endpoints: %v", c2.Endpoints()))
	kapi := client.NewKeysAPI(c2)
	kprefix := types.LegacyKubernetesPrefix
	kfound, _ := kapi.Get(context.Background(), kprefix, nil)
	if kfound != nil { // legacy v2 keyspace found
		return kprefix, nil
	}
	kprefix = types.KubernetesPrefix
	kfound, _ = kapi.Get(context.Background(), kprefix, nil)
	if kfound != nil { // modern v2 keyspace found
		return kprefix, nil
	}
	return "", nil
}

func getKubeDistrov2(endpoint string, secure bool) (distro types.KubernetesDistro, err error) {
	distro = types.NotADistro
	c2, err := util.NewClient2(endpoint, secure)
	if err != nil {
		return 0, err
	}
	kapi := client.NewKeysAPI(c2)
	_, err = kapi.Get(context.Background(), types.LegacyKubernetesPrefix, nil)
	if err == nil {
		distro = types.Vanilla
	}
	_, err = kapi.Get(context.Background(), types.KubernetesPrefix, nil)
	if err == nil {
		distro = types.Vanilla
	}
	_, err = kapi.Get(context.Background(), types.OpenShiftPrefix, nil)
	if err == nil {
		distro = types.OpenShift
	}
	return distro, nil
}

func getKubeDistrov3(endpoint string, secure bool) (distro types.KubernetesDistro, err error) {
	c3, cerr := util.NewClient3(endpoint, secure)
	if cerr != nil {
		return distro, fmt.Errorf("%s", cerr)
	}
	defer func() { _ = c3.Close() }()
	_, err = c3.Get(context.Background(), types.LegacyKubernetesPrefix)
	if err == nil {
		distro = types.Vanilla
	}
	_, err = c3.Get(context.Background(), types.KubernetesPrefix)
	if err != nil {
		distro = types.Vanilla
	}
	_, err = c3.Get(context.Background(), types.OpenShiftPrefix)
	if err != nil {
		distro = types.OpenShift
	}
	return distro, nil
}

func statsv2(endpoint string, secure bool, startkey string) (numkeys, totalsize int, err error) {
	c2, err := util.NewClient2(endpoint, secure)
	if err != nil {
		return 0, 0, err
	}
	kapi := client.NewKeysAPI(c2)
	log.WithFields(log.Fields{"func": "discovery.statsv2"}).Debug(fmt.Sprintf("Got etcd cluster using v2 API with endpoints: %v", c2.Endpoints()))
	err = Visit2(kapi, startkey, "", func(path, val string, arg interface{}) error {
		numkeys++
		totalsize += len(val)
		return nil
	}, "")
	if err != nil {
		return 0, 0, err
	}
	return numkeys, totalsize, nil
}

func statsv3(endpoint string, secure bool, startkey, endkey string) (numkeys, totalsize int, err error) {
	c3, cerr := util.NewClient3(endpoint, secure)
	if cerr != nil {
		return 0, 0, cerr
	}
	// defer func() { _ = c3.Close() }()
	log.WithFields(log.Fields{"func": "discovery.statsv3"}).Debug(fmt.Sprintf("Got etcd3 cluster with endpoints %v", c3.Endpoints()))
	err = Visit3(c3, "", startkey, endkey, func(path, val string, arg interface{}) error {
		numkeys++
		totalsize += len(val)
		return nil
	}, "")
	if err != nil {
		return 0, 0, err
	}
	return numkeys, totalsize, nil
}
