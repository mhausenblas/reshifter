package handler

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// ReleaseVersion which is set by main
	ReleaseVersion string
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
}
