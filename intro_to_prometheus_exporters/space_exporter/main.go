package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	coll := newSpaceCollector()

	registry := prometheus.NewRegistry() // HL
	registry.MustRegister(coll)          // HL

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{}) // HL

	http.Handle("/space", handler)

	log.Println("Listening on", "127.0.0.1:9999")

	http.ListenAndServe("127.0.0.1:9999", nil)
}

type spaceCollector struct {
	peopleInSpace prometheus.Gauge
}

func newSpaceCollector() *spaceCollector {
	return &spaceCollector{
		peopleInSpace: prometheus.NewGauge(prometheus.GaugeOpts{ // HL
			Name: "people_in_space",                     // HL
			Help: "Number of people currently in space", // HL
		}), // HL
	}
}

func (sc *spaceCollector) Describe(ch chan<- *prometheus.Desc) {
	sc.peopleInSpace.Describe(ch) // HL
}

func (sc *spaceCollector) Collect(ch chan<- prometheus.Metric) {
	// TODO: get the number of people in space

	sc.peopleInSpace.Collect(ch) // HL
}

// END OMIT
