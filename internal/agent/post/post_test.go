package post

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/greyfox12/Metrics/internal/agent/getparam"
)

func TestPostCounter(t *testing.T) {
	var (
		ListCounter map[int]CounterMetric
		ListGauge   map[int]GaugeMetric
	)
	ListGauge = make(map[int]GaugeMetric)
	ListCounter = make(map[int]CounterMetric)

	ListGauge[1] = GaugeMetric{"Alloc", Gauge(5.5)}
	ListGauge[2] = GaugeMetric{"BuckHashSys", Gauge(6)}

	ListCounter[1] = CounterMetric{"PollCount", Counter(100)}
	updateTyp := "update"

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		fmt.Fprintf("NewServer w=%v expected=%v\n", w, expected)
	}))

	defer svr.Close()
	fmt.Printf("svr.URL=%s\n", svr.URL)
	c := NewClient(getparam.TConfig{Address: svr.URL})
	err := c.PostCounter(ListGauge, ListCounter, updateTyp)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}
}
