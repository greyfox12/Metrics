package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/greyfox12/Metrics/internal/server/compress"
	"github.com/greyfox12/Metrics/internal/server/filesave"
	"github.com/greyfox12/Metrics/internal/server/getparam"
	"github.com/greyfox12/Metrics/internal/server/handler"
	"github.com/greyfox12/Metrics/internal/server/logmy"
	"github.com/greyfox12/Metrics/internal/server/storage"
)

const (
	LenArr           = 10000
	defServerAdr     = "localhost:8080"
	defStoreInterval = 300
	defStorePath     = "metrics-db.json"
	defRestore       = true
)

func main() {

	vServerParam := getparam.ServerParam{IPAddress: defServerAdr,
		StoreInterval: defStoreInterval,
		FileStorePath: defStorePath,
		Restore:       defRestore}
	// запрашиваю параметры ключей-переменных окружения
	vServerParam = getparam.Param(&vServerParam)

	// Инициализирую логирование
	if ok := logmy.Initialize("info"); ok != nil {
		panic(ok)
	}

	gauge := new(storage.GaugeCounter)
	gauge.Init(LenArr)
	metric := new(storage.MetricCounter)
	metric.Init(LenArr)

	// Загрузка данных из файла
	if vServerParam.Restore {
		if err := filesave.LoadMetric(gauge, metric, vServerParam.FileStorePath); err != nil {
			panic(err)
		}
	}

	// запускаю сохранение данных в файл
	if vServerParam.StoreInterval > 0 {
		go func(*storage.GaugeCounter, *storage.MetricCounter, getparam.ServerParam) {
			for {
				time.Sleep(time.Duration(vServerParam.StoreInterval) * time.Second)
				filesave.SaveMetric(gauge, metric, vServerParam.FileStorePath)
			}
		}(gauge, metric, vServerParam)
	}

	r := chi.NewRouter()

	// определяем хендлер, который выводит определённую машину
	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.ListMetricPage(gauge, metric))
		r.Get("/value/gauge/{metricName}", handler.OneMetricPage(gauge, metric))
		r.Get("/value/counter/{metricName}", handler.OneMetricPage(gauge, metric))
		r.Get("/*", handler.ErrorPage())
		//		r.Route("/update", func(r chi.Router) {
		r.Post("/update/", handler.PostPage(gauge, metric, LenArr))
		r.Post("/value/", handler.OnePostMetricPage(gauge, metric))
		r.Post("/update", handler.PostPage(gauge, metric, LenArr))
		r.Post("/value", handler.OnePostMetricPage(gauge, metric))
		r.Post("/update/gauge/{metricName}/{metricVal}", handler.GaugePage(gauge, LenArr))
		r.Post("/update/counter/{metricName}/{metricVal}", handler.CounterPage(metric, LenArr))

		r.Post("/*", handler.ErrorPage())

		//		})
	})

	log.Fatal(http.ListenAndServe(vServerParam.IPAddress, compress.GzipHandle(r)))
}
