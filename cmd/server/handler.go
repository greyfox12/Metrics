package main

import (
	"net/http"
	"strconv"
	"strings"
)

func GaugePage(res http.ResponseWriter, req *http.Request) {
	if MemMetric.gauge == nil {
		MemMetric.gauge = make(map[string]float64, LenArr)
	}

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

		if aSt[1] != "update" || aSt[2] != "gauge" {
			res.WriteHeader(http.StatusBadRequest)
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

func CounterPage(res http.ResponseWriter, req *http.Request) {
	if MemMetric.counter == nil {
		MemMetric.counter = make(map[string]int64, LenArr)
	}
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

		if aSt[1] != "update" || aSt[2] != "counter" {
			res.WriteHeader(http.StatusBadRequest)
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

func ErrorPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusMethodNotAllowed)
}
