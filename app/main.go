package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		Endpoint: "",
	}
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		version := "0.1"
		fmt.Fprintf(w, "ReShifter in version %s", version)
	})
	http.HandleFunc("/v1/backup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(bc)
	})
	log.Println("Serving API")
	_ = http.ListenAndServe(":8080", nil)
}

func ui() {
	fs := http.FileServer(http.Dir("ui"))
	http.Handle("/", fs)
	log.Println("Serving UI")
	_ = http.ListenAndServe(":8080", nil)
}

func backup() {
	var url = ""
	var accessToken = ""
	var bearer = fmt.Sprintf("Bearer %s", accessToken)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authorization", bearer)

	// TBD: fetch content

	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	//
	// body, _ := ioutil.ReadAll(resp.Body)
}
