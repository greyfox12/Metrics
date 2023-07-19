package post

import (
	"fmt"
	"io"
	"net/http"
	"time"
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

func (c Client) PostCounter(ga map[int]GaugeMetric, co map[int]CounterMetric) error {
	fmt.Printf("Time: %v\n", time.Now().Unix())
	//	fmt.Printf("URL: %v\n", c.url)

	for _, val := range ga {
		s := fmt.Sprintf("%s/update/gauge/%s/%v", c.url, val.Name, val.Val)

		resp, err := http.Post(s, "Content-Type: text/plain", nil)
		if err != nil {
			return error(err)
		}

		defer resp.Body.Close()
		_, _ = io.ReadAll(resp.Body)
	}

	for _, val := range co {
		s := fmt.Sprintf("%s/update/counter/%s/%v", c.url, val.Name, val.Val)
		resp, err := http.Post(s, "Content-Type: text/plain", nil)

		if err != nil {
			return error(err)
		}

		defer resp.Body.Close()
		if _, err := io.ReadAll(resp.Body); err != nil {
			return error(err)
		}

	}
	return nil
}
