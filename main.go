// ReShifter enables backing up and restoring OpenShift clusters.
// The ReShifter app launches an API and a UI at port 8080.
// The API is instrumented, exposing Prometheus metrics.
// When launching the app with the defaults, the backups are created in the
// current directory and the temporary work files are placed in the /tmp directory.
package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
	log.Println("Serving API from :8080/v1")
	_ = http.ListenAndServe(":8080", nil)
}

func ui() {
	fs := http.FileServer(http.Dir("ui"))
	http.Handle("/", fs)
	log.Println("Serving UI from :8080/")
	_ = http.ListenAndServe(":8080", nil)
}
