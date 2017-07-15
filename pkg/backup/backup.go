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
	version, secure, err := discovery.ProbeEtcd(endpoint)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	distrotype, err := discovery.ProbeKubernetesDistro(endpoint)
	if err != nil {
		return "", fmt.Errorf("Can't determine Kubernetes distro: %s", err)
	}

	strategyName, strategy := pickStrategy()

	// deal with etcd3 servers:
	if strings.HasPrefix(version, "3") {
		c3, cerr := util.NewClient3(endpoint, secure)
		if cerr != nil {
			return "", fmt.Errorf("Can't connect to etcd3: %s", cerr)
		}
		defer func() { _ = c3.Close() }()
		log.WithFields(log.Fields{"func": "backup.Backup"}).Debug(fmt.Sprintf("Got etcd3 cluster with endpoints %v", c3.Endpoints()))
		kprefix := types.LegacyKubernetesPrefix
		_, gerr := c3.Get(context.Background(), types.KubernetesPrefix)
		if gerr == nil { // key found
			kprefix = types.KubernetesPrefix
		}
		err = discovery.Visit3(c3, kprefix, target, types.Vanilla, strategy, strategyName)
		if err != nil {
			return "", err
		}
		if distrotype == types.OpenShift {
			err = discovery.Visit3(c3, types.OpenShiftPrefix, target, types.OpenShift, strategy, strategyName)
			if err != nil {
				return "", err
			}
		}
	}
	// deal with etcd2 servers:
	if strings.HasPrefix(version, "2") {
		c2, cerr := util.NewClient2(endpoint, secure)
		if cerr != nil {
			return "", fmt.Errorf("Can't connect to etcd2: %s", cerr)
		}
		kapi := client.NewKeysAPI(c2)
		log.WithFields(log.Fields{"func": "backup.Backup"}).Debug(fmt.Sprintf("Got etcd2 cluster with %v", c2.Endpoints()))
		kprefix := types.LegacyKubernetesPrefix
		_, gerr := kapi.Get(context.Background(), types.KubernetesPrefix, nil)
		if gerr == nil { // key found
			kprefix = types.KubernetesPrefix
		}
		err = discovery.Visit2(kapi, kprefix, target, strategy, strategyName)
		if err != nil {
			return "", err
		}
		if distrotype == types.OpenShift {
			err = discovery.Visit2(kapi, types.OpenShiftPrefix, target, strategy, strategyName)
			if err != nil {
				return "", err
			}
		}
	}

	if strategyName == types.ReapFunctionRaw {
		// create ZIP file of the reaped content:
		_, err = arch(target)
		if err != nil {
			return "", err
		}
		if remote != "" {
			err = remotes.StoreInS3(remote, bucket, target, based)
			if err != nil {
				return "", err
			}
		}
	}

	return based, nil
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

func pickStrategy() (string, types.Reap) {
	backupstrategy := os.Getenv("RS_BACKUP_STRATEGY")
	switch backupstrategy {
	case "raw":
		return types.ReapFunctionRaw, raw
	case "render":
		return types.ReapFunctionRender, render
	default:
		return types.ReapFunctionRaw, raw
	}
}
