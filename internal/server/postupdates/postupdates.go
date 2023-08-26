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
		body := make([]byte, 10000)
		var err error

		var JSONMetrics []storage.Metrics
		var LastMess storage.Metrics

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
			fmt.Printf("Header gzip \n")
			bodyS, err = compress.Decompress(body, req.Header.Get("Content-Encoding"))
			if err != nil {
				fmt.Printf("Error ungzip %v\n", err)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		fmt.Printf("PostUpdates: n =%v, Body: %v \n", n, string(bodyS))

		err = json.Unmarshal(bodyS, &JSONMetrics)
		if err != nil {
			fmt.Printf("Error decode %v \n", err)
			logmy.OutLog(err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("PostUpdates JSONMetrics:  %v \n", JSONMetrics)

		for _, messJSON := range JSONMetrics {

			fmt.Printf(" %v \n", messJSON)
			mess, ok := DecodeMess(mgauge, mmetric, maxlen, messJSON)
			if ok != http.StatusOK {
				res.WriteHeader(ok)
				return
			}
			LastMess = *mess
		}

		// Ответ в JSON
		buf, err := json.Marshal(LastMess)
		if err != nil {
			fmt.Printf("PostUpdates: Error code response: %v \n", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("response: %v \n", string(buf))

		res.Header().Set("Content-Type", "application/json")
		// Сжимю, если нужно
		/*		if req.Header.Get("Content-Encoding") == "gzip" || req.Header.Get("Content-Encoding") == "flate" {
					fmt.Printf("Compress response: \n")
					buf, err = compress.Compress(buf)
					if err != nil {
						fmt.Printf("PostUpdates: Error Compress response: %v \n", err)
						res.WriteHeader(http.StatusBadRequest)
						return

					}
					//		res.Header().Set("Accept-Encoding", "gzip")
					res.Header().Set("Content-Encoding", "gzip")
				}
		*/ //		res.Header().Set("Accept-Encoding", "gzip")

		res.WriteHeader(http.StatusOK)
		res.Write(buf)
	}
}

// распаковываю JSON и записываю с память
func DecodeMess(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter, maxlen int, messJSON storage.Metrics) (*storage.Metrics, int) {
	var err error

	//		fmt.Printf("vMetrics: %v \n", vMetrics)
	if messJSON.ID == "" || (messJSON.MType != "gauge" && messJSON.MType != "counter") || len(messJSON.ID) > 100 {
		logmy.OutLog(err)
		return nil, http.StatusBadRequest
	}

	if messJSON.MType == "gauge" {
		// контроль длинны карты
		if _, ok := mgauge.Get(messJSON.ID); ok != nil && mgauge.Len() > maxlen {
			logmy.OutLog(err)
			return nil, http.StatusBadRequest
		}

		if messJSON.Value == nil {
			messJSON.Value = new(float64)
		}

		// Добавляю новую метрику
		mgauge.Set(messJSON.ID, *messJSON.Value)
		// Выбираю новое значение метрики

		if *messJSON.Value, err = mgauge.Get(messJSON.ID); err != nil {
			logmy.OutLog(err)
			return nil, http.StatusBadRequest
		}
	}

	if messJSON.MType == "counter" {
		// контроль длинны карты
		if _, ok := mmetric.Get(messJSON.ID); ok != nil && mmetric.Len() > maxlen {
			logmy.OutLog(err)
			return nil, http.StatusBadRequest
		}
		if messJSON.Delta == nil {
			messJSON.Delta = new(int64)
		}
		// Добавляю новую метрику
		mmetric.Set(messJSON.ID, *messJSON.Delta)

		// Выбираю новое значение метрики
		if *messJSON.Delta, err = mmetric.Get(messJSON.ID); err != nil {
			logmy.OutLog(err)
			return nil, http.StatusBadRequest
		}
	}

	return &messJSON, http.StatusOK
}
