package discovery

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ProbeEtcd probes an endpoint to figure which version of etcd
// and in which mode (secure or insecure) it is used.
func ProbeEtcd(endpoint string) (string, bool, error) {
	issecure := false
	version := "2"
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", false, fmt.Errorf("Can't parse endpoint %s: %s", endpoint, err)
	}
	if u.Scheme == "https" {
		issecure = true
	}
	res, err := http.Get(u.String())
	if err != nil {
		return "", false, fmt.Errorf("Can't connect to etcd at %s: %s", endpoint, err)
	}
	j, _ := ioutil.ReadAll(res.Body)
	_ = j
	_ = res.Body.Close()
	return version, issecure, nil
}
