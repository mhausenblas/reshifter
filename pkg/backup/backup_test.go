package backup

import (
	"context"
	"fmt"
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

func TestBackupv2(t *testing.T) {
	port := "4001"
	// testing insecure etcd 2:
	tetcd := "http://127.0.0.1:" + port
	// backing up to remote https://play.minio.io:9000:
	_ = os.Setenv("ACCESS_KEY_ID", "Q3AM3UQ867SPQQA43P2F")
	_ = os.Setenv("SECRET_ACCESS_KEY", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")
	etcd2Backup(t, port, tetcd, types.LegacyKubernetesPrefix+"/namespaces/kube-system")
	etcd2Backup(t, port, tetcd, types.KubernetesPrefix+"/namespaces/kube-system")
	etcd2Backup(t, port, tetcd, types.OpenShiftPrefix+"/builds")
}

func TestBackupv3(t *testing.T) {
	port := "4001"
	// testing insecure etcd3:
	tetcd := "http://127.0.0.1:" + port
	// backing up to remote https://play.minio.io:9000:
	_ = os.Setenv("ACCESS_KEY_ID", "Q3AM3UQ867SPQQA43P2F")
	_ = os.Setenv("SECRET_ACCESS_KEY", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")
	etcd3Backup(t, port, tetcd, types.LegacyKubernetesPrefix+"/namespaces/kube-system")
	etcd3Backup(t, port, tetcd, types.KubernetesPrefix+"/namespaces/kube-system")
	etcd3Backup(t, port, tetcd, types.OpenShiftPrefix+"/builds")
}

func TestBackupSecure(t *testing.T) {
	port := "4001"
	// testing secure etcd 2:
	tetcd := "https://127.0.0.1:" + port
	// backing up to remote https://play.minio.io:9000:
	_ = os.Setenv("ACCESS_KEY_ID", "Q3AM3UQ867SPQQA43P2F")
	_ = os.Setenv("SECRET_ACCESS_KEY", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")
	etcd2Backup(t, port, tetcd, types.LegacyKubernetesPrefix+"/namespaces/kube-system")
}

func etcd2Backup(t *testing.T, port, tetcd, key string) {
	defer func() { _ = util.EtcdDown() }()
	_ = os.Setenv("RS_ETCD_CLIENT_CERT", filepath.Join(util.Certsdir(), "client.pem"))
	_ = os.Setenv("RS_ETCD_CLIENT_KEY", filepath.Join(util.Certsdir(), "client-key.pem"))
	_ = os.Setenv("RS_ETCD_CA_CERT", filepath.Join(util.Certsdir(), "ca.pem"))
	secure, err := util.LaunchEtcd2(tetcd, port)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	c2, err := util.NewClient2(tetcd, secure)
	if err != nil {
		t.Errorf("Can't connect to local etcd2 at %s: %s", tetcd, err)
		return
	}
	kapi := client.NewKeysAPI(c2)
	val := "{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}"
	_, err = kapi.Set(context.Background(), key, val, &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
	if err != nil {
		t.Errorf("Can't create etcd entry %s=%s: %s", key, val, err)
		return
	}
	if strings.HasPrefix(key, types.OpenShiftPrefix) { // make sure to add vanilla key as well
		_, err = kapi.Set(context.Background(), types.KubernetesPrefix+"/namespaces/kube-system", val, &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
		if err != nil {
			t.Errorf("Can't create etcd entry %s=%s: %s", key, val, err)
			return
		}
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

func etcd3Backup(t *testing.T, port, tetcd, key string) {
	defer func() { _ = util.EtcdDown() }()
	_ = os.Setenv("ETCDCTL_API", "3")
	_ = os.Setenv("RS_ETCD_CLIENT_CERT", filepath.Join(util.Certsdir(), "client.pem"))
	_ = os.Setenv("RS_ETCD_CLIENT_KEY", filepath.Join(util.Certsdir(), "client-key.pem"))
	_ = os.Setenv("RS_ETCD_CA_CERT", filepath.Join(util.Certsdir(), "ca.pem"))
	secure, err := util.LaunchEtcd3(tetcd, port)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	c3, err := util.NewClient3(tetcd, secure)
	if err != nil {
		t.Errorf("Can't connect to local etcd3 at %s: %s", tetcd, err)
		return
	}
	val := "{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}"
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	_, err = c3.Put(context.Background(), key, val)
	if err != nil {
		t.Errorf("Can't create etcd entry %s=%s: %s", key, val, err)
		return
	}
	if strings.HasPrefix(key, types.OpenShiftPrefix) { // make sure to add vanilla key as well
		_, err = c3.Put(context.Background(), types.KubernetesPrefix+"/namespaces/kube-system", val)
		if err != nil {
			t.Errorf("Can't create etcd entry %s=%s: %s", key, val, err)
			return
		}
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

func genentry(etcdversion string, distro types.KubernetesDistro) (string, string, error) {
	switch distro {
	case types.Vanilla:
		if etcdversion == "2" {
			return types.LegacyKubernetesPrefix + "/namespaces/kube-system", "{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}", nil
		}
		return types.KubernetesPrefix + "/namespaces/kube-system", "{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}", nil
	case types.OpenShift:
		return types.OpenShiftPrefix + "/builds", "{\"kind\":\"Build\",\"apiVersion\":\"v1\"}", nil
	default:
		return "", "", fmt.Errorf("That's not a Kubernetes distro")
	}
}
