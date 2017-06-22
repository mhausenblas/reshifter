package etcd

import (
	"os"
	"testing"

	"github.com/coreos/etcd/client"
)

func TestRestore(t *testing.T) {
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
	based, err := Backup(tetcd)
	if err != nil {
		t.Errorf("Error during backup: %s", err)
	}

	target := "/tmp"
	afile := based
	_, err = Restore(afile, target, tetcd)
	if err != nil {
		t.Errorf("Error during restore: %s", err)
	}

	// make sure to clean up:
	_ = os.Remove(based + ".zip")
	_ = etcddown()
}
