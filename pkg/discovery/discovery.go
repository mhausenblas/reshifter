package discovery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mhausenblas/reshifter/pkg/types"
)

// ProbeEtcd probes an endpoint at path /version to figure
// which version of etcd it is and in which mode (secure or insecure)
// it is used.
func ProbeEtcd(endpoint string) (string, bool, error) {
	issecure := false
	u, err := url.Parse(endpoint + "/version")
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
	var etcd2r types.Etcd2Response
	err = json.NewDecoder(res.Body).Decode(&etcd2r)
	_ = res.Body.Close()
	return etcd2r.EtcdServerVersion, issecure, nil
}
