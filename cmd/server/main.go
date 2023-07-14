package main

import (
	"net/http"
	"strconv"
	"strings"
)

const LenArr = 10000

type tMetric struct {
	gauge   map[string]float64
	counter map[string]int64
}

var MemMetric tMetric

func gaugePage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		//		res.WriteHeader(http.StatusOK)
		st := req.URL.Path
		// Проверка корректности
		if len(st) > 100 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		aSt := strings.Split(st, "/")
		if len(aSt) != 5 || len(aSt[3]) == 0 || aSt[2] != "gauge" {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		metricVal, err := strconv.ParseFloat(aSt[4], 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// контроль длинны карты
		if _, ok := MemMetric.gauge[aSt[3]]; !ok && len(MemMetric.gauge) > LenArr {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// Добавляю новую метрику

		MemMetric.gauge[aSt[3]] = metricVal

		return

	}
	res.WriteHeader(http.StatusMethodNotAllowed)
}

func counterPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		//		res.WriteHeader(http.StatusOK)
		st := req.URL.Path
		//		io.WriteString(res, req.URL.Path)

		//		return
		// Проверка корректности
		if len(st) > 100 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		aSt := strings.Split(st, "/")
		if len(aSt) != 5 || len(aSt[3]) == 0 || aSt[2] != "counter" {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		metricVal, err := strconv.ParseInt(aSt[4], 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// контроль длинны карты
		if _, ok := MemMetric.counter[aSt[3]]; !ok && len(MemMetric.counter) > LenArr {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// Добавляю новую метрику

		MemMetric.counter[aSt[3]] += metricVal

		return
	}
	res.WriteHeader(http.StatusMethodNotAllowed)
}

func errorPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/gauge/`, gaugePage)
	mux.HandleFunc(`/update/counter/`, counterPage)
	mux.HandleFunc(`/`, errorPage)
	MemMetric.gauge = make(map[string]float64, LenArr)
	MemMetric.counter = make(map[string]int64, LenArr)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
