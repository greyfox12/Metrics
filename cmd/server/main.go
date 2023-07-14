package main

import (
	"net/http"
)

const LenArr = 10000

type tMetric struct {
	gauge   map[string]float64
	counter map[string]int64
}

var MemMetric tMetric

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/gauge/`, GaugePage)
	mux.HandleFunc(`/update/counter/`, CounterPage)
	mux.HandleFunc(`/`, ErrorPage)
	MemMetric.gauge = make(map[string]float64, LenArr)
	MemMetric.counter = make(map[string]int64, LenArr)
	//gaugeMetric.val = m

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
