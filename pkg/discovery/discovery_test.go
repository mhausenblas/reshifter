package discovery

import (
	"strings"
	"testing"

	"github.com/mhausenblas/reshifter/pkg/util"
)

var (
	probetests = []struct {
		scheme  string
		port    string
		version string
		secure  bool
	}{
		{"http", "2379", "2", false},
		{"http", "2379", "3", false},
	}
)

func TestProbeEtcd(t *testing.T) {
	defer func() {
		_ = util.Etcddown()
	}()
	// etcd2 discovery:
	p := probetests[0].port
	tetcd := "localhost:" + p
	err := util.Etcd2up(p)
	if err != nil {
		t.Errorf("Can't launch etcd at %s: %s", tetcd, err)
		return
	}
	v, s, err := ProbeEtcd(probetests[0].scheme + "://" + tetcd)
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.HasPrefix(v, probetests[0].version) || s != probetests[0].secure {
		t.Errorf("discovery.ProbeEtcd(%s) => (%s, %t) want (%s.x.x, %t)", tetcd, v, s, probetests[0].version, probetests[0].secure)
	}
	_ = util.Etcddown()
	// etcd3 discovery:
	p = probetests[1].port
	tetcd = "localhost:" + p
	err = util.Etcd3up(p)
	if err != nil {
		t.Errorf("Can't launch etcd at %s: %s", tetcd, err)
		return
	}
	v, s, err = ProbeEtcd(probetests[1].scheme + "://" + tetcd)
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.HasPrefix(v, probetests[1].version) || s != probetests[1].secure {
		t.Errorf("discovery.ProbeEtcd(%s) => (%s, %t) want (%s.x.x, %t)", tetcd, v, s, probetests[1].version, probetests[1].secure)
	}
}
