package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/pierrre/archivefile/zip"
)

// raw is a reap strategy that stores val in a path in directory target.
func raw(path, val string, target interface{}) error {
	t, ok := target.(string)
	if !ok {
		return fmt.Errorf("Can't use target %v, it should be a string!", target)
	}
	_, err := store(t, path, val)
	if err != nil {
		return fmt.Errorf("Can't store %s%s in %s: %s", t, path, val, err)
	}
	return nil
}

// render is a reap strategy that writes path/val to a writer.
func render(path, val string, writer interface{}) error {
	w, ok := writer.(io.Writer)
	if !ok {
		return fmt.Errorf("Can't use writer %v, it should be an io.Writer!", writer)
	}
	_, err := w.Write([]byte(fmt.Sprintf("%s = %s\n", path, val)))
	if err != nil {
		return fmt.Errorf("Can't write %s%s to %v: %s", path, val, w, err)
	}
	return nil
}

// filter is a reap strategy that stores val in a path in directory target,
// if a certain condition is  met (set of white-listed paths).
func filter(path, val string, target interface{}) error {
	t, ok := target.(string)
	if !ok {
		return fmt.Errorf("Can't use target %v, it should be a string!", target)
	}
	log.WithFields(log.Fields{"func": "backup.filter"}).Debug(fmt.Sprintf("On path %s", path))
	if strings.Contains(path, "deployments") {
		_, err := store(t, path, val)
		if err != nil {
			return fmt.Errorf("Can't store %s%s in %s: %s", t, path, val, err)
		}
	}
	return nil
}

// store creates a file at based+path with val as its content.
// based is the output directory to use and path can be
// any valid etcd key (with ':'' characters being escaped automatically).
func store(based string, path string, val string) (string, error) {
	log.WithFields(log.Fields{"func": "backup.store"}).Debug(fmt.Sprintf("Trying to store %s with value=%s in %s", path, val, based))
	// make sure we're dealing with a valid path
	// that is, non-empty and has to start with /:
	if path == "" || (strings.Index(path, "/") != 0) {
		return "", fmt.Errorf("Path has to be non-empty")
	}
	// escape ":" in the path so that we have no issues storing it in the filesystem:
	fpath, _ := filepath.Abs(filepath.Join(based, strings.Replace(path, ":", types.EscapeColon, -1)))
	if path == "/" {
		log.WithFields(log.Fields{"func": "backup.store"}).Debug(fmt.Sprintf("Rewriting root"))
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
	log.WithFields(log.Fields{"func": "backup.store"}).Debug(fmt.Sprintf("Stored %s in %s with %d bytes", path, fpath, nbytes))
	return cpath, nil
}

// arch creates a ZIP archive of the content store() has generated
func arch(based string) (string, error) {
	defer func() {
		_ = os.RemoveAll(based)
	}()
	log.WithFields(log.Fields{"func": "backup.arch"}).Debug(fmt.Sprintf("Trying to pack backup into %s.zip", based))
	opath := based + ".zip"
	err := zip.ArchiveFile(based, opath, func(apath string) {
		log.WithFields(log.Fields{"func": "backup.arch"}).Debug(fmt.Sprintf("%s", apath))
	})
	if err != nil {
		return "", fmt.Errorf("Can't create archive or no content to back up: %s", err)
	}
	return opath, nil
}
