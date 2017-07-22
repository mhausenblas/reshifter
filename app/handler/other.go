package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/discovery"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
)

// Version responds to HTTP GET requests of the form:
//		http GET localhost:8080/v1/version
func Version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ReShifter in version %s", ReleaseVersion)
}

// EPstats responds to HTTP GET requests such as:
//		http GET localhost:8080/v1/epstats?endpoint=http%3A%2F%2Flocalhost%3A2379
func EPstats(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Query().Get("endpoint")
	if endpoint == "" || strings.Index(endpoint, "http") != 0 {
		merr := "The endpoint is malformed"
		http.Error(w, merr, http.StatusBadRequest)
		log.Error(merr)
		return
	}
	vlk, vls, err := discovery.CountKeysFor(endpoint, types.LegacyKubernetesPrefix, types.LegacyKubernetesPrefixLast)
	if err != nil {
		log.Info("Didn't find legacy keys, trying modern keys now")
	}
	vk, vs, err := discovery.CountKeysFor(endpoint, types.KubernetesPrefix, types.KubernetesPrefixLast)
	if err != nil {
		merr := fmt.Sprintf("Having problems calculating stats, no vanialla keys found: %s", err)
		http.Error(w, merr, http.StatusInternalServerError)
		log.Error(merr)
		return
	}
	log.Debugf("vanilla [keys:%d, size:%d]", vlk+vk, vls+vs)
	// note: ignoring error here since we're adding up the stats
	// and if this happens to be a non-OpenShift distro we simply
	// add 0 to the overall count, and it's still fine:
	osk, oss, _ := discovery.CountKeysFor(endpoint, types.OpenShiftPrefix, types.OpenShiftPrefixLast)
	log.Debugf("openshift [keys:%d, size:%d]", osk, oss)
	_ = json.NewEncoder(w).Encode(struct {
		NumKeys         int `json:"numkeys"`
		TotalSizeValues int `json:"totalsizevalbytes"`
	}{
		vlk + vk + osk,
		vls + vs + oss,
	})
}

// Explorer responds to HTTP GET requests such as:
//		http GET localhost:8080/v1/explorer?endpoint=http%3A%2F%2Flocalhost%3A2379
func Explorer(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Query().Get("endpoint")
	if endpoint == "" || strings.Index(endpoint, "http") != 0 {
		merr := "The endpoint is malformed"
		http.Error(w, merr, http.StatusBadRequest)
		log.Error(merr)
		return
	}
	version, apiversion, issecure, err := discovery.ProbeEtcd(endpoint)
	if err != nil {
		merr := fmt.Sprintf("%s", err)
		http.Error(w, merr, http.StatusBadRequest)
		log.Error(err)
		return
	}
	distrotype, err := discovery.ProbeKubernetesDistro(endpoint)
	if err != nil {
		merr := fmt.Sprintf("Can't determine Kubernetes distro: %s", err)
		http.Error(w, merr, http.StatusBadRequest)
		log.Error(err)
		return
	}
	secure := "insecure etcd, no SSL/TLS configured"
	if issecure {
		secure = "secure etcd, SSL/TLS configured"
	}
	_ = json.NewEncoder(w).Encode(struct {
		EtcdVersion  string `json:"etcdversion"`
		APIVersion   string `json:"apiversion"`
		EtcdSecurity string `json:"etcdsec"`
		K8SDistro    string `json:"k8sdistro"`
	}{
		version,
		apiversion,
		secure,
		util.LookupDistro(distrotype),
	})
}
