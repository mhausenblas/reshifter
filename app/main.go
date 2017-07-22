// The ReShifter app launches an API and a UI at port 8080 and uses the
// ReShifter library to communicate and manipulate etcd. The API is
// instrumented, exposing Prometheus metrics. When launching the app with
// the defaults, the backups are created in /tmp/reshifter and via environment
// variables you can define the etcd certs/keys as well as the S3 backend
// credentials, see also http://reshifter.info/#configuration for more infos.
package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/mhausenblas/reshifter/app/handler"
	"github.com/mhausenblas/reshifter/pkg/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if envd := os.Getenv("DEBUG"); envd != "" {
		log.SetLevel(log.DebugLevel)
	}
	port := "8080"
	host, _ := util.ExternalIP()
	r := mux.NewRouter()
	// the HTTP API:
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/v1/version", handler.Version).Methods("GET")
	r.HandleFunc("/v1/explorer", handler.Explorer).Methods("GET")
	r.HandleFunc("/v1/epstats", handler.EPstats).Methods("GET")
	r.HandleFunc("/v1/backup", handler.BackupCreate).Methods("POST")
	r.HandleFunc("/v1/backup/all", handler.BackupList).Methods("GET")
	r.HandleFunc("/v1/backup/{backupid:[0-9]+}", handler.BackupRetrieve).Methods("GET")
	r.HandleFunc("/v1/restore", handler.Restore).Methods("POST")
	r.HandleFunc("/v1/restore/upload", handler.RestoreUpload).Methods("POST")
	log.Printf("Serving API from: %s:%s/v1", host, port)
	// the Web UI:
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/")))
	log.Printf("Serving UI from: %s:%s/", host, port)
	http.Handle("/", r)
	// the app server:
	srv := &http.Server{Handler: r, Addr: "0.0.0.0:" + port}
	log.Fatal(srv.ListenAndServe())
}
