package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	operationSuccess = "success"
	operationFail    = "fail"
)

var (
	backupTotal  *prometheus.CounterVec
	keysRestored prometheus.Counter
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
}
