package post

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Counter int64
type Gauge float64

type CounterMetric struct {
	Name string
	Val  Counter
}

type GaugeMetric struct {
	Name string
	Val  Gauge
}

type Client struct {
	url string
}

func NewClient(url string) Client {
	return Client{url}
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (c Client) PostCounter(ga map[int]GaugeMetric, co map[int]CounterMetric) error {
	//	fmt.Printf("Time: %v\n", time.Now().Unix())

	adrstr := fmt.Sprintf("%s/update", c.url)
	for _, val := range ga {
		st := Metrics{ID: val.Name, MType: "gauge", Value: (*float64)(&val.Val)}

		jsonData, err := json.Marshal(st)
		if err != nil {
			return error(err)
		}

		//		fmt.Printf("Ответ сервера: %v\n", string(jsonData))

		resp, err := http.Post(adrstr, "Content-Type: application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return error(err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return error(err)
		}
		fmt.Println("response Body:", string(body))
	}

	for _, val := range co {

		st := Metrics{ID: val.Name, MType: "counter", Delta: (*int64)(&val.Val)}

		jsonData, err := json.Marshal(st)
		if err != nil {
			return error(err)
		}

		//		fmt.Println(string(jsonData))

		resp, err := http.Post(adrstr, "Content-Type: application/json", bytes.NewBuffer(jsonData))

		if err != nil {
			return error(err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return error(err)
		}
		fmt.Println("response Body:", string(body))

	}
	return nil
}
