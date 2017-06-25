package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	"github.com/mhausenblas/reshifter/pkg/discovery"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
	"github.com/pierrre/archivefile/zip"
	"golang.org/x/net/context"
)

// Backup traverses all paths of an etcd server starting from the root
// and creates a ZIP archive of the content in the current directory.
// On success, it returns the backup ID, which is the Unix time encoded
// point in time the backup operation was started, for example 1498050161.
// Example:
//
//		bID, err := etcd.Backup("http://localhost:2379")
func Backup(endpoint string) (string, error) {
	// if envd := os.Getenv("DEBUG"); envd != "" {
	// 	log.SetLevel(log.DebugLevel)
	// }
	based := fmt.Sprintf("%d", time.Now().Unix())
	version, secure, err := discovery.ProbeEtcd(endpoint)
	if err != nil {
		return "", fmt.Errorf("Can't understand endpoint %s: %s", endpoint, err)
	}
	// deal with etcd3 servers:
	if strings.HasPrefix(version, "3") {
		c3, cerr := util.NewClient3(endpoint, secure)
		if cerr != nil {
			return "", fmt.Errorf("Can't connect to etcd: %s", cerr)
		}
		defer func() { _ = c3.Close() }()
		log.WithFields(log.Fields{"func": "Backup"}).Debug(fmt.Sprintf("Got etcd3 cluster with %v", c3.Endpoints()))
		err = visit3(c3, types.KubernetesPrefix, func(path string, val string) error {
			_, err = store(based, path, val)
			if err != nil {
				return fmt.Errorf("Can't store backup locally: %s", err)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}
	// deal with etcd2 servers:
	if strings.HasPrefix(version, "2") {
		c2, cerr := util.NewClient2(endpoint, secure)
		if cerr != nil {
			return "", fmt.Errorf("Can't connect to etcd: %s", cerr)
		}
		kapi := client.NewKeysAPI(c2)
		log.WithFields(log.Fields{"func": "Backup"}).Debug(fmt.Sprintf("Got etcd2 cluster with %v", c2.Endpoints()))

		err = visit2(kapi, "/", func(path string, val string) error {
			_, err = store(based, path, val)
			if err != nil {
				return fmt.Errorf("Can't store backup locally: %s", err)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}
	// create ZIP file of the reaped content:
	_, err = arch(based)
	if err != nil {
		return "", err
	}
	return based, nil
}

// visit2 recursively visits an etcd2 server from the root and applies
// the reap function on leaf nodes (keys that don't have sub-keys),
// otherwise descents the tree.
func visit2(kapi client.KeysAPI, path string, fn types.Reap) error {
	log.WithFields(log.Fields{"func": "backup.visit2"}).Debug(fmt.Sprintf("On node %s", path))
	copts := client.GetOptions{
		Recursive: true,
		Sort:      false,
		Quorum:    true,
	}
	res, err := kapi.Get(context.Background(), path, &copts)
	if err != nil {
		return err
	}
	if res.Node.Dir { // there are children
		log.WithFields(log.Fields{"func": "backup.visit2"}).Debug(fmt.Sprintf("%s has %d children", path, len(res.Node.Nodes)))
		for _, node := range res.Node.Nodes {
			log.WithFields(log.Fields{"func": "backup.visit2"}).Debug(fmt.Sprintf("Next visiting child %s", node.Key))
			_ = visit2(kapi, node.Key, fn)
		}
		return nil
	}
	// otherwise we're on a leaf node:
	return fn(res.Node.Key, string(res.Node.Value))
}

// visit3 visits all paths of an etcd3 server and applies the reap function
// on the keys.
func visit3(c3 *clientv3.Client, path string, fn types.Reap) error {
	log.WithFields(log.Fields{"func": "backup.visit3"}).Debug(fmt.Sprintf("Processing %s", path))
	res, err := c3.Get(context.Background(), types.KubernetesPrefix+"*", clientv3.WithRange(types.KubernetesPrefixLast))
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"func": "backup.visit3"}).Debug(fmt.Sprintf("Got %v", res.Kvs))
	for _, ev := range res.Kvs {
		log.WithFields(log.Fields{"func": "backup.visit3"}).Debug(fmt.Sprintf("key: %s, value: %s", ev.Key, ev.Value))
		err = fn(string(ev.Key), string(ev.Value))
		if err != nil {
			return err
		}
	}
	return nil
}

// store creates a file at based+path with val as its content.
// based is the relative base directory to use and path can be
// any valid etcd key (with : characters being escaped automatically).
func store(based string, path string, val string) (string, error) {
	// make sure we're dealing with a valid path
	// that is, non-empty and has to start with /:
	if path == "" || (strings.Index(path, "/") != 0) {
		return "", fmt.Errorf("Path has to be non-empty")
	}
	cwd, _ := os.Getwd()
	// escape ":" in the path so that we have no issues storing it in the filesystem:
	fpath, _ := filepath.Abs(filepath.Join(cwd, based, strings.Replace(path, ":", types.EscapeColon, -1)))
	if path == "/" {
		log.WithFields(log.Fields{"func": "store"}).Debug(fmt.Sprintf("Rewriting root"))
		fpath, _ = filepath.Abs(filepath.Join(cwd, based))
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
	cwd, _ := os.Getwd()
	opath := filepath.Join(cwd, based+".zip")
	ipath := filepath.Join(cwd, based, "/")
	err := zip.ArchiveFile(ipath, opath, func(apath string) {
		log.WithFields(log.Fields{"func": "arch"}).Debug(fmt.Sprintf("%s", apath))
	})
	if err != nil {
		return "", fmt.Errorf("Can't create archive or no content to back up: %s", err)
	}
	return opath, nil
}
