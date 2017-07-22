package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/restore"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
)

// RestoreUpload responds to HTTP POST requests at:
//		http POST localhost:8080/v1/restore/upload
func RestoreUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(types.UploadInMemoryBufferSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	m := r.MultipartForm
	backupfiles := m.File["backupfile"]
	if len(backupfiles) == 0 {
		nbfe := fmt.Sprintf("No backup file uploaded. Aborting …")
		http.Error(w, nbfe, http.StatusInternalServerError)
		log.Error(nbfe)
		return
	}
	log.Infof("Got %v as backup file", backupfiles)
	var overallwritten int64
	for _, bf := range backupfiles {
		fn := bf.Filename
		bid := fn[0 : len(fn)-len(filepath.Ext(fn))]
		log.Infof("Verifying backup ID %s and then trying to upload content …", bid)
		if !util.IsBackupID(bid) {
			abortreason := fmt.Sprintf("Aborting upload: %s is not a valid backup ID. Must be a Unix timestamp formatted one such as 1499588813.zip …", bid)
			http.Error(w, abortreason, http.StatusConflict)
			log.Error(abortreason)
			return
		}
		src, err := bf.Open()
		defer func() { _ = src.Close() }()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		dst, err := os.Create(filepath.Join(types.DefaultWorkDir, bf.Filename))
		defer func() { _ = dst.Close() }()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		written, err := io.Copy(dst, src)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		overallwritten += written
	}
	log.Infof("Uploading backup file done, written %d bytes", overallwritten)
	_ = json.NewEncoder(w).Encode(struct {
		Received int64 `json:"received"`
	}{
		overallwritten,
	})
}

// Restore responds to HTTP POST requests such as:
//		http POST localhost:8080/v1/restore endpoint=http://localhost:2379 backupid=1498230556
func Restore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only HTTP POST is supported", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var rreq RestoreRequest
	err := decoder.Decode(&rreq)
	if err != nil {
		http.Error(w, "The restore request is malformed", http.StatusBadRequest)
		return
	}
	target := types.DefaultWorkDir
	if !util.IsBackupID(rreq.BackupID) {
		abortreason := fmt.Sprintf("Aborting restore: %s is not a valid backup ID", rreq.BackupID)
		http.Error(w, abortreason, http.StatusConflict)
		log.Error(abortreason)
		return
	}
	log.Infof("Starting restore to %s from backup %s", rreq.Endpoint, rreq.BackupID)
	w.Header().Set("Content-Type", "application/json")
	rr := RestoreResult{
		Outcome:      operationSuccess,
		KeysRestored: 0,
		ElapsedTime:  0,
	}
	krestored, etime, err := restore.Restore(rreq.Endpoint, rreq.BackupID, target, rreq.Remote, rreq.Bucket)
	if err != nil {
		rr.Outcome = operationFail
		log.Error(err)
	}
	rr.KeysRestored = krestored
	rr.ElapsedTime = etime.Seconds()
	keysRestored.Add(float64(krestored))
	log.Infof("Completed restore from %s to %s: %v", rreq.BackupID, rreq.Endpoint, rr)
	if rr.Outcome == "fail" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	_ = json.NewEncoder(w).Encode(rr)
}
