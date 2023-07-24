package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/greyfox12/Metrics/internal/server/logmy"
	"github.com/greyfox12/Metrics/internal/server/storage"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func GaugePage(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter, maxlen int) http.HandlerFunc {
	return logmy.RequestLogger(func(res http.ResponseWriter, req *http.Request) {

		//		fmt.Printf("GaugePage \n")
		var vMetrics Metrics
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		aSt := strings.Split(req.URL.Path, "/")

		if len(aSt) != 2 || aSt[1] != "update" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(req.Body)

		if err := decoder.Decode(&vMetrics); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		//		fmt.Printf("vMetrics: %v \n", vMetrics)

		if vMetrics.ID == "" || (vMetrics.MType != "gauge" && vMetrics.MType != "counter") || len(vMetrics.ID) > 100 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if vMetrics.MType == "gauge" {
			// контроль длинны карты
			if _, ok := mgauge.Get(vMetrics.ID); ok != nil && mgauge.Len() > maxlen {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			if vMetrics.Value == nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			// Добавляю новую метрику
			mgauge.Set(vMetrics.ID, *vMetrics.Value)
			// Выбираю новое значение метрики
			var ok error
			if *vMetrics.Value, ok = mgauge.Get(vMetrics.ID); ok != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if vMetrics.MType == "counter" {
			// контроль длинны карты
			if _, ok := mmetric.Get(vMetrics.ID); ok != nil && mmetric.Len() > maxlen {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			if vMetrics.Delta == nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			// Добавляю новую метрику
			mmetric.Set(vMetrics.ID, *vMetrics.Delta)
			// Выбираю новое значение метрики
			var ok error
			if *vMetrics.Delta, ok = mmetric.Get(vMetrics.ID); ok != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		jsonData, err := json.Marshal(vMetrics)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return

		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(jsonData))
	})
}

func ErrorPage() http.HandlerFunc {
	return logmy.RequestLogger(func(res http.ResponseWriter, req *http.Request) {
		//		fmt.Printf("Error page %v\n", http.MethodPost)
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
	})
}

func ListMetricPage(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter) http.HandlerFunc {
	return logmy.RequestLogger(func(res http.ResponseWriter, req *http.Request) {
		//		fmt.Printf("ListMetric page \n")

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
		res.WriteHeader(http.StatusOK)

		io.WriteString(res, strings.Join(body, "\n"))
	})
}

func OneMetricPage(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter) http.HandlerFunc {
	return logmy.RequestLogger(func(res http.ResponseWriter, req *http.Request) {

		//		fmt.Printf("OneMetricPage \n")
		var Val string
		var retInt int64
		var retFloat float64
		var ok error
		aSt := strings.Split(req.URL.Path, "/")

		metricName := aSt[3]
		if aSt[2] == "gauge" {
			if retFloat, ok = mgauge.Get(metricName); ok != nil {
				res.WriteHeader(http.StatusNotFound)
				return
			} else {
				Val = fmt.Sprintf("%v", retFloat)
			}
		} else {
			if retInt, ok = mmetric.Get(metricName); ok != nil {
				res.WriteHeader(http.StatusNotFound)
				return
			}
			Val = fmt.Sprintf("%v", retInt)
		}

		res.WriteHeader(http.StatusOK)
		io.WriteString(res, fmt.Sprintf("%v", Val))

	})
}

func OnePostMetricPage(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter) http.HandlerFunc {
	return logmy.RequestLogger(func(res http.ResponseWriter, req *http.Request) {

		//		fmt.Printf("OneMetricPage \n")
		var vMetrics Metrics
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		decoder := json.NewDecoder(req.Body)

		if err := decoder.Decode(&vMetrics); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		//		fmt.Printf("vMetrics: %v \n", vMetrics)

		if vMetrics.ID == "" || (vMetrics.MType != "gauge" && vMetrics.MType != "counter") || len(vMetrics.ID) > 100 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if vMetrics.MType == "gauge" {
			var ok error
			if *vMetrics.Value, ok = mgauge.Get(vMetrics.ID); ok != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if vMetrics.MType == "counter" {
			// контроль длинны карты
			var ok error
			if *vMetrics.Delta, ok = mmetric.Get(vMetrics.ID); ok != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		jsonData, err := json.Marshal(vMetrics)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return

		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(jsonData))
	})
}
