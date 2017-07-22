package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/mhausenblas/reshifter/pkg/backup"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
)

// BackupCreate responds to HTTP POST requests such as:
//		http GET localhost:8080/v1/backup endpoint=http://localhost:2379
func BackupCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var breq BackupRequest
	err := decoder.Decode(&breq)
	if err != nil {
		mreq := "The backup request is malformed"
		http.Error(w, mreq, http.StatusBadRequest)
		log.Error(mreq)
		return
	}
	log.Infof("Starting backup from %s", breq.Endpoint)
	w.Header().Set("Content-Type", "application/json")
	bres := BackupResult{
		Outcome:  operationSuccess,
		BackupID: "0",
	}
	fmt.Printf("%v\n", breq)
	if breq.Filter != "" {
		_ = os.Setenv("RS_BACKUP_STRATEGY", fmt.Sprintf("filter:%s", breq.Filter))
		log.Infof("Using filter backup strategy with following whitelist: %s", breq.Filter)
	}
	if breq.APIversion != "" {
		_ = os.Setenv("RS_ETCD_API_VERSION", fmt.Sprintf("%s", breq.APIversion))
		log.Infof("Using etcd API version: %s", breq.APIversion)
	}
	target := types.DefaultWorkDir
	bid, err := backup.Backup(breq.Endpoint, target, breq.Remote, breq.Bucket)
	if err != nil {
		bres.Outcome = operationFail
		log.Error(err)
	}
	bres.BackupID = bid
	backupTotal.WithLabelValues(bres.Outcome).Inc()
	log.Infof("Completed backup from %s: %v", breq.Endpoint, bres)
	if bres.Outcome == operationFail {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	_ = json.NewEncoder(w).Encode(bres)
}

// BackupRetrieve responds to HTTP GET requests such as:
//		http GET localhost:8080/v1/backup/1498230556
func BackupRetrieve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	backupid := vars["backupid"]
	if !util.IsBackupID(backupid) {
		abortreason := fmt.Sprintf("Aborting backup retrieve: %s is not a valid backup ID", backupid)
		http.Error(w, abortreason, http.StatusConflict)
		log.Error(abortreason)
		return
	}
	target := types.DefaultWorkDir
	c, err := ioutil.ReadFile(filepath.Join(target, backupid) + ".zip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	fmt.Fprint(w, string(c))
}

// BackupList responds to HTTP GET requests at:
//		http GET localhost:8080/v1/backup/all?remote=play.minio.io:9000&bucket=test123
func BackupList(w http.ResponseWriter, r *http.Request) {
	remote := r.URL.Query().Get("remote")
	bucket := r.URL.Query().Get("bucket")
	backupIDs, err := backup.List(remote, bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	_ = json.NewEncoder(w).Encode(struct {
		NumBackups   int      `json:"numbackups"`
		EtcdSecurity []string `json:"backupids"`
	}{
		len(backupIDs),
		backupIDs,
	})
}
