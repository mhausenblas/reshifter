package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/mhausenblas/reshifter/pkg/backup"
	"github.com/mhausenblas/reshifter/pkg/discovery"
	"github.com/mhausenblas/reshifter/pkg/restore"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
)

// versionHandler responds to HTTP GET requests of the form:
//		http GET localhost:8080/v1/version
func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ReShifter in version %s", releaseVersion)
}

// epstatsHandler responds to HTTP GET requests such as:
//		http GET localhost:8080/v1/epstats?endpoint=http%3A%2F%2Flocalhost%3A2379
func epstatsHandler(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Query().Get("endpoint")
	if endpoint == "" || strings.Index(endpoint, "http") != 0 {
		merr := "The endpoint is malformed"
		http.Error(w, merr, http.StatusBadRequest)
		log.Error(merr)
		return
	}
	vk, vs, err := discovery.CountKeysFor(endpoint, types.Vanilla)
	if err != nil {
		merr := fmt.Sprintf("Having problems calculating stats: %s", err)
		http.Error(w, merr, http.StatusInternalServerError)
		log.Error(merr)
		return
	}
	log.Debugf("vanilla [keys:%d, size:%d]", vk, vs)
	// note: ignoring error here since we're adding up the stats
	// and if this happens to be a non-OpenShift distro we simply
	// add 0 to the overall count, and it's still fine:
	osk, oss, _ := discovery.CountKeysFor(endpoint, types.OpenShift)
	log.Debugf("openshift [keys:%d, size:%d]", osk, oss)
	_ = json.NewEncoder(w).Encode(struct {
		NumKeys         int `json:"numkeys"`
		TotalSizeValues int `json:"totalsizevalbytes"`
	}{
		vk + osk,
		vs + oss,
	})
}

// explorerHandler responds to HTTP GET requests such as:
//		http GET localhost:8080/v1/explorer?endpoint=http%3A%2F%2Flocalhost%3A2379
func explorerHandler(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Query().Get("endpoint")
	if endpoint == "" || strings.Index(endpoint, "http") != 0 {
		merr := "The endpoint is malformed"
		http.Error(w, merr, http.StatusBadRequest)
		log.Error(merr)
		return
	}
	version, issecure, err := discovery.ProbeEtcd(endpoint)
	if err != nil {
		merr := fmt.Sprintf("%s", err)
		http.Error(w, merr, http.StatusBadRequest)
		log.Error(err)
		return
	}
	distrotype, err := discovery.ProbeKubernetesDistro(endpoint)
	if err != nil {
		merr := fmt.Sprintf("Can't determine Kubernetes distro: %s", err)
		http.Error(w, merr, http.StatusBadRequest)
		log.Error(err)
		return
	}
	secure := "insecure etcd, no SSL/TLS configured"
	if issecure {
		secure = "secure etcd, SSL/TLS configure"
	}
	_ = json.NewEncoder(w).Encode(struct {
		EtcdVersion  string `json:"etcdversion"`
		EtcdSecurity string `json:"etcdsec"`
		K8SDistro    string `json:"k8sdistro"`
	}{
		version,
		secure,
		util.LookupDistro(distrotype),
	})
}

// backupCreateHandler responds to HTTP POST requests such as:
//		http GET localhost:8080/v1/backup endpoint=http://localhost:2379
func backupCreateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var breq types.BackupRequest
	err := decoder.Decode(&breq)
	if err != nil {
		mreq := "The backup request is malformed"
		http.Error(w, mreq, http.StatusBadRequest)
		log.Error(mreq)
		return
	}
	log.Infof("Starting backup from %s", breq.Endpoint)
	w.Header().Set("Content-Type", "application/json")
	bres := types.BackupResult{
		Outcome:  operationSuccess,
		BackupID: "0",
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

// backupRetrieveHandler responds to HTTP GET requests such as:
//		http GET localhost:8080/v1/backup/1498230556
func backupRetrieveHandler(w http.ResponseWriter, r *http.Request) {
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

// backupListHandler responds to HTTP GET requests at:
//		http GET localhost:8080/v1/backup/all?remote=play.minio.io:9000&bucket=test123
func backupListHandler(w http.ResponseWriter, r *http.Request) {
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

// / restoreUploadHandler responds to HTTP POST requests at:
//		http POST localhost:8080/v1/restore/upload
func restoreUploadHandler(w http.ResponseWriter, r *http.Request) {
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

// restoreHandler responds to HTTP POST requests such as:
//		http POST localhost:8080/v1/restore endpoint=http://localhost:2379 backupid=1498230556
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
	target := types.DefaultWorkDir
	if !util.IsBackupID(rreq.BackupID) {
		abortreason := fmt.Sprintf("Aborting restore: %s is not a valid backup ID", rreq.BackupID)
		http.Error(w, abortreason, http.StatusConflict)
		log.Error(abortreason)
		return
	}
	log.Infof("Starting restore to %s from backup %s", rreq.Endpoint, rreq.BackupID)
	w.Header().Set("Content-Type", "application/json")
	rr := types.RestoreResult{
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
