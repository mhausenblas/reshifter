package etcd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/pierrre/archivefile/zip"
)

// Restore takes archive file afile (without file extension) and
// unpacks it into a target directory. It then traverses the target directory
// in the local filesystem and populates an etcd server at endpoint with the
// content of the sub-directories.
//		Restore("1498055655", "localhost:2379")
func Restore(afile, target string, endpoint string) error {
	err := unarch(afile+".zip", target)
	if err != nil {
		log.WithFields(log.Fields{"func": "Restore"}).Error(err)
		return err
	}
	c2, err := newClient2(endpoint, false)
	if err != nil {
		log.WithFields(log.Fields{"func": "Restore"}).Error(fmt.Sprintf("Can't connect to etcd: %s", err))
		return fmt.Errorf("Can't connect to etcd: %s", err)
	}
	kapi := client.NewKeysAPI(c2)
	numrestored := 0
	err = filepath.Walk(target, func(path string, f os.FileInfo, err error) error {
		if f.Name() == afile {
			return nil
		}
		base, _ := filepath.Abs(filepath.Join(target, afile))
		key, _ := filepath.Rel(base, path)
		// append the root "/" to make it a key and unescape ":"
		key = "/" + strings.Replace(key, EscapeColon, ":", -1)
		if f.IsDir() {
			cfile, _ := filepath.Abs(filepath.Join(path, ContentFile))
			_, err = os.Stat(cfile)
			if err != nil { // empty directory (no content file) inserting a non-leaf key
				_, err = kapi.Set(context.Background(), key, "", &client.SetOptions{Dir: true, PrevExist: client.PrevNoExist})
				if err != nil {
					log.WithFields(log.Fields{"func": "Restore"}).Error(fmt.Sprintf("Can't restore key %s: %s", key, err))
					return nil
				}
				log.WithFields(log.Fields{"func": "Restore"}).Info(fmt.Sprintf("Restored %s", key))
				numrestored++
				return nil
			}
			// there is a content file at this path, inserting a leaf key:
			c, err := ioutil.ReadFile(cfile)
			if err != nil {
				log.WithFields(log.Fields{"func": "Restore"}).Error(fmt.Sprintf("Can't read content file %s: %s", cfile, err))
				return nil
			}
			_, err = kapi.Set(context.Background(), key, string(c), &client.SetOptions{Dir: false, PrevExist: client.PrevNoExist})
			if err != nil {
				return nil
			}
			log.WithFields(log.Fields{"func": "Restore"}).Info(fmt.Sprintf("Restored %s", key))
			numrestored++
		}
		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{"func": "Restore"}).Error(fmt.Sprintf("Can't traverse directory %s: %s", target, err))
		fmt.Errorf("Can't traverse directory %s: %s", target, err)
	}
	return nil
}

// unarch takes archive file afile and unpacks it into a target directory.
func unarch(afile, target string) error {
	err := zip.UnarchiveFile(afile, target, func(apath string) {
		log.WithFields(log.Fields{"func": "unarch"}).Debug(fmt.Sprintf("%s", apath))
	})
	if err != nil {
		return fmt.Errorf("Can't unpack archive %s: %s", afile, err)
	}
	return nil
}
