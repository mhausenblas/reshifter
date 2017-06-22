package discovery

import (
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
	}
)

func TestProbeEtcd(t *testing.T) {
	defer func() {
		_ = util.Etcddown()
	}()
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
	if v != probetests[0].version || s != probetests[0].secure {
		t.Errorf("discovery.ProbeEtcd(%s) => (%s, %t) want (%s, %t)", tetcd, v, s, probetests[0].version, probetests[0].secure)
	}

}
