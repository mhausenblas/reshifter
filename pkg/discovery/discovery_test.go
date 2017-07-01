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
	probetests = []struct {
		launchfunc func(string) error
		scheme     string
		port       string
		version    string
		secure     bool
	}{
		{util.Etcd2Up, "http", "4001", "2", false},
		{util.Etcd3Up, "http", "4001", "3", false},
		{util.Etcd2SecureUp, "https", "4001", "2", true},
		{util.Etcd3SecureUp, "https", "4001", "3", true},
	}
	notadistro      = []string{"/something"}
	vanilladistro   = []string{types.KubernetesPrefix}
	openshiftdistro = []string{types.KubernetesPrefix, types.OpenShiftPrefix}
	k8stests        = []struct {
		keys    []string
		version string
		secure  bool
		distro  types.KubernetesDistro
	}{
		{notadistro, "2", false, types.NotADistro},
		{vanilladistro, "2", false, types.Vanilla},
		{openshiftdistro, "2", false, types.OpenShift},
	}
)

func TestCountKeysFor(t *testing.T) {
	defer func() {
		_ = util.EtcdDown()
	}()
	port := "4001"
	wantk := 2
	wants := 11
	tetcd := "http://127.0.0.1:" + port
	err := util.Etcd2Up(port)
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
	_, err = kapi.Set(context.Background(), types.KubernetesPrefix+"/namespaces/kube-system", ".", &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
	if err != nil {
		t.Errorf("Can't create etcd entry: %s", err)
		return
	}
	_, err = kapi.Set(context.Background(), types.KubernetesPrefix+"/namespaces/default", "..........", &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
	if err != nil {
		t.Errorf("Can't create etcd entry: %s", err)
		return
	}
	k, s, err := CountKeysFor(tetcd, types.Vanilla)
	if err != nil {
		t.Error(err)
		return
	}
	if k != wantk {
		t.Errorf("discovery.CountKeysFor(%s) => (%d, %d) want (%d, %d)", tetcd, k, s, wantk, wants)
	}
}

func TestProbeEtcd(t *testing.T) {
	_ = os.Setenv("RS_ETCD_CLIENT_CERT", filepath.Join(util.Certsdir(), "client.pem"))
	_ = os.Setenv("RS_ETCD_CLIENT_KEY", filepath.Join(util.Certsdir(), "client-key.pem"))
	for _, tt := range probetests {
		testEtcdX(t, tt.launchfunc, tt.scheme, tt.port, tt.version, tt.secure)
	}
	_, _, err := ProbeEtcd("127.0.0.1")
	if err == nil {
		t.Error(err)
	}
	_, _, err = ProbeEtcd("127.0.0.1:2379")
	if err == nil {
		t.Error(err)
	}
}

func testEtcdX(t *testing.T, etcdLaunchFunc func(string) error, scheme string, port string, version string, secure bool) {
	defer func() {
		_ = util.EtcdDown()
	}()
	tetcd := "127.0.0.1:" + port
	err := etcdLaunchFunc(port)
	if err != nil {
		t.Errorf("Can't launch etcd at %s: %s", tetcd, err)
		return
	}
	v, s, err := ProbeEtcd(scheme + "://" + tetcd)
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.HasPrefix(v, version) || s != secure {
		t.Errorf("discovery.ProbeEtcd(%s://%s) => (%s, %t) want (%s.x.x, %t)", scheme, tetcd, v, s, version, secure)
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
	err := util.Etcd2Up("4001")
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
	for _, key := range keys {
		err = util.SetKV2(kapi, key, ".")
		if err != nil {
			t.Errorf("Can't create key %s: %s", key, err)
			return
		}
	}
	distrotype, err := ProbeKubernetesDistro(tetcd)
	if err != nil {
		t.Errorf("Can't determine Kubernetes distro: %s", err)
		return
	}
	if distrotype != distro {
		t.Errorf("discovery.ProbeKubernetesDistro(%s) with keys %s => %v want %v", tetcd, keys, distrotype, distro)
	}
}
