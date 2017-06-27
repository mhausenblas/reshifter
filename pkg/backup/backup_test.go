package backup

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	// backing up to remote https://play.minio.io:9000:
	_ = os.Setenv("ACCESS_KEY_ID", "Q3AM3UQ867SPQQA43P2F")
	_ = os.Setenv("SECRET_ACCESS_KEY", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")
	etcd2Backup(t, port, tetcd)
	etcd3Backup(t, port, tetcd)
	// testing secure etcd 2 and 3:
	tetcd = "https://localhost:" + port
	etcd2Backup(t, port, tetcd)
	etcd3Backup(t, port, tetcd)
}

func etcd2Backup(t *testing.T, port, tetcd string) {
	defer func() { _ = util.EtcdDown() }()
	secure := false
	switch {
	case strings.Index(tetcd, "https") == 0:
		err := util.Etcd2SecureUp(port)
		secure = true
		_ = os.Setenv("RS_ETCD_CLIENT_CERT", filepath.Join(util.Certsdir(), "client.pem"))
		_ = os.Setenv("RS_ETCD_CLIENT_KEY", filepath.Join(util.Certsdir(), "client-key.pem"))
		_ = os.Setenv("RS_ETCD_CA_CERT", filepath.Join(util.Certsdir(), "ca.pem"))
		if err != nil {
			t.Errorf("Can't launch secure etcd2 at %s: %s", tetcd, err)
			return
		}
	case strings.Index(tetcd, "http") == 0:
		err := util.Etcd2Up(port)
		if err != nil {
			t.Errorf("Can't launch insecure etcd2 at %s: %s", tetcd, err)
			return
		}
	default:
		t.Errorf("That's not a valid etcd2 endpoint: %s", tetcd)
		return
	}
	c2, err := util.NewClient2(tetcd, secure)
	if err != nil {
		t.Errorf("Can't connect to local etcd2 at %s: %s", tetcd, err)
		return
	}
	kapi := client.NewKeysAPI(c2)
	_, err = kapi.Set(context.Background(),
		types.KubernetesPrefix+"/namespaces/kube-system",
		"{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}",
		&client.SetOptions{Dir: false, PrevExist: client.PrevNoExist},
	)
	if err != nil {
		t.Errorf("Can't create key %snamespaces/kube-system: %s", types.KubernetesPrefix, err)
		return
	}
	backupid, err := Backup(tetcd, types.DefaultWorkDir, "play.minio.io:9000", "reshifter-test-cluster")
	if err != nil {
		t.Errorf("Error during backup: %s", err)
		return
	}
	opath, _ := filepath.Abs(filepath.Join(types.DefaultWorkDir, backupid))
	_, err = os.Stat(opath + ".zip")
	if err != nil {
		t.Errorf("No archive found: %s", err)
	}
	// make sure to clean up:
	_ = os.Remove(opath + ".zip")
}

func etcd3Backup(t *testing.T, port, tetcd string) {
	defer func() { _ = util.EtcdDown() }()
	_ = os.Setenv("ETCDCTL_API", "3")
	secure := false
	switch {
	case strings.Index(tetcd, "https") == 0:
		err := util.Etcd3SecureUp(port)
		secure = true
		_ = os.Setenv("RS_ETCD_CLIENT_CERT", filepath.Join(util.Certsdir(), "client.pem"))
		_ = os.Setenv("RS_ETCD_CLIENT_KEY", filepath.Join(util.Certsdir(), "client-key.pem"))
		_ = os.Setenv("RS_ETCD_CA_CERT", filepath.Join(util.Certsdir(), "ca.pem"))
		if err != nil {
			t.Errorf("Can't launch secure etcd3 at %s: %s", tetcd, err)
			return
		}
	case strings.Index(tetcd, "http") == 0:
		err := util.Etcd3Up(port)
		if err != nil {
			t.Errorf("Can't launch insecure etcd3 at %s: %s", tetcd, err)
			return
		}
	default:
		t.Errorf("That's not a valid etcd2 endpoint: %s", tetcd)
		return
	}
	c3, err := util.NewClient3(tetcd, secure)
	if err != nil {
		t.Errorf("Can't connect to local etcd3 at %s: %s", tetcd, err)
		return
	}
	_, err = c3.Put(context.Background(),
		types.KubernetesPrefix+"/namespaces/kube-system",
		"{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}",
	)
	if err != nil {
		t.Errorf("Can't create key %snamespaces/kube-system: %s", types.KubernetesPrefix, err)
		return
	}
	backupid, err := Backup(tetcd, types.DefaultWorkDir, "play.minio.io:9000", "reshifter-test-cluster")
	if err != nil {
		t.Errorf("Error during backup: %s", err)
		return
	}
	opath, _ := filepath.Abs(filepath.Join(types.DefaultWorkDir, backupid))
	_, err = os.Stat(opath + ".zip")
	if err != nil {
		t.Errorf("No archive found: %s", err)
	}
	// make sure to clean up:
	_ = os.Remove(opath + ".zip")
}
