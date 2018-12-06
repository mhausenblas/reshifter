package backup

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
	"go.etcd.io/etcd/client"
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
	filtertests = []struct {
		filter string
		path   string
		exists bool
	}{
		{"filter:abc", "/abc", true},
		{"filter:abc", "/def", false},
		{"filter:name", "/some/name", true},
		{"filter:abc,def", "/some/other/def", true},
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
	etc3BackupWithLargeSubkeys(t, port, tetcd, types.KubernetesPrefix+"/test/bigsubkeys")
}

func TestBackupv2inv3(t *testing.T) {
	defer func() { _ = util.EtcdDown() }()
	port := "4001"
	// testing insecure etcd 3:
	tetcd := "http://127.0.0.1:" + port
	// backing up to remote https://play.minio.io:9000:
	_ = os.Setenv("ACCESS_KEY_ID", "Q3AM3UQ867SPQQA43P2F")
	_ = os.Setenv("SECRET_ACCESS_KEY", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")

	// adding key using v2 API in an etcd3
	_, err := util.LaunchEtcd3(tetcd, port)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	c2, err := util.NewClient2(tetcd, false)
	if err != nil {
		t.Errorf("Can't connect to local etcd3 at %s: %s", tetcd, err)
		return
	}
	kapi := client.NewKeysAPI(c2)
	key := types.LegacyKubernetesPrefix + "/namespaces/kube-system"
	val := "{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}"
	_, err = kapi.Set(context.Background(), key, val, &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
	if err != nil {
		t.Errorf("Can't create etcd entry %s=%s: %s", key, val, err)
		return
	}
	key = types.LegacyKubernetesPrefix + "/ThirdPartyResourceData/stable.example.com/crontabs/dohnto/my-new-cron-object"
	val = "{\"kind\":\"ThirdPartyResource\",\"apiVersion\":\"v1\"}"
	_, err = kapi.Set(context.Background(), key, val, &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
	if err != nil {
		t.Errorf("Can't create etcd entry %s=%s: %s", key, val, err)
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

func TestBackupSecure(t *testing.T) {
	port := "4001"
	// testing secure etcd 2:
	tetcd := "https://127.0.0.1:" + port
	// backing up to remote https://play.minio.io:9000:
	_ = os.Setenv("ACCESS_KEY_ID", "Q3AM3UQ867SPQQA43P2F")
	_ = os.Setenv("SECRET_ACCESS_KEY", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")
	etcd2Backup(t, port, tetcd, types.LegacyKubernetesPrefix+"/namespaces/kube-system")
}

func TestFilters(t *testing.T) {
	for _, ft := range filtertests {
		_ = os.Setenv("RS_BACKUP_STRATEGY", ft.filter)
		_ = filter(ft.path, ".", "/tmp")
		got := true
		if _, err := os.Stat("/tmp" + ft.path); os.IsNotExist(err) {
			got = false
		}
		want := ft.exists
		if got != want {
			t.Errorf("backup.filter(\"%s\", \".\", \"/tmp\") with %s => %t, want %t", ft.path, ft.filter, got, want)
		}
	}
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

func etc3BackupWithLargeSubkeys(t *testing.T, port, tetcd, key string) {
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
	stringSize := int(0.5 * 1024 * 1024)
	for i := 0; i < 10; i++ {
		subkey := fmt.Sprintf("%s/%d", key, i)
		_, err = c3.Put(context.Background(), subkey, generateRandomString(stringSize))

		if err != nil {
			t.Errorf("Can't create etcd entry %s: %s", subkey, err)
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

// src: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func generateRandomString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
