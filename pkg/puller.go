package pkg

import (
	"github.com/prometheus/client_golang/prometheus"
)

func InitPuller(cs ...prometheus.Collector) {
	collectors := []prometheus.Collector{}
	collectors = append(collectors, cs...)
	prometheus.MustRegister(collectors...)
}
