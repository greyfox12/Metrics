package storage

import (
	"errors"
	"sort"
	"sync"
)

type MetricCounter struct {
	Val   map[string]int64
	Mutex sync.RWMutex
}

func (m MetricCounter) Set(key string, v int64) {
	m.Mutex.Lock()

	defer m.Mutex.Unlock()

	m.Val[key] += v
}

func (m MetricCounter) Get(key string) (int64, error) {
	m.Mutex.RLock()

	defer m.Mutex.RUnlock()

	if v, ok := m.Val[key]; ok {
		return v, nil
	}

	return 0, errors.New("not found")
}

func (m *MetricCounter) Init(defLen int) {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	m.Val = make(map[string]int64, defLen)
}

func (m MetricCounter) Len() int {
	m.Mutex.RLock()

	defer m.Mutex.RUnlock()

	return len(m.Val)
}

func (m MetricCounter) Keylist() []string {
	var ret []string
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()

	for key, _ := range m.Val {
		ret = append(ret, key)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})

	return ret
}

// ///////////////////////////////
type GaugeCounter struct {
	Val   map[string]float64
	Mutex sync.RWMutex
}

func (g GaugeCounter) Set(key string, v float64) {
	g.Mutex.Lock()

	defer g.Mutex.Unlock()

	g.Val[key] = v
}

func (g GaugeCounter) Get(key string) (float64, error) {
	g.Mutex.RLock()

	defer g.Mutex.RUnlock()

	if v, ok := g.Val[key]; ok {
		return v, nil
	}

	return 0, errors.New("not found")
}

func (g GaugeCounter) Len() int {
	g.Mutex.RLock()

	defer g.Mutex.RUnlock()

	return len(g.Val)
}

func (g *GaugeCounter) Init(defLen int) {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()

	g.Val = make(map[string]float64, defLen)
}

func (g GaugeCounter) Keylist() []string {
	var ret []string
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()

	for key, _ := range g.Val {
		ret = append(ret, key)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})

	return ret
}
