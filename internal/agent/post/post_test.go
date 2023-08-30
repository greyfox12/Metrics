package post

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/greyfox12/Metrics/internal/agent/getparam"
)

func TestPostCounter(t *testing.T) {

	List := make(map[int]CollectMetr)

	List[1] = CollectMetr{ID: "Alloc", MType: "gauge", Value: 5.5}
	List[2] = CollectMetr{ID: "BuckHashSys", MType: "gauge", Value: 6}

	List[3] = CollectMetr{ID: "PollCount", MType: "counter", Value: 100}
	updateTyp := "update"

	jobs := make(chan map[int]CollectMetr, 1)
	jobs <- List
	results := make(chan error, 1)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		fmt.Fprintf("NewServer w=%v expected=%v\n", w, expected)
	}))

	defer svr.Close()
	fmt.Printf("svr.URL=%s\n", svr.URL)
	c := NewClient(getparam.TConfig{Address: svr.URL})
	go c.PostCounter(jobs, results, updateTyp)
	err := <-results
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}
}
