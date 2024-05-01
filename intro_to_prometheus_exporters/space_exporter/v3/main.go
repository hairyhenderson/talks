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
	peopleInSpace *prometheus.GaugeVec
}

func newSpaceCollector() *spaceCollector {
	return &spaceCollector{
		peopleInSpace: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "people_in_space",
			Help: "Number of people currently in space",
		}, []string{"craft"}), // HL
	}
}

func (sc *spaceCollector) Describe(ch chan<- *prometheus.Desc) {
	sc.peopleInSpace.Describe(ch)
}

func (sc *spaceCollector) Collect(ch chan<- prometheus.Metric) {
	err := sc.observe()
	if err != nil {
		log.Println("Error observing data:", err)
		return
	}

	sc.peopleInSpace.Collect(ch)
}

type astroResponse struct {
	People []struct {
		Craft string
	}
}

func (sc *spaceCollector) observe() error {
	// ...
	resp, err := http.Get("http://api.open-notify.org/astros.json") // OMIT
	if err != nil {                                                 // OMIT
		return fmt.Errorf("fetching data: %w", err) // OMIT
	} // OMIT
	defer resp.Body.Close() // OMIT
	// OMIT
	astros := astroResponse{}                                          // OMIT
	if err := json.NewDecoder(resp.Body).Decode(&astros); err != nil { // OMIT
		return fmt.Errorf("decoding data: %w", err) // OMIT
	} // OMIT
	// OMIT
	// map of craft to count of astronauts
	m := map[string]int{}
	for _, person := range astros.People {
		m[person.Craft]++
	}

	for craft, value := range m {
		labels := prometheus.Labels{"craft": craft} // HL

		sc.peopleInSpace.With(labels).Set(float64(value)) // HL
	}

	return nil
}

// END OMIT
