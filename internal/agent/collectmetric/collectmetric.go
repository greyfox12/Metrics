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

	list := make(map[int]post.CollectMetr)

	runtime.ReadMemStats(&m)
	fmt.Printf("Start %v\n", PollCount)
	list[1] = post.CollectMetr{ID: "Alloc", MType: "gauge", Value: post.Gauge(m.Alloc)}
	list[2] = post.CollectMetr{ID: "BuckHashSys", MType: "gauge", Value: post.Gauge(m.BuckHashSys)}
	list[3] = post.CollectMetr{ID: "Frees", MType: "gauge", Value: post.Gauge(m.Frees)}
	list[4] = post.CollectMetr{ID: "GCCPUFraction", MType: "gauge", Value: post.Gauge(m.GCCPUFraction)}
	list[5] = post.CollectMetr{ID: "GCSys", MType: "gauge", Value: post.Gauge(m.GCSys)}
	list[6] = post.CollectMetr{ID: "HeapAlloc", MType: "gauge", Value: post.Gauge(m.HeapAlloc)}
	list[7] = post.CollectMetr{ID: "HeapIdle", MType: "gauge", Value: post.Gauge(m.HeapIdle)}
	list[8] = post.CollectMetr{ID: "HeapObjects", MType: "gauge", Value: post.Gauge(m.HeapObjects)}
	list[9] = post.CollectMetr{ID: "HeapReleased", MType: "gauge", Value: post.Gauge(m.HeapReleased)}
	list[10] = post.CollectMetr{ID: "HeapSys", MType: "gauge", Value: post.Gauge(m.HeapSys)}
	list[11] = post.CollectMetr{ID: "LastGC", MType: "gauge", Value: post.Gauge(m.LastGC)}
	list[12] = post.CollectMetr{ID: "Lookups", MType: "gauge", Value: post.Gauge(m.Lookups)}
	list[13] = post.CollectMetr{ID: "MCacheInuse", MType: "gauge", Value: post.Gauge(m.MCacheInuse)}
	list[14] = post.CollectMetr{ID: "MCacheSys", MType: "gauge", Value: post.Gauge(m.MCacheSys)}
	list[15] = post.CollectMetr{ID: "Mallocs", MType: "gauge", Value: post.Gauge(m.Mallocs)}
	list[16] = post.CollectMetr{ID: "NextGC", MType: "gauge", Value: post.Gauge(m.NextGC)}
	list[17] = post.CollectMetr{ID: "NumForcedGC", MType: "gauge", Value: post.Gauge(m.NumForcedGC)}
	list[18] = post.CollectMetr{ID: "NumGC", MType: "gauge", Value: post.Gauge(m.NumGC)}
	list[19] = post.CollectMetr{ID: "OtherSys", MType: "gauge", Value: post.Gauge(m.OtherSys)}
	list[20] = post.CollectMetr{ID: "PauseTotalNs", MType: "gauge", Value: post.Gauge(m.PauseTotalNs)}
	list[21] = post.CollectMetr{ID: "StackInuse", MType: "gauge", Value: post.Gauge(m.StackInuse)}
	list[22] = post.CollectMetr{ID: "StackSys", MType: "gauge", Value: post.Gauge(m.StackSys)}
	list[23] = post.CollectMetr{ID: "Sys", MType: "gauge", Value: post.Gauge(m.Sys)}
	list[24] = post.CollectMetr{ID: "TotalAlloc", MType: "gauge", Value: post.Gauge(m.TotalAlloc)}
	list[25] = post.CollectMetr{ID: "RandomValue", MType: "gauge", Value: post.Gauge(rand.Float64())}
	list[26] = post.CollectMetr{ID: "HeapInuse", MType: "gauge", Value: post.Gauge(m.HeapInuse)}
	list[27] = post.CollectMetr{ID: "MSpanInuse", MType: "gauge", Value: post.Gauge(m.MSpanInuse)}
	list[28] = post.CollectMetr{ID: "MSpanSys", MType: "gauge", Value: post.Gauge(m.MSpanSys)}

	list[29] = post.CollectMetr{ID: "PollCount", MType: "gauge", Delta: post.Counter(PollCount)}

	job <- list
	//	return list
}

func CollectAdd(PollCount int64, job chan<- map[int]post.CollectMetr) {

	//	var list map[int]post.CollectMetr
	list := make(map[int]post.CollectMetr)

	fmt.Printf("Start Add:%v\n", PollCount)
	v, _ := mem.VirtualMemory()
	h, _ := load.Avg()
	list[1] = post.CollectMetr{ID: "TotalMemory", MType: "gauge", Value: post.Gauge(v.Total)}
	list[2] = post.CollectMetr{ID: "FreeMemory", MType: "gauge", Value: post.Gauge(v.Free)}
	list[3] = post.CollectMetr{ID: "CPUutilization1", MType: "gauge", Value: post.Gauge(h.Load1)}
	job <- list
}
