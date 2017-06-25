package backup

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/coreos/etcd/client"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
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
			t.Errorf("backup.store(\".\", %q, %q) => %q, want %q", tt.path, tt.val, got, want)
		}
	}
	// make sure to clean up remaining directories:
	_ = os.RemoveAll(tmpTestDir)
}

func TestBackup(t *testing.T) {
	port := "4001"
	// testing insecure etcd 2 and 3:
	tetcd := "http://localhost:" + port
	etcd2Backup(t, port, tetcd)
	etcd3Backup(t, port, tetcd)
	// testing secure etcd 2 and 3:
	// tetcd := "https://localhost:" + port
	// TBD
}

func etcd2Backup(t *testing.T, port, tetcd string) {
	defer func() { _ = util.EtcdDown() }()
	err := util.Etcd2Up(port)
	if err != nil {
		t.Errorf("Can't launch local etcd at %s: %s", tetcd, err)
		return
	}
	c2, err := util.NewClient2(tetcd, false)
	if err != nil {
		t.Errorf("Can't connect to local etcd2 at %s: %s", tetcd, err)
		return
	}
	kapi := client.NewKeysAPI(c2)
	err = util.SetKV2(kapi, "/foo", "some")
	if err != nil {
		t.Errorf("Can't create key /foo: %s", err)
		return
	}
	err = util.SetKV2(kapi, "/that/here", "moar")
	if err != nil {
		t.Errorf("Can't create key /that/here: %s", err)
		return
	}
	based, err := Backup(tetcd, "/tmp")
	if err != nil {
		t.Errorf("Error during backup: %s", err)
		return
	}
	opath, _ := filepath.Abs(filepath.Join("/tmp", based))
	_, err = os.Stat(opath + ".zip")
	if err != nil {
		t.Errorf("No archive found: %s", err)
	}
	// make sure to clean up:
	_ = os.Remove(opath + ".zip")
}

func etcd3Backup(t *testing.T, port, tetcd string) {
	defer func() { _ = util.EtcdDown() }()
	err := util.Etcd3Up(port)
	if err != nil {
		t.Errorf("Can't launch local etcd at %s: %s", tetcd, err)
		return
	}
	c3, err := util.NewClient3(tetcd, false)
	if err != nil {
		t.Errorf("Can't connect to local etcd3 at %s: %s", tetcd, err)
		return
	}

	pr, err := c3.Put(context.Background(), types.KubernetesPrefix+"namespaces/kube-system", "{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}")
	if err != nil {
		t.Errorf("Can't create key %snamespaces/kube-system: %s", types.KubernetesPrefix, err)
		return
	}
	_ = pr
	// t.Logf(fmt.Sprintf("PUT response: %v", pr))

	// res, err := c3.Get(context.Background(), types.KubernetesPrefix+"*", clientv3.WithRange(types.KubernetesPrefixLast))
	// if err != nil {
	// 	t.Errorf("No value retrieved: %s", err)
	// 	return
	// }
	// t.Logf(fmt.Sprintf("GET response: %v", res.Kvs))

	based, err := Backup(tetcd, "/tmp")
	if err != nil {
		t.Errorf("Error during backup: %s", err)
		return
	}
	opath, _ := filepath.Abs(filepath.Join("/tmp", based))
	_, err = os.Stat(opath + ".zip")
	if err != nil {
		t.Errorf("No archive found: %s", err)
	}
	// make sure to clean up:
	_ = os.Remove(opath + ".zip")

}
