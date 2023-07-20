package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/greyfox12/Metrics/internal/server/getparam"
	"github.com/greyfox12/Metrics/internal/server/handler"
	"github.com/greyfox12/Metrics/internal/server/storage"
)

const (
	LenArr       = 10000
	defServerAdr = "localhost:8080"
)

func main() {

	// запрашиваю параметры ключей-переменных окружения
	IPAddress := getparam.Param(defServerAdr)

	gauge := new(storage.GaugeCounter)
	gauge.Init(LenArr)
	metric := new(storage.MetricCounter)
	metric.Init(LenArr)

	r := chi.NewRouter()

	// определяем хендлер, который выводит определённую машину
	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.ListMetricPage(gauge, metric))
		r.Get("/value/gauge/{metricName}", handler.OneMetricPage(gauge, metric))
		r.Get("/value/counter/{metricName}", handler.OneMetricPage(gauge, metric))
		r.Route("/update", func(r chi.Router) {
			r.Post("/gauge/{metricName}/{metricVal}", handler.GaugePage(gauge, LenArr))
			r.Post("/counter/{metricName}/{metricVal}", handler.CounterPage(metric, LenArr))
			r.Post("/*", handler.ErrorPage)
		})
	})

	log.Fatal(http.ListenAndServe(IPAddress, r))
}
