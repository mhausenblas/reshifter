package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/mhausenblas/reshifter/pkg/discovery"
	"github.com/mhausenblas/reshifter/pkg/remotes"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
	"github.com/pierrre/archivefile/zip"
)

// Backup iterates over well-known Kubernetes (distro) keys in an etcd server
// and creates a ZIP archive of the content in the target directory.
// On success, it returns the backup ID, which is the Unix time encoded
// point in time the backup operation was started, for example 1498050161.
// If remote and bucket is provided, the backup will be additional stored
// in this S3-compatible object store.
//
// Example:
//
//		bID, err := backup.Backup("http://localhost:2379", "/tmp", "play.minio.io:9000", "reshifter-test-cluster")
func Backup(endpoint, target, remote, bucket string) (string, error) {
	based := fmt.Sprintf("%d", time.Now().Unix())
	target, _ = filepath.Abs(filepath.Join(target, based))
	version, secure, err := discovery.ProbeEtcd(endpoint)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	distrotype, err := discovery.ProbeKubernetesDistro(endpoint)
	if err != nil {
		return "", fmt.Errorf("Can't determine Kubernetes distro: %s", err)
	}
	// deal with etcd3 servers:
	if strings.HasPrefix(version, "3") {
		c3, cerr := util.NewClient3(endpoint, secure)
		if cerr != nil {
			return "", fmt.Errorf("Can't connect to etcd3: %s", cerr)
		}
		defer func() { _ = c3.Close() }()
		log.WithFields(log.Fields{"func": "Backup"}).Debug(fmt.Sprintf("Got etcd3 cluster with endpoints %v", c3.Endpoints()))
		err = discovery.Visit3(c3, types.KubernetesPrefix, types.Vanilla, func(path string, val string) error {
			_, err = store(target, path, val)
			if err != nil {
				return fmt.Errorf("Can't store backup locally: %s", err)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
		if distrotype == types.OpenShift {
			err = discovery.Visit3(c3, types.OpenShiftPrefix, types.OpenShift, func(path string, val string) error {
				_, err = store(target, path, val)
				if err != nil {
					return fmt.Errorf("Can't store backup locally: %s", err)
				}
				return nil
			})
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
		log.WithFields(log.Fields{"func": "Backup"}).Debug(fmt.Sprintf("Got etcd2 cluster with %v", c2.Endpoints()))
		err = discovery.Visit2(kapi, types.KubernetesPrefix, func(path string, val string) error {
			_, err = store(target, path, val)
			if err != nil {
				return fmt.Errorf("Can't store backup locally: %s", err)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
		if distrotype == types.OpenShift {
			err = discovery.Visit2(kapi, types.OpenShiftPrefix, func(path string, val string) error {
				_, err = store(target, path, val)
				if err != nil {
					return fmt.Errorf("Can't store backup locally: %s", err)
				}
				return nil
			})
			if err != nil {
				return "", err
			}
		}
	}
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
	return based, nil
}

// store creates a file at based+path with val as its content.
// based is the output directory to use and path can be
// any valid etcd key (with ':'' characters being escaped automatically).
func store(based string, path string, val string) (string, error) {
	// make sure we're dealing with a valid path
	// that is, non-empty and has to start with /:
	if path == "" || (strings.Index(path, "/") != 0) {
		return "", fmt.Errorf("Path has to be non-empty")
	}
	// escape ":" in the path so that we have no issues storing it in the filesystem:
	fpath, _ := filepath.Abs(filepath.Join(based, strings.Replace(path, ":", types.EscapeColon, -1)))
	if path == "/" {
		log.WithFields(log.Fields{"func": "store"}).Debug(fmt.Sprintf("Rewriting root"))
		fpath = based
	}
	err := os.MkdirAll(fpath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	cpath, _ := filepath.Abs(filepath.Join(fpath, types.ContentFile))
	c, err := os.Create(cpath)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	defer func() {
		_ = c.Close()
	}()
	nbytes, err := c.WriteString(val)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	log.WithFields(log.Fields{"func": "store"}).Debug(fmt.Sprintf("Stored %s in %s with %d bytes", path, fpath, nbytes))
	return cpath, nil
}

// arch creates a ZIP archive of the content store() has generated
func arch(based string) (string, error) {
	defer func() {
		_ = os.RemoveAll(based)
	}()
	opath := based + ".zip"
	err := zip.ArchiveFile(based, opath, func(apath string) {
		log.WithFields(log.Fields{"func": "arch"}).Debug(fmt.Sprintf("%s", apath))
	})
	if err != nil {
		return "", fmt.Errorf("Can't create archive or no content to back up: %s", err)
	}
	return opath, nil
}
