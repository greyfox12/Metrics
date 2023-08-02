package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	defStoreInterval = 10
	defStorePath     = "/tmp/metrics-db.json"
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
			fmt.Printf("%v\n", err)
		}
	}

	// запускаю сохранение данных в файл
	if vServerParam.StoreInterval > 0 {
		go func(*storage.GaugeCounter, *storage.MetricCounter, getparam.ServerParam) {
			ticker := time.NewTicker(time.Second * time.Duration(vServerParam.StoreInterval))
			defer ticker.Stop()
			for {
				<-ticker.C
				filesave.SaveMetric(gauge, metric, vServerParam.FileStorePath)
			}
		}(gauge, metric, vServerParam)
	}

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(handler.SavePage(gauge, metric, LenArr, vServerParam)) // автосохранение данных

	// определяем хендлер
	r.Route("/", func(r chi.Router) {
		r.Get("/", logmy.RequestLogger(handler.ListMetricPage(gauge, metric)))
		r.Get("/value/gauge/{metricName}", logmy.RequestLogger(handler.OneMetricPage(gauge, metric)))
		r.Get("/value/counter/{metricName}", logmy.RequestLogger(handler.OneMetricPage(gauge, metric)))
		r.Get("/*", logmy.RequestLogger(handler.ErrorPage))

		r.Post("/update", logmy.RequestLogger(handler.PostPage(gauge, metric, LenArr, vServerParam)))
		r.Post("/value", logmy.RequestLogger(handler.OnePostMetricPage(gauge, metric)))
		r.Post("/update/gauge/{metricName}/{metricVal}", logmy.RequestLogger(handler.GaugePage(gauge, metric, LenArr, vServerParam)))
		r.Post("/update/counter/{metricName}/{metricVal}", logmy.RequestLogger(handler.CounterPage(gauge, metric, LenArr, vServerParam)))

		r.Post("/*", logmy.RequestLogger(handler.ErrorPage))

	})

	fmt.Printf("Start Server %v\n", vServerParam.IPAddress)
	log.Fatal(http.ListenAndServe(vServerParam.IPAddress, compress.GzipHandle(r)))
}
