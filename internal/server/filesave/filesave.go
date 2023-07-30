package filesave

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/greyfox12/Metrics/internal/server/handler"
	"github.com/greyfox12/Metrics/internal/server/storage"
)

func SaveMetric(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter, fileName string) error {
	//	var vMetrics handler.Metrics

	fmt.Printf("Save file\n")
	var buf []byte
	for _, val := range mgauge.Keylist() {
		v, err := mgauge.Get(val)
		if err != nil {
			continue
		}
		st := handler.Metrics{ID: val, MType: "gauge", Value: (*float64)(&v)}

		jsonData, err := json.Marshal(st)
		if err != nil {
			return error(err)
		}
		buf = append(buf, jsonData...)
		buf = append(buf, '\n')
		//		os.WriteFile(fileName, jsonData, 0666)
	}

	for _, val := range mmetric.Keylist() {
		v, err := mmetric.Get(val)
		if err != nil {
			continue
		}
		st := handler.Metrics{ID: val, MType: "metrica", Delta: (*int64)(&v)}

		jsonData, err := json.Marshal(st)
		if err != nil {
			return error(err)
		}
		buf = append(buf, jsonData...)
		buf = append(buf, '\n')
	}
	os.WriteFile(fileName, buf, 0666)
	return nil
}

// // Загрузка данных из Файла
func LoadMetric(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter, fileName string) error {
	var scanner *bufio.Scanner
	fmt.Printf("Load data from file %s\n", fileName)

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Printf("LoadMetric - Error open file =%v\n", fileName)
		return err
	}
	scanner = bufio.NewScanner(file)

	for scanner.Scan() {

		// читаем данные из scanner
		data := scanner.Bytes()

		metric := handler.Metrics{}
		err := json.Unmarshal(data, &metric)
		//		fmt.Printf("metric=%v\n", metric)
		if err != nil {
			fmt.Printf("Error Unmarshal string =%v\n", err)
			return err
		}
		if metric.ID == "" || (metric.Value == nil && metric.Delta == nil) {
			return errors.New("LoadMetric - Error field")
		}

		switch metric.MType {
		case "gauge":
			mgauge.Set(metric.ID, float64(*metric.Value))
		case "metrica":
			mmetric.Set(metric.ID, int64(*metric.Delta))
		default:
			fmt.Printf("Unknow Mtype in file =%v\n", metric.MType)
		}
	}

	return nil
}
