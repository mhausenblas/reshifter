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
func ProbeEtcd(endpoint string) (string, bool, error) {
	u, err := url.Parse(endpoint + "/version")
	if err != nil {
		return "", false, fmt.Errorf("Can't parse endpoint %s: %s", endpoint, err)
	}
	if u.Scheme == "https" { // secure etcd
		clientcert, clientkey, err := util.ClientCertAndKeyFromEnv()
		if err != nil {
			return "", false, err
		}
		version, verr := getVersionSecure(u.String(), clientcert, clientkey)
		if verr != nil {
			return "", false, verr
		}
		return version, true, nil
	}
	version, verr := getVersion(u.String())
	if verr != nil {
		return "", false, verr
	}
	return version, false, nil
}

// ProbeKubernetesDistro probes an etcd cluster for which Kubernetes
// distribution is present by scanning the available keys.
func ProbeKubernetesDistro(endpoint string) (types.KubernetesDistro, error) {
	distro := types.NotADistro
	version, secure, err := ProbeEtcd(endpoint)
	if err != nil {
		return types.NotADistro, fmt.Errorf("%s", err)
	}
	// deal with etcd3 servers:
	if strings.HasPrefix(version, "3") {
		c3, cerr := util.NewClient3(endpoint, secure)
		if cerr != nil {
			return types.NotADistro, fmt.Errorf("Can't connect to etcd: %s", cerr)
		}
		defer func() { _ = c3.Close() }()
		_, err := c3.Get(context.Background(), types.KubernetesPrefix)
		if err != nil {
			return distro, nil
		}
		distro = types.Vanilla
		_, err = c3.Get(context.Background(), types.OpenShiftPrefix)
		if err != nil {
			return distro, nil
		}
		distro = types.OpenShift
		return distro, nil
	}
	// deal with etcd2 servers:
	if strings.HasPrefix(version, "2") {
		c2, cerr := util.NewClient2(endpoint, secure)
		if cerr != nil {
			return types.NotADistro, fmt.Errorf("Can't connect to etcd: %s", cerr)
		}
		kapi := client.NewKeysAPI(c2)
		_, err := kapi.Get(context.Background(), types.KubernetesPrefix, nil)
		if err != nil {
			return distro, nil
		}
		distro = types.Vanilla
		_, err = kapi.Get(context.Background(), types.OpenShiftPrefix, nil)
		if err != nil {
			return distro, nil
		}
		distro = types.OpenShift
		return distro, nil
	}
	return types.NotADistro, fmt.Errorf("Can't determine Kubernetes distro")
}

// CountKeysFor iterates over well-known keys of a given Kubernetes distro
// and returns the number of keys and their values total size, in the
// respective key range/subtree.
func CountKeysFor(endpoint string, distro types.KubernetesDistro) (int, int, error) {
	numkeys := 0
	totalsize := 0
	version, secure, err := ProbeEtcd(endpoint)
	if err != nil {
		return 0, 0, fmt.Errorf("%s", err)
	}
	// deal with etcd3 servers:
	if strings.HasPrefix(version, "3") {
		c3, cerr := util.NewClient3(endpoint, secure)
		if cerr != nil {
			return 0, 0, fmt.Errorf("Can't connect to etcd: %s", cerr)
		}
		defer func() { _ = c3.Close() }()
		log.WithFields(log.Fields{"func": "discovery.CountKeysFor"}).Debug(fmt.Sprintf("Got etcd3 cluster with endpoints %v", c3.Endpoints()))
		switch distro {
		case types.Vanilla:
			err = Visit3(c3, types.KubernetesPrefix, "", types.Vanilla, func(path, val string, arg interface{}) error {
				numkeys++
				totalsize += len(val)
				return nil
			}, "")
			if err != nil {
				return 0, 0, err
			}
		case types.OpenShift:
			err = Visit3(c3, types.OpenShiftPrefix, "", types.OpenShift, func(path, val string, arg interface{}) error {
				numkeys++
				totalsize += len(val)
				return nil
			}, "")
			if err != nil {
				return 0, 0, err
			}
		}
	}
	// deal with etcd2 servers:
	if strings.HasPrefix(version, "2") {
		c2, cerr := util.NewClient2(endpoint, secure)
		if cerr != nil {
			return 0, 0, fmt.Errorf("Can't connect to etcd: %s", cerr)
		}
		kapi := client.NewKeysAPI(c2)
		log.WithFields(log.Fields{"func": "discovery.CountKeysFor"}).Debug(fmt.Sprintf("Got etcd2 cluster with %v", c2.Endpoints()))
		switch distro {
		case types.Vanilla:
			err = Visit2(kapi, types.KubernetesPrefix, "", func(path, val string, arg interface{}) error {
				numkeys++
				totalsize += len(val)
				return nil
			}, "")
			if err != nil {
				return 0, 0, err
			}
		case types.OpenShift:
			err = Visit2(kapi, types.OpenShiftPrefix, "", func(path, val string, arg interface{}) error {
				numkeys++
				totalsize += len(val)
				return nil
			}, "")
			if err != nil {
				return 0, 0, err
			}
		}
	}
	return numkeys, totalsize, nil
}

func getVersion(endpoint string) (string, error) {
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
	return etcdr.EtcdServerVersion, nil
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
