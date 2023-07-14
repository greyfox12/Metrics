package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

const LenArr = 10000
const DefServerAdr = "localhost:8080"

type tMetric struct {
	gauge   map[string]float64
	counter map[string]int64
}

var MemMetric tMetric

func main() {
	var cfg string
	var ok bool

	if cfg, ok = os.LookupEnv("ADDRESS"); !ok {
		cfg = DefServerAdr
	}

	fmt.Printf("ServerAdr = %v\n", cfg)

	IPAddress := flag.String("a", cfg, "Endpoint server IP address host:port")
	flag.Parse()

	r := chi.NewRouter()

	// определяем хендлер, который выводит определённую машину
	r.Route("/", func(r chi.Router) {
		r.Get("/", ListMetricPage)
		r.Get("/value/gauge/{metricName}", OneMetricPage)
		r.Get("/value/counter/{metricName}", OneMetricPage)
		r.Route("/update", func(r chi.Router) {
			r.Post("/gauge/{metricName}/{metricVal}", GaugePage)
			r.Post("/counter/{metricName}/{metricVal}", CounterPage)
			r.Post("/*", ErrorPage)
		})
	})

	log.Fatal(http.ListenAndServe(*IPAddress, r))
}
