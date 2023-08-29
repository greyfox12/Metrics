package collectmetric

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/greyfox12/Metrics/internal/agent/post"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func CollectGauge(PollCount int64, job chan<- map[int]post.CollectMetr) {
	var m runtime.MemStats

	//	var List map[int]post.CollectMetr
	List := make(map[int]post.CollectMetr)

	runtime.ReadMemStats(&m)
	fmt.Printf("Start %v\n", PollCount)
	List[1] = post.CollectMetr{ID: "Alloc", MType: "gauge", Value: post.Gauge(m.Alloc)}
	List[2] = post.CollectMetr{ID: "BuckHashSys", MType: "gauge", Value: post.Gauge(m.BuckHashSys)}
	List[3] = post.CollectMetr{ID: "Frees", MType: "gauge", Value: post.Gauge(m.Frees)}
	List[4] = post.CollectMetr{ID: "GCCPUFraction", MType: "gauge", Value: post.Gauge(m.GCCPUFraction)}
	List[5] = post.CollectMetr{ID: "GCSys", MType: "gauge", Value: post.Gauge(m.GCSys)}
	List[6] = post.CollectMetr{ID: "HeapAlloc", MType: "gauge", Value: post.Gauge(m.HeapAlloc)}
	List[7] = post.CollectMetr{ID: "HeapIdle", MType: "gauge", Value: post.Gauge(m.HeapIdle)}
	List[8] = post.CollectMetr{ID: "HeapObjects", MType: "gauge", Value: post.Gauge(m.HeapObjects)}
	List[9] = post.CollectMetr{ID: "HeapReleased", MType: "gauge", Value: post.Gauge(m.HeapReleased)}
	List[10] = post.CollectMetr{ID: "HeapSys", MType: "gauge", Value: post.Gauge(m.HeapSys)}
	List[11] = post.CollectMetr{ID: "LastGC", MType: "gauge", Value: post.Gauge(m.LastGC)}
	List[12] = post.CollectMetr{ID: "Lookups", MType: "gauge", Value: post.Gauge(m.Lookups)}
	List[13] = post.CollectMetr{ID: "MCacheInuse", MType: "gauge", Value: post.Gauge(m.MCacheInuse)}
	List[14] = post.CollectMetr{ID: "MCacheSys", MType: "gauge", Value: post.Gauge(m.MCacheSys)}
	List[15] = post.CollectMetr{ID: "Mallocs", MType: "gauge", Value: post.Gauge(m.Mallocs)}
	List[16] = post.CollectMetr{ID: "NextGC", MType: "gauge", Value: post.Gauge(m.NextGC)}
	List[17] = post.CollectMetr{ID: "NumForcedGC", MType: "gauge", Value: post.Gauge(m.NumForcedGC)}
	List[18] = post.CollectMetr{ID: "NumGC", MType: "gauge", Value: post.Gauge(m.NumGC)}
	List[19] = post.CollectMetr{ID: "OtherSys", MType: "gauge", Value: post.Gauge(m.OtherSys)}
	List[20] = post.CollectMetr{ID: "PauseTotalNs", MType: "gauge", Value: post.Gauge(m.PauseTotalNs)}
	List[21] = post.CollectMetr{ID: "StackInuse", MType: "gauge", Value: post.Gauge(m.StackInuse)}
	List[22] = post.CollectMetr{ID: "StackSys", MType: "gauge", Value: post.Gauge(m.StackSys)}
	List[23] = post.CollectMetr{ID: "Sys", MType: "gauge", Value: post.Gauge(m.Sys)}
	List[24] = post.CollectMetr{ID: "TotalAlloc", MType: "gauge", Value: post.Gauge(m.TotalAlloc)}
	List[25] = post.CollectMetr{ID: "RandomValue", MType: "gauge", Value: post.Gauge(rand.Float64())}
	List[26] = post.CollectMetr{ID: "HeapInuse", MType: "gauge", Value: post.Gauge(m.HeapInuse)}
	List[27] = post.CollectMetr{ID: "MSpanInuse", MType: "gauge", Value: post.Gauge(m.MSpanInuse)}
	List[28] = post.CollectMetr{ID: "MSpanSys", MType: "gauge", Value: post.Gauge(m.MSpanSys)}

	List[29] = post.CollectMetr{ID: "PollCount", MType: "gauge", Delta: post.Counter(PollCount)}

	job <- List
	//	return List
}

func CollectAdd(PollCount int64, job chan<- map[int]post.CollectMetr) {

	var List map[int]post.CollectMetr
	List = make(map[int]post.CollectMetr)

	fmt.Printf("Start Add:%v\n", PollCount)
	v, _ := mem.VirtualMemory()
	h, _ := load.Avg()
	List[1] = post.CollectMetr{ID: "TotalMemory", MType: "gauge", Value: post.Gauge(v.Total)}
	List[2] = post.CollectMetr{ID: "FreeMemory", MType: "gauge", Value: post.Gauge(v.Free)}
	List[3] = post.CollectMetr{ID: "CPUutilization1", MType: "gauge", Value: post.Gauge(h.Load1)}
	job <- List
}
