package etcd

import "testing"

var storetests = []struct {
	inbased string
	inpath  string
	inval   string
	out     string
}{
	{".", "test", "a", "./test/a"},
	{".", "test", "b", "./test/b"},
}

func TestStore(t *testing.T) {
	for _, tt := range storetests {
		store(tt.inbased, tt.inpath, tt.inval)
		s := "./test/a"
		if s != tt.out {
			t.Errorf("etcd.store(%q, %q, %q) => %q, want %q", tt.inbased, tt.inpath, tt.inval, s, tt.out)
		}
	}
}
