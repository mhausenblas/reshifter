package discovery

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
)

var (
	statstests = []struct {
		version    string
		apiversion string
		start      string
		end        string
	}{
		{"2.x", "v2", types.LegacyKubernetesPrefix, types.LegacyKubernetesPrefixLast},
		{"2.x", "v2", types.KubernetesPrefix, types.KubernetesPrefixLast},
		{"3.x", "v2", types.LegacyKubernetesPrefix, types.LegacyKubernetesPrefixLast},
		{"3.x", "v2", types.KubernetesPrefix, types.KubernetesPrefixLast},
		{"3.x", "v3", types.LegacyKubernetesPrefix, types.LegacyKubernetesPrefixLast},
		{"3.x", "v3", types.KubernetesPrefix, types.KubernetesPrefixLast},
	}
	probetests = []struct {
		launchfunc func(string, string) (bool, error)
		scheme     string
		port       string
		version    string
		secure     bool
	}{
		{util.LaunchEtcd2, "http", "4001", "2", false},
		{util.LaunchEtcd3, "http", "4001", "3", false},
		{util.LaunchEtcd2, "https", "4001", "2", true},
		{util.LaunchEtcd3, "https", "4001", "3", true},
	}
	k8stests = []struct {
		keys    []string
		version string
		secure  bool
		distro  types.KubernetesDistro
	}{
		{[]string{""}, "2", false, types.NotADistro},
		{[]string{"/something"}, "2", false, types.NotADistro},
		{[]string{types.LegacyKubernetesPrefix}, "2", false, types.Vanilla},
		{[]string{types.KubernetesPrefix}, "2", false, types.Vanilla},
		{[]string{types.LegacyKubernetesPrefix, types.OpenShiftPrefix}, "2", false, types.OpenShift},
	}
)

func TestCountKeysFor(t *testing.T) {
	for _, tt := range statstests {
		testCountX(t, tt.version, tt.apiversion, tt.start, tt.end)
	}
}

func testCountX(t *testing.T, version, apiversion, start, end string) {
	defer func() {
		_ = util.EtcdDown()
	}()
	port := "4001"
	tetcd := "http://127.0.0.1:" + port
	wantkeys := 2
	wantsize := 11
	switch {
	case strings.HasPrefix(version, "3"): // etcd3 server
		_, err := util.LaunchEtcd3(tetcd, port)
		if err != nil {
			t.Errorf("Can't launch etcd at %s: %s", tetcd, err)
			return
		}
		switch apiversion {
		case "v2": // a v2 API in an etcd3
			c2, cerr := util.NewClient2(tetcd, false)
			if cerr != nil {
				t.Errorf("Can't connect to local etcd2 at %s: %s", tetcd, err)
				return
			}
			kapi := client.NewKeysAPI(c2)
			_, err = kapi.Set(context.Background(), start+"/namespaces/kube-system", ".", &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
			if err != nil {
				t.Errorf("Can't create etcd entry: %s", err)
				return
			}
			_, err = kapi.Set(context.Background(), start+"/namespaces/default", "..........", &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
			if err != nil {
				t.Errorf("Can't create etcd entry: %s", err)
				return
			}
		case "v3":
			c3, err := util.NewClient3(tetcd, false)
			if err != nil {
				t.Errorf("Can't connect to local etcd3 at %s: %s", tetcd, err)
				return
			}
			_, err = c3.Put(context.Background(), start+"/namespaces/kube-system", ".")
			if err != nil {
				t.Errorf("Can't create key %s/namespaces/kube-system: %s", start, err)
				return
			}
			_, err = c3.Put(context.Background(), start+"/namespaces/default", "..........")
			if err != nil {
				t.Errorf("Can't create key %s/namespaces/default: %s", start, err)
				return
			}
		}
	case strings.HasPrefix(version, "2"): // etcd2 server
		_, err := util.LaunchEtcd2(tetcd, port)
		if err != nil {
			t.Errorf("Can't launch etcd at %s: %s", tetcd, err)
			return
		}
		c2, err := util.NewClient2(tetcd, false)
		if err != nil {
			t.Errorf("Can't connect to local etcd2 at %s: %s", tetcd, err)
			return
		}
		kapi := client.NewKeysAPI(c2)
		_, err = kapi.Set(context.Background(), start+"/namespaces/kube-system", ".", &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
		if err != nil {
			t.Errorf("Can't create etcd entry: %s", err)
			return
		}
		_, err = kapi.Set(context.Background(), start+"/namespaces/default", "..........", &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
		if err != nil {
			t.Errorf("Can't create etcd entry: %s", err)
			return
		}
	default:
		t.Errorf("Can't understand etcd version, seems to be neither v3 nor v2 :(")
		return
	}
	k, s, err := CountKeysFor(tetcd, start, end)
	if err != nil {
		t.Error(err)
		return
	}
	if k != wantkeys {
		t.Errorf("discovery.CountKeysFor(%s, %s, %s) for an etcd %s using API %s => got (%d, %d) but want (%d, %d)", tetcd, start, end, version, apiversion, k, s, wantkeys, wantsize)
	}
}

func TestProbeEtcd(t *testing.T) {
	_ = os.Setenv("RS_ETCD_CLIENT_CERT", filepath.Join(util.Certsdir(), "client.pem"))
	_ = os.Setenv("RS_ETCD_CLIENT_KEY", filepath.Join(util.Certsdir(), "client-key.pem"))
	for _, tt := range probetests {
		testEtcdX(t, tt.launchfunc, tt.scheme, tt.port, tt.version, tt.secure)
	}
	_, _, _, err := ProbeEtcd("127.0.0.1")
	if err == nil {
		t.Error(err)
	}
	_, _, _, err = ProbeEtcd("127.0.0.1:2379")
	if err == nil {
		t.Error(err)
	}
}

func testEtcdX(t *testing.T, etcdLaunchFunc func(string, string) (bool, error), scheme string, port string, version string, secure bool) {
	defer func() {
		_ = util.EtcdDown()
	}()
	tetcd := scheme + "://127.0.0.1:" + port
	_, err := etcdLaunchFunc(tetcd, port)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	v, _, s, err := ProbeEtcd(tetcd)
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.HasPrefix(v, version) || s != secure {
		t.Errorf("discovery.ProbeEtcd(%s) => got (%s, %t) but want (%s.x.x, %t)", tetcd, v, s, version, secure)
	}
}

func TestProbeKubernetesDistro(t *testing.T) {
	for _, tt := range k8stests {
		testK8SX(t, tt.keys, tt.version, tt.secure, tt.distro)
		time.Sleep(time.Second * 1)
	}
}

func testK8SX(t *testing.T, keys []string, version string, secure bool, distro types.KubernetesDistro) {
	defer func() {
		_ = util.EtcdDown()
	}()
	tetcd := "http://127.0.0.1:4001"
	_, err := util.LaunchEtcd2(tetcd, "4001")
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	c2, err := util.NewClient2(tetcd, false)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	kapi := client.NewKeysAPI(c2)
	for _, key := range keys {
		if key != "" {
			err = util.SetKV2(kapi, key, ".")
			if err != nil {
				t.Errorf("%s", err)
				return
			}
		}
	}
	distrotype, err := ProbeKubernetesDistro(tetcd)
	if err != nil {
		t.Errorf("Can't determine Kubernetes distro: %s", err)
		return
	}
	if distrotype != distro {
		t.Errorf("discovery.ProbeKubernetesDistro(%s) with keys %s => got '%s' but want '%s'", tetcd, keys, util.LookupDistro(distrotype), util.LookupDistro(distro))
	}
}
