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

func (c Client) PostCounter(ga map[int]GaugeMetric, co map[int]CounterMetric, updateTyp string) error {
	//	fmt.Printf("Time: %v\n", time.Now().Unix())
	var stArr []Metrics

	adrstr := fmt.Sprintf("%s/%s", c.url, updateTyp)

	for _, val := range ga {
		st := Metrics{ID: val.Name, MType: "gauge", Value: (*float64)(&val.Val)}
		stArr = append(stArr, st)

		if updateTyp == "update" {
			if ok := postMess(st, adrstr); ok != nil {
				fmt.Printf("Error post: %v, %v\n", st, ok)
			}
		}
	}

	for _, val := range co {

		st := Metrics{ID: val.Name, MType: "counter", Delta: (*int64)(&val.Val)}
		stArr = append(stArr, st)

		if updateTyp == "update" {
			if ok := postMess(st, adrstr); ok != nil {
				fmt.Printf("Error post: %v, %v\n", st, ok)
			}
		}
	}

	if updateTyp == "update" {
		return nil
	}

	//	fmt.Printf("stArr: %v\n", stArr)
	if ok := postUpdates(stArr, adrstr); ok != nil {
		fmt.Printf("Error posts:  %v\n", ok)
	}

	return nil
}

// Вывод по одной записи
func postMess(st Metrics, adrstr string) error {
	var err error
	jsonData, err := json.Marshal(st)
	if err != nil {
		return error(err)
	}

	err = ActPost(jsonData, adrstr)
	if err != nil {
		return error(err)
	}
	//	fmt.Println("response Body:", body)
	return nil
}

// Вывод слайса
func postUpdates(stArr []Metrics, adrstr string) error {
	var err error
	var buf []byte

	for _, st := range stArr {
		jsonData, err := json.Marshal(st)
		if err != nil {
			return error(err)
		}
		buf = append(buf, jsonData...)
		//		buf = append(buf, []byte{'\n'}...)
	}

	err = ActPost(buf, adrstr)
	if err != nil {
		return error(err)
	}
	//	fmt.Println("response Body:", body)
	return nil
}

func ActPost(buf []byte, adrstr string) error {

	fmt.Printf("JSON: %v\n", string(buf))

	jsonZip, err := compress.Compress(buf)
	//		fmt.Printf("jsonZip: %+v\n", jsonZip)
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

	_, err = io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return error(err)
	}
	//	fmt.Println("response Body:", body)
	return nil
}
