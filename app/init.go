package main

import (
	"os"

	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	operationSuccess = "success"
	operationFail    = "fail"
)

var (
	releaseVersion string
	backupTotal    *prometheus.CounterVec
	keysRestored   prometheus.Counter
)

func init() {
	backupTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reshifter",
			Subsystem: "api",
			Name:      "backup_total",
			Help:      "Number of completed backup operations.",
		},
		[]string{"outcome"},
	)
	keysRestored = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "reshifter",
			Subsystem: "api",
			Name:      "keys_restored",
			Help:      "Number of keys restored.",
		})
	prometheus.MustRegister(backupTotal)
	prometheus.MustRegister(keysRestored)

	// make sure that our working directory exists:
	if _, err := os.Stat(types.DefaultWorkDir); os.IsNotExist(err) {
		_ = os.MkdirAll(types.DefaultWorkDir, 0777)
	}
}
