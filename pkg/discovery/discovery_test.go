package discovery

import (
	"strings"
	"testing"

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
		{util.Etcd2Up, "http", "2379", "2", false},
		{util.Etcd3Up, "http", "2379", "3", false},
		{util.Etcd2SecureUp, "https", "2379", "2", true},
	}
)

func TestProbeEtcd(t *testing.T) {
	for _, tt := range probetests {
		testEtcdX(t, tt.launchfunc, tt.scheme, tt.port, tt.version, tt.secure)
	}
	_, _, err := ProbeEtcd("localhost")
	if err == nil {
		t.Error(err)
	}
	_, _, err = ProbeEtcd("localhost:2379")
	if err == nil {
		t.Error(err)
	}
}

func testEtcdX(t *testing.T, etcdLaunchFunc func(string) error, scheme string, port string, version string, secure bool) {
	defer func() {
		_ = util.EtcdDown()
	}()
	tetcd := "localhost:" + port
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
