package discovery

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"net/http"
	"net/url"

	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
)

// ProbeEtcd probes an endpoint at path /version to figure
// which version of etcd it is and in which mode (secure or insecure)
// it is used. Example:
//
//		version, secure, err := ProbeEtcd("http://localhost:2379")
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

func getVersion(endpoint string) (string, error) {
	var etcdr types.EtcdResponse
	res, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("Can't query %s endpoint: %s", endpoint, err)
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
		return "", fmt.Errorf("Can't query %s endpoint: %s", endpoint, err)
	}
	err = json.NewDecoder(res.Body).Decode(&etcdr)
	if err != nil {
		return "", fmt.Errorf("Can't decode response from etcd: %s", err)
	}
	_ = res.Body.Close()
	return etcdr.EtcdServerVersion, nil
}
