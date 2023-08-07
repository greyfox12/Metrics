package post

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/greyfox12/Metrics/internal/agent/compress"
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
		jsonZip, err := compress.Compress(jsonData)
		fmt.Printf("jsonZip: %+v\n", jsonZip)
		if err != nil {
			return error(err)
		}
		//		fmt.Printf("jsonZip=%v\n", jsonZip)
		//	http.Header.Set("Content-Encoding", "gzip")

		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, err := http.NewRequest("POST", adrstr, bytes.NewBuffer(jsonZip))
		if err != nil {
			return error(err)
		}
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Add("Content-Type", "application/json")
		response, err := client.Do(req)
		if err != nil {
			return error(err)
		}
		body, err := io.ReadAll(response.Body)
		defer response.Body.Close()

		if err != nil {
			return error(err)
		}

		fmt.Println("response Body:", body)
	}

	for _, val := range co {

		st := Metrics{ID: val.Name, MType: "counter", Delta: (*int64)(&val.Val)}

		jsonData, err := json.Marshal(st)
		if err != nil {
			return error(err)
		}

		jsonZip, err := compress.Compress(jsonData)
		fmt.Printf("jsonZip: %+v\n", jsonZip)
		if err != nil {
			return error(err)
		}
		//		fmt.Printf("jsonZip=%v\n", jsonZip)

		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, err := http.NewRequest("POST", adrstr, bytes.NewBuffer(jsonZip))
		if err != nil {
			return error(err)
		}
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Add("Content-Type", "application/json")
		response, err := client.Do(req)
		if err != nil {
			return error(err)
		}

		body, err := io.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			return error(err)
		}

		fmt.Println("response Body:", body)

	}

	return nil
}
