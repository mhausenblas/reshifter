package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/etcd"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	operationSuccess = "success"
	operationFail    = "fail"
)

var (
	backupTotal *prometheus.CounterVec
)

func init() {
	backupTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "dev",
			Subsystem: "app_server",
			Name:      "backup_total",
			Help:      "The count of backup attempts.",
		},
		[]string{"outcome"},
	)
	prometheus.MustRegister(backupTotal)

	if envd := os.Getenv("DEBUG"); envd != "" {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	go api()
	go ui()
	select {}
}

func api() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/v1/version", versionHandler)
	http.HandleFunc("/v1/backup", backupHandler)
	http.HandleFunc("/v1/restore", restoreHandler)
	log.Println("Serving API from /v1")
	_ = http.ListenAndServe(":8080", nil)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	version := "0.1.16"
	fmt.Fprintf(w, "ReShifter in version %s", version)
}

func backupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ep := etcd.Endpoint{
		Version: "2",
		URL:     "localhost:2379",
	}
	br := etcd.BackupResult{
		Outcome:  operationSuccess,
		BackupID: "0",
	}
	bid, err := etcd.Backup(ep.URL)
	if err != nil {
		br.Outcome = operationFail
		log.Error(err)
	}
	br.BackupID = bid
	backupTotal.WithLabelValues(br.Outcome).Inc()
	log.Infof("Completed backup from %s: %v", ep.URL, br)
	if br.Outcome == operationFail {
		http.Error(w, err.Error(), 409)
		return
	}
	_ = json.NewEncoder(w).Encode(br)
}

func restoreHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ep := etcd.Endpoint{
		Version: "2",
		URL:     "localhost:2379",
	}
	rr := etcd.RestoreResult{
		Outcome:      operationSuccess,
		KeysRestored: 0,
	}
	target := "/tmp"
	afile := r.URL.Query().Get("archive")
	if !etcd.IsBackupID(afile) {
		abortreason := fmt.Sprintf("Aborting restore: %s is not a valid backup ID", afile)
		http.Error(w, abortreason, 409)
		log.Error(abortreason)
		return
	}
	krestored, err := etcd.Restore(afile, target, ep.URL)
	if err != nil {
		rr.Outcome = operationFail
		log.Error(err)
	}
	rr.KeysRestored = krestored
	log.Infof("Completed restore from %s to %s: %v", afile, ep.URL, rr)
	if rr.Outcome == "fail" {
		http.Error(w, err.Error(), 409)
		return
	}
	_ = json.NewEncoder(w).Encode(rr)
}

func ui() {
	fs := http.FileServer(http.Dir("ui"))
	http.Handle("/", fs)
	log.Println("Serving UI from /")
	_ = http.ListenAndServe(":8080", nil)
}
