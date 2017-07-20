package backup

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	"github.com/mhausenblas/reshifter/pkg/discovery"
	"github.com/mhausenblas/reshifter/pkg/remotes"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
)

// Backup iterates over well-known Kubernetes (distro) keys in an etcd server
// and creates a ZIP archive of the content in the target directory.
// On success, it returns the backup ID, which is the Unix time encoded
// point in time the backup operation was started, for example 1498050161.
// If remote and bucket is provided, the backup will be additional stored
// in this S3-compatible object store.
func Backup(endpoint, target, remote, bucket string) (string, error) {
	based := fmt.Sprintf("%d", time.Now().Unix())
	if _, err := os.Stat(target); os.IsNotExist(err) {
		_ = os.Mkdir(target, 0700)
	}
	target, _ = filepath.Abs(filepath.Join(target, based))
	version, apiversion, secure, err := discovery.ProbeEtcd(endpoint)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	distrotype, err := discovery.ProbeKubernetesDistro(endpoint)
	if err != nil {
		return "", fmt.Errorf("Can't determine Kubernetes distro: %s", err)
	}
	switch {
	case strings.HasPrefix(version, "3"): // etcd3 server
		if apiversion == "v2" { // a v2 API in an etcd3
			err = backupv2(endpoint, target, secure, distrotype)
			if err != nil {
				return "", err
			}
			break
		}
		err = backupv3(endpoint, target, secure, distrotype)
		if err != nil {
			return "", err
		}
	case strings.HasPrefix(version, "2"): // etcd2 server
		err = backupv2(endpoint, target, secure, distrotype)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("Can't understand etcd version, seems to be neither v3 nor v2 :(")
	}
	strategyName, _ := pickStrategy()
	if strategyName != types.ReapFunctionRender {
		// create ZIP file of the reaped content:
		_, err = arch(target)
		if err != nil {
			return "", err
		}
		// store in S3-compatible remote storage, if requested:
		if remote != "" {
			err = remotes.StoreInS3(remote, bucket, target, based)
			if err != nil {
				return "", err
			}
		}
	}
	return based, nil
}

func backupv3(endpoint, target string, secure bool, distrotype types.KubernetesDistro) error {
	c3, err := util.NewClient3(endpoint, secure)
	if err != nil {
		return fmt.Errorf("Can't connect to etcd3: %s", err)
	}
	defer func() { _ = c3.Close() }()
	log.WithFields(log.Fields{"func": "backup.backupv3"}).Debug(fmt.Sprintf("Got etcd3 cluster with endpoints %v", c3.Endpoints()))
	kprefix := types.LegacyKubernetesPrefix
	klprefix := types.LegacyKubernetesPrefixLast
	kgv, _ := c3.Get(context.Background(), kprefix+"/*", clientv3.WithRange(klprefix))
	if kgv.Count == 0 { // legacy key not found, must be new key space
		kprefix = types.KubernetesPrefix
		klprefix = types.KubernetesPrefixLast
	}
	strategyName, strategy := pickStrategy()
	err = discovery.Visit3(c3, target, kprefix, klprefix, strategy, strategyName)
	if err != nil {
		return err
	}
	if distrotype == types.OpenShift {
		err = discovery.Visit3(c3, target, types.OpenShiftPrefix, types.OpenShiftPrefixLast, strategy, strategyName)
		if err != nil {
			return err
		}
	}
	return nil
}

func backupv2(endpoint, target string, secure bool, distrotype types.KubernetesDistro) error {
	c2, err := util.NewClient2(endpoint, secure)
	if err != nil {
		return fmt.Errorf("Can't connect to etcd2: %s", err)
	}
	log.WithFields(log.Fields{"func": "backup.checkv2"}).Debug(fmt.Sprintf("Got etcd2 cluster with %v", c2.Endpoints()))
	kapi := client.NewKeysAPI(c2)
	kprefix, err := checkv2(endpoint, secure)
	if err != nil {
		return err
	}
	if kprefix == "" {
		return fmt.Errorf("Can't find well-known v2 keyspaces")
	}
	strategyName, strategy := pickStrategy()
	err = discovery.Visit2(kapi, kprefix, target, strategy, strategyName)
	if err != nil {
		return err
	}
	if distrotype == types.OpenShift {
		err = discovery.Visit2(kapi, types.OpenShiftPrefix, target, strategy, strategyName)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkv2(endpoint string, secure bool) (string, error) {
	c2, err := util.NewClient2(endpoint, secure)
	if err != nil {
		return "", fmt.Errorf("Can't connect to etcd2: %s", err)
	}
	log.WithFields(log.Fields{"func": "backup.checkv2"}).Debug(fmt.Sprintf("Got etcd2 cluster with %v", c2.Endpoints()))
	kapi := client.NewKeysAPI(c2)
	kprefix := types.LegacyKubernetesPrefix
	kfound, _ := kapi.Get(context.Background(), kprefix, nil)
	if kfound != nil { // legacy v2 keyspace found
		return kprefix, nil
	}
	kprefix = types.KubernetesPrefix
	kfound, _ = kapi.Get(context.Background(), kprefix, nil)
	if kfound != nil { // modern v2 keyspace found
		return kprefix, nil
	}
	return "", nil
}

func pickStrategy() (string, types.Reap) {
	backupstrategy := os.Getenv("RS_BACKUP_STRATEGY")
	switch backupstrategy {
	case "raw":
		return types.ReapFunctionRaw, raw
	case "render":
		return types.ReapFunctionRender, render
	case "filter":
		return types.ReapFunctionFilter, filter
	default:
		return types.ReapFunctionRaw, raw
	}
}

// List generates a list of available backup IDs.
// If remote and bucket is provided, the S3-compatible object store
// will be queried rather than the local work directory.
func List(remote, bucket string) ([]string, error) {
	var backupIDs []string

	if remote == "" {
		files, err := ioutil.ReadDir(types.DefaultWorkDir)
		if err != nil {
			return nil, fmt.Errorf("Can't read backup IDs from local: %s", err)
		}
		for _, file := range files {
			re := regexp.MustCompile("\\d{10}.zip")
			fn := file.Name()
			bid := fn[0 : len(fn)-len(filepath.Ext(fn))]
			if re.Match([]byte(fn)) {
				backupIDs = append(backupIDs, bid)
			}
		}
		return backupIDs, nil
	}
	// we gonna query the remote S3-compatible object store:
	backupIDs, err := remotes.ListObjectsInS3Bucket(remote, bucket)
	if err != nil {
		return nil, fmt.Errorf("Can't read backup IDs from remote: %s", err)
	}
	return backupIDs, nil
}
