// ReShifter enables backing up and restoring OpenShift clusters.
// The ReShifter app launches an API and a UI at port 8080.
// The API is instrumented, exposing Prometheus metrics.
// When launching the app with the defaults, the backups are created in the
// current directory and the temporary work files are placed in the /tmp directory.
package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/mhausenblas/reshifter/pkg/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := "8080"
	host, _ := util.ExternalIP()
	if envd := os.Getenv("DEBUG"); envd != "" {
		log.SetLevel(log.DebugLevel)
	}
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/v1/version", versionHandler)
	r.HandleFunc("/v1/explorer", explorerHandler)
	r.HandleFunc("/v1/epstats", keystatsHandler)
	r.HandleFunc("/v1/backup", backupCreateHandler).Methods("POST")
	r.HandleFunc("/v1/backup/{backupid:[0-9]+}", backupRetrieveHandler).Methods("GET")
	r.HandleFunc("/v1/restore", restoreHandler)
	log.Printf("Serving API from: %s:%s/v1", host, port)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/")))
	log.Printf("Serving UI from: %s:%s/", host, port)
	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:" + port,
	}
	log.Fatal(srv.ListenAndServe())
}
