package postupdates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/greyfox12/Metrics/internal/server/compress"
	"github.com/greyfox12/Metrics/internal/server/getparam"
	"github.com/greyfox12/Metrics/internal/server/logmy"
	"github.com/greyfox12/Metrics/internal/server/storage"
)

func PostUpdates(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter, maxlen int, cfg getparam.ServerParam) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		fmt.Printf("PostUpdates \n")
		body := make([]byte, 1000)
		var err error
		var resp []byte // Ответ клиенту

		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		aSt := strings.Split(req.URL.Path, "/")
		if len(aSt) < 2 || aSt[1] != "updates" {
			fmt.Printf("Error tags %v \n", len(aSt))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		n, err := req.Body.Read(body)
		if err != nil && n <= 0 {
			fmt.Printf("Error req.Body.Read(body):%v: \n", err)
			fmt.Printf("n =%v, Body: %v \n", n, body)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		bodyS := body[0:n]

		if req.Header.Get("Content-Encoding") == "gzip" || req.Header.Get("Content-Encoding") == "flate" {
			//			fmt.Printf("Header gzip \n")
			bodyS, err = compress.Decompress(body, req.Header.Get("Content-Encoding"))
			if err != nil {
				fmt.Printf("Error ungzip %v\n", err)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		fmt.Printf("PostUpdates: n =%v, Body: %v \n", n, bodyS)

		for _, messJSON := range strings.Split(string(bodyS), "}") {

			if messJSON == "" {
				continue
			}
			messJSON = messJSON + "}"
			fmt.Printf(" %v \n", messJSON)
			mess, ok := decodeMess(mgauge, mmetric, maxlen, messJSON)
			if ok != http.StatusOK {
				res.WriteHeader(ok)
				return
			}
			resp = append(resp, mess...)
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(resp))
	}
}

func decodeMess(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter, maxlen int, messJSON string) ([]byte, int) {
	var err error
	var vMetrics storage.Metrics
	var jsonData []byte

	err = json.Unmarshal([]byte(messJSON), &vMetrics)
	if err != nil {
		fmt.Printf("Error decode %v \n", err)
		logmy.OutLog(err)
		return nil, http.StatusBadRequest
	}
	//		fmt.Printf("vMetrics: %v \n", vMetrics)

	if vMetrics.ID == "" || (vMetrics.MType != "gauge" && vMetrics.MType != "counter") || len(vMetrics.ID) > 100 {
		logmy.OutLog(err)
		return nil, http.StatusBadRequest
	}

	if vMetrics.MType == "gauge" {
		// контроль длинны карты
		if _, ok := mgauge.Get(vMetrics.ID); ok != nil && mgauge.Len() > maxlen {
			logmy.OutLog(err)
			return nil, http.StatusBadRequest
		}

		if vMetrics.Value == nil {
			vMetrics.Value = new(float64)
		}

		// Добавляю новую метрику
		mgauge.Set(vMetrics.ID, *vMetrics.Value)
		// Выбираю новое значение метрики

		if *vMetrics.Value, err = mgauge.Get(vMetrics.ID); err != nil {
			logmy.OutLog(err)
			return nil, http.StatusBadRequest
		}
	}

	if vMetrics.MType == "counter" {
		// контроль длинны карты
		if _, ok := mmetric.Get(vMetrics.ID); ok != nil && mmetric.Len() > maxlen {
			logmy.OutLog(err)
			return nil, http.StatusBadRequest
		}
		if vMetrics.Delta == nil {
			vMetrics.Delta = new(int64)
		}
		// Добавляю новую метрику
		mmetric.Set(vMetrics.ID, *vMetrics.Delta)

		// Выбираю новое значение метрики
		if *vMetrics.Delta, err = mmetric.Get(vMetrics.ID); err != nil {
			logmy.OutLog(err)
			return nil, http.StatusBadRequest
		}
	}

	if jsonData, err = json.Marshal(vMetrics); err != nil {
		logmy.OutLog(err)
		return nil, http.StatusBadRequest
	}

	return jsonData, http.StatusOK
}
