package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	coll := newSpaceCollector()

	registry := prometheus.NewRegistry()
	registry.MustRegister(coll)

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	http.Handle("/space", handler)

	log.Println("Listening on", "127.0.0.1:9999")

	http.ListenAndServe("127.0.0.1:9999", nil)
}

type spaceCollector struct {
	peopleInSpace prometheus.Gauge
}

func newSpaceCollector() *spaceCollector {
	return &spaceCollector{
		peopleInSpace: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "people_in_space",
			Help: "Number of people currently in space",
		}),
	}
}

func (sc *spaceCollector) Describe(ch chan<- *prometheus.Desc) {
	sc.peopleInSpace.Describe(ch)
}

func (sc *spaceCollector) Collect(ch chan<- prometheus.Metric) {
	err := sc.observe() // HL
	if err != nil {
		log.Println("Error observing data:", err)
		return
	}

	sc.peopleInSpace.Collect(ch)
}

type astroResponse struct{ Number int }

func (sc *spaceCollector) observe() error {
	resp, err := http.Get("http://api.open-notify.org/astros.json") // HL
	if err != nil {
		return fmt.Errorf("fetching data: %w", err)
	}
	defer resp.Body.Close()

	astros := astroResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&astros); err != nil { // HL
		return fmt.Errorf("decoding data: %w", err)
	}

	sc.peopleInSpace.Set(float64(astros.Number)) // HL

	return nil
}

// END OMIT
