package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/mhausenblas/reshifter/app/pkg"
)

// BurryConfig holds all relevant config parameters for burry to run
type BurryConfig struct {
	Target   string `json:"target"`
	Endpoint string `json:"endpoint"`
}

func main() {
	go ui()
	go api()
	select {}
}

func api() {
	bc := BurryConfig{
		Target:   "local",
		Endpoint: "localhost:2379",
	}
	http.HandleFunc("/v1/version", func(w http.ResponseWriter, r *http.Request) {
		version := "0.1"
		fmt.Fprintf(w, "ReShifter in version %s", version)
	})
	http.HandleFunc("/v1/backup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := etcd.Backup(bc.Endpoint)
		if err != nil {
			log.Error(err)
		}
		_ = json.NewEncoder(w).Encode(bc)
	})
	log.Println("Serving API from /v1")
	_ = http.ListenAndServe(":8080", nil)
}

func ui() {
	fs := http.FileServer(http.Dir("ui"))
	http.Handle("/", fs)
	log.Println("Serving UI from /")
	_ = http.ListenAndServe(":8080", nil)
}
