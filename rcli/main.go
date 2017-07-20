// Package main implements the ReShifter CLI tool rcli.
package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/rcli/cmd"
)

func main() {
	if envd := os.Getenv("DEBUG"); envd != "" {
		log.SetLevel(log.DebugLevel)
	}
	cmd.Execute()
}
