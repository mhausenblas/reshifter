package etcd

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

func TestBackup(t *testing.T) {
	tetcd := "localhost:2379"
	version := fmt.Sprintf("http://%s/version", tetcd)
	// check if local test etcd is available, otherwise abort right here:
	res, err := http.Get(version)
	if err != nil {
		t.Errorf("Can't connect to local etcd at %s. Run e2e-test/etcd-up.sh to launch it and try again.", tetcd)
	}
	j, _ := ioutil.ReadAll(res.Body)
	t.Logf("Got %s from %s", j, version)
	_ = res.Body.Close()
	err = Backup(tetcd)
	if err != nil {
		t.Errorf("Error during backup: %s.", err)
	}
	// TODO: pull in setup (populate etcd with content), check if content is as expected, hint at shutdown
}

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
