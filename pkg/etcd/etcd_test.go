package etcd

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	tmpTestDir = "test/"
	storetests = []struct {
		path string
		val  string
	}{
		{"", ""},
		{"non-valid-key", ""},
		{"/", "root"},
		{"/" + tmpTestDir, "some"},
		{"/" + tmpTestDir + "/first-level", "another"},
		{"/" + tmpTestDir + "/this:also", "escaped"},
	}
)

func TestStore(t *testing.T) {
	for _, tt := range storetests {
		p, err := store(".", tt.path, tt.val)
		if err != nil {
			continue
		}
		got := readcontent(p)
		want := tt.val
		if got != want {
			t.Errorf("etcd.store(\".\", %q, %q) => %q, want %q", tt.path, tt.val, got, want)
		}
	}
	// make sure to clean up remaining directories:
	_ = os.RemoveAll(tmpTestDir)
}

func readcontent(path string) string {
	// make sure to clean up individual files
	defer func() {
		if path != "." {
			_ = os.Remove(path)
		}
	}()
	content, _ := ioutil.ReadFile(path)
	return string(content)
}
