package etcd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/coreos/etcd/client"
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
		c, _ := ioutil.ReadFile(p)
		got := string(c)
		if tt.path == "/" {
			_ = os.Remove(p)
		}
		want := tt.val
		if got != want {
			t.Errorf("etcd.store(\".\", %q, %q) => %q, want %q", tt.path, tt.val, got, want)
		}
	}
	// make sure to clean up remaining directories:
	_ = os.RemoveAll(tmpTestDir)
}

func TestBackup(t *testing.T) {
	port := "2379"
	tetcd := "localhost:" + port
	err := etcd2up(port)
	if err != nil {
		t.Errorf("Can't launch local etcd at %s: %s", tetcd, err)
		return
	}
	// create some key-value pairs:
	c2, err := newClient2(tetcd, false)
	if err != nil {
		t.Errorf("Can't connect to local etcd2 at %s: %s", tetcd, err)
		return
	}

	kapi := client.NewKeysAPI(c2)
	err = setKV2(kapi, "/foo", "some")
	if err != nil {
		t.Errorf("Can't create key /foo: %s", err)
		return
	}
	err = setKV2(kapi, "/that", "")
	if err != nil {
		t.Errorf("Can't create key /that: %s", err)
		return
	}
	err = setKV2(kapi, "/that/here", "moar")
	if err != nil {
		t.Errorf("Can't create key /that/here: %s", err)
	}

	based, err := Backup(tetcd)
	if err != nil {
		t.Errorf("Error during backup: %s", err)
	}
	// TODO: check if content is as expected
	_, err = os.Stat(based + ".zip")
	if err != nil {
		t.Errorf("No archive found: %s", err)
	}
	// make sure to clean up:
	_ = os.Remove(based + ".zip")
	_ = etcddown()
}
