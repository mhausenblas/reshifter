package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/backup"
	"github.com/mhausenblas/reshifter/pkg/restore"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
)

func versionHandler(w http.ResponseWriter, r *http.Request) {
	version := "0.1.23"
	fmt.Fprintf(w, "ReShifter in version %s", version)
}

// backupHandler responds to HTTP POST requests such as:
//		http POST localhost:8080/v1/backup endpoint=http://localhost:2379
func backupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only HTTP POST is supported", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var breq types.BackupRequest
	err := decoder.Decode(&breq)
	if err != nil {
		http.Error(w, "The backup request is malformed", http.StatusBadRequest)
		return
	}
	log.Infof("Starting backup from %s", breq.Endpoint)
	w.Header().Set("Content-Type", "application/json")
	bres := types.BackupResult{
		Outcome:  operationSuccess,
		BackupID: "0",
	}
	bid, err := backup.Backup(breq.Endpoint)
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

// restoreHandler responds to HTTP POST requests such as:
//		http POST localhost:8080/v1/restore endpoint=http://localhost:2379 archive=1498230556
func restoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only HTTP POST is supported", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var rreq types.RestoreRequest
	err := decoder.Decode(&rreq)
	if err != nil {
		http.Error(w, "The restore request is malformed", http.StatusBadRequest)
		return
	}
	target := "/tmp"
	if !util.IsBackupID(rreq.Archive) {
		abortreason := fmt.Sprintf("Aborting restore: %s is not a valid backup ID", rreq.Archive)
		http.Error(w, abortreason, http.StatusConflict)
		log.Error(abortreason)
		return
	}
	log.Infof("Starting restore to %s from backup %s", rreq.Endpoint, rreq.Archive)
	w.Header().Set("Content-Type", "application/json")
	rr := types.RestoreResult{
		Outcome:      operationSuccess,
		KeysRestored: 0,
	}
	krestored, err := restore.Restore(rreq.Archive, target, rreq.Endpoint)
	if err != nil {
		rr.Outcome = operationFail
		log.Error(err)
	}
	rr.KeysRestored = krestored
	keysRestored.Add(float64(krestored))
	log.Infof("Completed restore from %s to %s: %v", rreq.Archive, rreq.Endpoint, rr)
	if rr.Outcome == "fail" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	_ = json.NewEncoder(w).Encode(rr)
}
