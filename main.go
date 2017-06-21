package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/etcd"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
}

func main() {
	go ui()
	go api()
	select {}
}

func api() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/v1/version", versionHandler)
	http.HandleFunc("/v1/backup", backupHandler)
	log.Println("Serving API from /v1")
	_ = http.ListenAndServe(":8080", nil)
}

func backupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ep := etcd.Endpoint{
		Version: "2",
		URL:     "localhost:2379",
	}
	outcome := "success"
	b, err := etcd.Backup(ep.URL)
	if err != nil {
		outcome = "failed"
		log.Error(err)
	}
	log.Println("Created backup from %s in %s", ep.URL, b)
	_ = json.NewEncoder(w).Encode(ep)
	backupTotal.WithLabelValues(outcome).Inc()
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	version := "0.1"
	fmt.Fprintf(w, "ReShifter in version %s", version)
}

func ui() {
	fs := http.FileServer(http.Dir("ui"))
	http.Handle("/", fs)
	log.Println("Serving UI from /")
	_ = http.ListenAndServe(":8080", nil)
}
