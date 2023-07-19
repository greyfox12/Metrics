package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/greyfox12/Metrics/internal/server/storage"
)

func GaugePage(mgauge storage.GaugeCounter, maxlen int) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		aSt := strings.Split(req.URL.Path, "/")

		if len(aSt) != 5 || aSt[1] != "update" || aSt[2] != "gauge" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		metricName := aSt[3]
		metricVal := aSt[4]
		if metricName == "" || metricVal == "" || len(metricName) > 100 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		metricCn, err := strconv.ParseFloat(metricVal, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// контроль длинны карты
		if _, ok := mgauge.Get(metricName); ok != nil && mgauge.Len() > maxlen {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// Добавляю новую метрику
		mgauge.Set(metricName, metricCn)

		return
	}
}

func CounterPage(mmetric storage.MetricCounter, maxlen int) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		aSt := strings.Split(req.URL.Path, "/")
		if len(aSt) != 5 || aSt[1] != "update" || aSt[2] != "counter" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		metricName := aSt[3]
		metricVal := aSt[4]
		if metricName == "" || metricVal == "" || len(metricName) > 100 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// Проверка корректности
		metricCn, err := strconv.ParseInt(metricVal, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// контроль длинны карты
		if _, ok := mmetric.Get(metricName); ok != nil && mmetric.Len() > maxlen {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// Добавляю новую метрику
		mmetric.Set(metricName, metricCn)

		return
	}
}

func ErrorPage(res http.ResponseWriter, req *http.Request) {
	//	fmt.Printf("req.Method3=%v\n", req.Method)
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	aSt := strings.Split(req.URL.Path, "/")
	if aSt[1] == "update" && aSt[2] != "counter" && aSt[2] != "gauge" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusNotFound)
	return
}

func ListMetricPage(mgauge storage.GaugeCounter, mmetric storage.MetricCounter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var body []string

		for _, val := range mgauge.Keylist() {
			if v, ok := mgauge.Get(val); ok == nil {
				body = append(body, fmt.Sprintf("%s = %v", val, v))
			}
		}

		for _, val := range mmetric.Keylist() {
			if v, ok := mmetric.Get(val); ok == nil {
				body = append(body, fmt.Sprintf("%s = %v", val, v))
			}
		}

		io.WriteString(res, strings.Join(body, "\n"))
	}
}

func OneMetricPage(mgauge storage.GaugeCounter, mmetric storage.MetricCounter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var Val string
		var ret_int int64
		var ret_float float64
		var ok error
		aSt := strings.Split(req.URL.Path, "/")

		metricName := aSt[3]
		if aSt[2] == "gauge" {
			if ret_float, ok = mgauge.Get(metricName); ok != nil {
				res.WriteHeader(http.StatusNotFound)
				return
			} else {
				Val = fmt.Sprintf("%v", ret_float)
			}
		} else {
			if ret_int, ok = mmetric.Get(metricName); ok != nil {
				res.WriteHeader(http.StatusNotFound)
				return
			}
			Val = fmt.Sprintf("%v", ret_int)
		}

		io.WriteString(res, fmt.Sprintf("%v", Val))
	}
}
