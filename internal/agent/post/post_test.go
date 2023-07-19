package post

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestPostCounter(t *testing.T) {
	var ListCounter map[int]CounterMetric
	var ListGauge map[int]GaugeMetric
	ListGauge = make(map[int]GaugeMetric)
	ListCounter = make(map[int]CounterMetric)

	ListGauge[1] = GaugeMetric{"Alloc", Gauge(5.5)}
	ListGauge[2] = GaugeMetric{"BuckHashSys", Gauge(6)}

	ListCounter[1] = CounterMetric{"PollCount", Counter(100)}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		fmt.Fprintf("NewServer w=%v expected=%v\n", w, expected)
	}))

	defer svr.Close()
	fmt.Printf("svr.URL=%s\n", svr.URL)
	c := NewClient(svr.URL)
	err := c.PostCounter(ListGauge, ListCounter)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}
}
