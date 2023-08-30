package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	//	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/greyfox12/Metrics/internal/server/compress"
	"github.com/greyfox12/Metrics/internal/server/dbstore"
	"github.com/greyfox12/Metrics/internal/server/filesave"
	"github.com/greyfox12/Metrics/internal/server/getparam"
	"github.com/greyfox12/Metrics/internal/server/getping"
	"github.com/greyfox12/Metrics/internal/server/handler"
	"github.com/greyfox12/Metrics/internal/server/hash"
	"github.com/greyfox12/Metrics/internal/server/logmy"
	"github.com/greyfox12/Metrics/internal/server/postupdates"
	"github.com/greyfox12/Metrics/internal/server/storage"
)

const (
	LenArr           = 10000
	defServerAdr     = "localhost:8080"
	defStoreInterval = 100
	defStorePath     = "/tmp/metrics-db.json"
	defRestore       = true
	defDSN           = "host=localhost user=videos password=videos dbname=postgres sslmode=disable"
)

func main() {

	serverStart()
}

// Запускаю сервер
func serverStart() {
	var db *sql.DB

	vServerParam := getparam.ServerParam{IPAddress: defServerAdr,
		StoreInterval: defStoreInterval,
		FileStorePath: defStorePath,
		Restore:       defRestore,
		DSN:           defDSN,
	}
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
	if vServerParam.Restore && vServerParam.OnFile {
		if err := filesave.LoadMetric(gauge, metric, vServerParam.FileStorePath); err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	// Подключение к БД
	if vServerParam.OnDSN {
		var err error
		fmt.Printf("DSN: %v\n", vServerParam.DSN)
		db, err = sql.Open("pgx", vServerParam.DSN)
		if err != nil {
			logmy.OutLog(err)
			fmt.Printf("Error connect DB: %v\n", err)
		}
		defer db.Close()

		if err = dbstore.CreateDB(db); err != nil {
			vServerParam.OnDSN = false
		}

		if vServerParam.Restore {
			if err := dbstore.LoadMetric(gauge, metric, db); err != nil {
				logmy.OutLog(err)
				//				fmt.Printf("%v\n", err)
			}
		}
	}

	// запускаю сохранение данных в файл
	if (vServerParam.OnFile || vServerParam.OnDSN) && vServerParam.StoreInterval > 0 {
		go func(*storage.GaugeCounter, *storage.MetricCounter, *sql.DB, getparam.ServerParam) {
			ticker := time.NewTicker(time.Second * time.Duration(vServerParam.StoreInterval))
			defer ticker.Stop()
			for {
				<-ticker.C
				if vServerParam.OnFile {
					filesave.SaveMetric(gauge, metric, vServerParam.FileStorePath)
				}
				if vServerParam.OnDSN {
					dbstore.SaveMetric(gauge, metric, db)
				}
			}
		}(gauge, metric, db, vServerParam)

	}

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	if vServerParam.OnFile || vServerParam.OnDSN { // автосохранение данных
		r.Use(handler.SavePage(gauge, metric, db, vServerParam))
	}

	// определяем хендлер
	r.Route("/", func(r chi.Router) {
		r.Get("/", logmy.RequestLogger(handler.ListMetricPage(gauge, metric, vServerParam)))
		r.Get("/value/gauge/{metricName}", logmy.RequestLogger(handler.OneMetricPage(gauge, metric, vServerParam)))
		r.Get("/value/counter/{metricName}", logmy.RequestLogger(handler.OneMetricPage(gauge, metric, vServerParam)))
		r.Get("/ping", logmy.RequestLogger(getping.GetPing(db)))
		r.Get("/*", logmy.RequestLogger(handler.ErrorPage))

		r.Post("/updates", logmy.RequestLogger(postupdates.PostUpdates(gauge, metric, LenArr, vServerParam)))
		r.Post("/update", logmy.RequestLogger(handler.PostPage(gauge, metric, LenArr, vServerParam)))
		r.Post("/value", logmy.RequestLogger(handler.OnePostMetricPage(gauge, metric, vServerParam)))
		r.Post("/update/gauge/{metricName}/{metricVal}", logmy.RequestLogger(handler.GaugePage(gauge, metric, LenArr, vServerParam)))
		r.Post("/update/counter/{metricName}/{metricVal}", logmy.RequestLogger(handler.CounterPage(gauge, metric, LenArr, vServerParam)))

		r.Post("/*", logmy.RequestLogger(handler.ErrorPage))

	})

	fmt.Printf("Start Server %v\n", vServerParam.IPAddress)

	hd := compress.GzipHandle(compress.GzipRead(hash.HashHandle(hash.HashWriteHandle(r, vServerParam), vServerParam)))
	log.Fatal(http.ListenAndServe(vServerParam.IPAddress, hd))
}
