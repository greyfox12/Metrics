package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/greyfox12/Metrics/internal/agent/getparam"
	"github.com/greyfox12/Metrics/internal/agent/post"
)

func main() {

	// Читаю окружение и ключи командной строки
	Config := getparam.Param()
	fmt.Printf("ServerAdr = %v\n", Config.Address)
	fmt.Printf("PollInterval = %v\n", Config.PollInterval)
	fmt.Printf("ReportInterval = %v\n", Config.ReportInterval)

	if Config.PollInterval > Config.ReportInterval {
		panic("ReportInterval должен быть больше PollInterval")
	}

	var ListCounter map[int]post.CounterMetric
	var ListGauge map[int]post.GaugeMetric

	var m runtime.MemStats

	PollCount := post.Counter(0) //Счетчик циклов опроса

	ListGauge = make(map[int]post.GaugeMetric)
	ListCounter = make(map[int]post.CounterMetric)

	client := post.NewClient(Config.Address)

	for {
		runtime.ReadMemStats(&m)

		ListGauge[1] = post.GaugeMetric{Name: "Alloc", Val: post.Gauge(m.Alloc)}
		ListGauge[2] = post.GaugeMetric{Name: "BuckHashSys", Val: post.Gauge(m.BuckHashSys)}
		ListGauge[3] = post.GaugeMetric{Name: "Frees", Val: post.Gauge(m.Frees)}
		ListGauge[4] = post.GaugeMetric{Name: "GCCPUFraction", Val: post.Gauge(m.GCCPUFraction)}
		ListGauge[5] = post.GaugeMetric{Name: "GCSys", Val: post.Gauge(m.GCSys)}
		ListGauge[6] = post.GaugeMetric{Name: "HeapAlloc", Val: post.Gauge(m.HeapAlloc)}
		ListGauge[7] = post.GaugeMetric{Name: "HeapIdle", Val: post.Gauge(m.HeapIdle)}
		ListGauge[8] = post.GaugeMetric{Name: "HeapObjects", Val: post.Gauge(m.HeapObjects)}
		ListGauge[9] = post.GaugeMetric{Name: "HeapReleased", Val: post.Gauge(m.HeapReleased)}
		ListGauge[10] = post.GaugeMetric{Name: "HeapSys", Val: post.Gauge(m.HeapSys)}
		ListGauge[11] = post.GaugeMetric{Name: "LastGC", Val: post.Gauge(m.LastGC)}
		ListGauge[12] = post.GaugeMetric{Name: "Lookups", Val: post.Gauge(m.Lookups)}
		ListGauge[13] = post.GaugeMetric{Name: "MCacheInuse", Val: post.Gauge(m.MCacheInuse)}
		ListGauge[14] = post.GaugeMetric{Name: "MCacheSys", Val: post.Gauge(m.MCacheSys)}
		ListGauge[15] = post.GaugeMetric{Name: "Mallocs", Val: post.Gauge(m.Mallocs)}
		ListGauge[16] = post.GaugeMetric{Name: "NextGC", Val: post.Gauge(m.NextGC)}
		ListGauge[17] = post.GaugeMetric{Name: "NumForcedGC", Val: post.Gauge(m.NumForcedGC)}
		ListGauge[18] = post.GaugeMetric{Name: "NumGC", Val: post.Gauge(m.NumGC)}
		ListGauge[19] = post.GaugeMetric{Name: "OtherSys", Val: post.Gauge(m.OtherSys)}
		ListGauge[20] = post.GaugeMetric{Name: "PauseTotalNs", Val: post.Gauge(m.PauseTotalNs)}
		ListGauge[21] = post.GaugeMetric{Name: "StackInuse", Val: post.Gauge(m.StackInuse)}
		ListGauge[22] = post.GaugeMetric{Name: "StackSys", Val: post.Gauge(m.StackSys)}
		ListGauge[23] = post.GaugeMetric{Name: "Sys", Val: post.Gauge(m.Sys)}
		ListGauge[24] = post.GaugeMetric{Name: "TotalAlloc", Val: post.Gauge(m.TotalAlloc)}
		ListGauge[25] = post.GaugeMetric{Name: "RandomValue", Val: post.Gauge(rand.Float64())}
		ListGauge[26] = post.GaugeMetric{Name: "HeapInuse", Val: post.Gauge(m.HeapInuse)}
		ListGauge[27] = post.GaugeMetric{Name: "MSpanInuse", Val: post.Gauge(m.MSpanInuse)}
		ListGauge[28] = post.GaugeMetric{Name: "MSpanSys", Val: post.Gauge(m.MSpanSys)}

		ListCounter[1] = post.CounterMetric{Name: "PollCount", Val: post.Counter(PollCount)}

		if int(PollCount)%(Config.ReportInterval/Config.PollInterval) == 0 {
			if ok := client.PostCounter(ListGauge, ListCounter, "updates"); ok != nil {
				fmt.Printf("Error Post metrics: %v\n", ok)
			}
		}

		time.Sleep(time.Duration(Config.PollInterval) * time.Second)

		PollCount++
	}

}
