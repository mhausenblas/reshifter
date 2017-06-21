package etcd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/pierrre/archivefile/zip"
	"golang.org/x/net/context"
)

// Backup traverses all paths of an etcd server starting from the root
// and creates a ZIP archive of the content in the current directory
func Backup(endpoint string) (string, error) {
	based := fmt.Sprintf("%d", time.Now().Unix())
	cfg := client.Config{
		Endpoints:               []string{"http://" + endpoint},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.WithFields(log.Fields{"func": "backup"}).Error(fmt.Sprintf("Can't connect to etcd: %s", err))
		return "", fmt.Errorf("Can't connect to etcd: %s", err)
	}

	kapi := client.NewKeysAPI(c)
	visit(kapi, "/", func(path string, val string) {
		_, _ = store(based, path, val)
		// if err != nil {
		// 	return "", fmt.Errorf("Can't store value locally: %s", err)
		// }
	})
	_, err = arch(based)
	if err != nil {
		return "", err
	}
	return based, nil
}

// visit recursively visits a path in the etcd tree and applies the reap function
// on a node, if it is a leaf node, otherwise descents the tree
func visit(kapi client.KeysAPI, path string, fn reap) {
	log.WithFields(log.Fields{"func": "visit"}).Debug(fmt.Sprintf("On node %s", path))
	copts := client.GetOptions{
		Recursive: true,
		Sort:      false,
		Quorum:    true,
	}
	if resp, err := kapi.Get(context.Background(), path, &copts); err != nil {
		log.WithFields(log.Fields{"func": "visit"}).Error(fmt.Sprintf("%s", err))
	} else {
		if resp.Node.Dir { // there are children
			log.WithFields(log.Fields{"func": "visit"}).Debug(fmt.Sprintf("%s has %d children", path, len(resp.Node.Nodes)))
			for _, node := range resp.Node.Nodes {
				log.WithFields(log.Fields{"func": "visit"}).Debug(fmt.Sprintf("Next visiting child %s", node.Key))
				visit(kapi, node.Key, fn)
			}
		} else { // we're on a leaf node
			fn(resp.Node.Key, string(resp.Node.Value))
		}
	}
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
	fpath, _ := filepath.Abs(filepath.Join(cwd, based, strings.Replace(path, ":", EscapeColon, -1)))
	if path == "/" {
		log.WithFields(log.Fields{"func": "store"}).Debug(fmt.Sprintf("Rewriting root"))
		fpath, _ = filepath.Abs(filepath.Join(cwd, based))
	}
	err := os.MkdirAll(fpath, os.ModePerm)
	if err != nil {
		log.WithFields(log.Fields{"func": "store"}).Error(fmt.Sprintf("%s", err))
		return "", fmt.Errorf("%s", err)
	}
	cpath, _ := filepath.Abs(filepath.Join(fpath, ContentFile))
	c, err := os.Create(cpath)
	if err != nil {
		log.WithFields(log.Fields{"func": "store"}).Error(fmt.Sprintf("%s", err))
		return "", fmt.Errorf("%s", err)
	}
	defer func() {
		_ = c.Close()
	}()
	nbytes, err := c.WriteString(val)
	if err != nil {
		log.WithFields(log.Fields{"func": "store"}).Error(fmt.Sprintf("%s", err))
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
	progress := func(apath string) {
		log.WithFields(log.Fields{"func": "arch"}).Debug(fmt.Sprintf("%s", apath))
	}
	err := zip.ArchiveFile(ipath, opath, progress)
	if err != nil {
		log.WithFields(log.Fields{"func": "arch"}).Error(fmt.Sprintf("Can't create archive: %s", err))
		return "", fmt.Errorf("Can't create archive: %s", err)
	}
	return opath, nil
}
