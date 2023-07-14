package main

import (
	"fmt"
	"net/http"
)

type Client struct {
	url string
}

func NewClient(url string) Client {
	return Client{url}
}

func (c Client) PostCounter(ga map[int]GaugeMetric, co map[int]CounterMetric) int {
	//	fmt.Printf("Time: %v\n", time.Now().Unix())
	fmt.Printf("URL: %v\n", c.url)

	for _, val := range ga {
		s := fmt.Sprintf("%s/update/gauge/%s/%v", c.url, val.name, val.Val)
		resp, err := http.Post(s, "Content-Type: text/plain", nil)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return -1
			//			panic(err)
		}
		defer resp.Body.Close()
	}

	for _, val := range co {
		s := fmt.Sprintf("%s/update/counter/%s/%v", c.url, val.name, val.Val)
		resp, err := http.Post(s, "Content-Type: text/plain", nil)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return -1
			//			panic(err)
		}
		defer resp.Body.Close()

	}
	return 0
}
