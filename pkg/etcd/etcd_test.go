package etcd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

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

func TestBackup(t *testing.T) {
	tetcd := "localhost:2379"
	version := fmt.Sprintf("http://%s/version", tetcd)
	// check if local test etcd is available, otherwise abort right here:
	res, err := http.Get(version)
	if err != nil {
		t.Errorf("Can't connect to local etcd at %s. Run e2e-test/etcd-up.sh to launch it and try again", tetcd)
		return
	}
	j, _ := ioutil.ReadAll(res.Body)
	t.Logf("Got %s from %s", j, version)
	_ = res.Body.Close()
	// create some key-value pairs
	kapi, _ := newKeysAPI(tetcd)
	setkv(t, kapi, "/foo", "some")
	setkv(t, kapi, "/that", "")
	setkv(t, kapi, "/that/here", "moar")
	based, err := Backup(tetcd)
	if err != nil {
		t.Errorf("Error during backup: %s", err)
	}
	// TODO: check if content is as expected
	_, err = os.Stat(based + ".zip")
	if err != nil {
		t.Errorf("No archive found: %s", err)
	}
	t.Logf("Note: you can now run e2e-test/etcd-down.sh to shut down local etcd")
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

func setkv(t *testing.T, kapi client.KeysAPI, key, val string) {
	if val == "" {
		_, err := kapi.Set(context.Background(), key, "", &client.SetOptions{Dir: true, PrevExist: client.PrevIgnore})
		if err != nil {
			t.Errorf("Can't set key %s: %s", key, err)
		}
		return
	}
	_, err := kapi.Set(context.Background(), key, val, &client.SetOptions{Dir: false, PrevExist: client.PrevIgnore})
	if err != nil {
		t.Errorf("Can't set key %s with value %s: %s", key, val, err)
	}
}

func newKeysAPI(endpoint string) (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints:               []string{"http://" + endpoint},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("Can't create etcd client: %s", err)
	}
	kapi := client.NewKeysAPI(c)
	return kapi, nil
}
