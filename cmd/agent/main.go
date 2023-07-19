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

		ListGauge[1] = post.GaugeMetric{"Alloc", post.Gauge(m.Alloc)}
		ListGauge[2] = post.GaugeMetric{"BuckHashSys", post.Gauge(m.BuckHashSys)}
		ListGauge[3] = post.GaugeMetric{"Frees", post.Gauge(m.Frees)}
		ListGauge[4] = post.GaugeMetric{"GCCPUFraction", post.Gauge(m.GCCPUFraction)}
		ListGauge[5] = post.GaugeMetric{"GCSys", post.Gauge(m.GCSys)}
		ListGauge[6] = post.GaugeMetric{"HeapAlloc", post.Gauge(m.HeapAlloc)}
		ListGauge[7] = post.GaugeMetric{"HeapIdle", post.Gauge(m.HeapIdle)}
		ListGauge[8] = post.GaugeMetric{"HeapObjects", post.Gauge(m.HeapObjects)}
		ListGauge[9] = post.GaugeMetric{"HeapReleased", post.Gauge(m.HeapReleased)}
		ListGauge[10] = post.GaugeMetric{"HeapSys", post.Gauge(m.HeapSys)}
		ListGauge[11] = post.GaugeMetric{"LastGC", post.Gauge(m.LastGC)}
		ListGauge[12] = post.GaugeMetric{"Lookups", post.Gauge(m.Lookups)}
		ListGauge[13] = post.GaugeMetric{"MCacheInuse", post.Gauge(m.MCacheInuse)}
		ListGauge[14] = post.GaugeMetric{"MCacheSys", post.Gauge(m.MCacheSys)}
		ListGauge[15] = post.GaugeMetric{"Mallocs", post.Gauge(m.Mallocs)}
		ListGauge[16] = post.GaugeMetric{"NextGC", post.Gauge(m.NextGC)}
		ListGauge[17] = post.GaugeMetric{"NumForcedGC", post.Gauge(m.NumForcedGC)}
		ListGauge[18] = post.GaugeMetric{"NumGC", post.Gauge(m.NumGC)}
		ListGauge[19] = post.GaugeMetric{"OtherSys", post.Gauge(m.OtherSys)}
		ListGauge[20] = post.GaugeMetric{"PauseTotalNs", post.Gauge(m.PauseTotalNs)}
		ListGauge[21] = post.GaugeMetric{"StackInuse", post.Gauge(m.StackInuse)}
		ListGauge[22] = post.GaugeMetric{"StackSys", post.Gauge(m.StackSys)}
		ListGauge[23] = post.GaugeMetric{"Sys", post.Gauge(m.Sys)}
		ListGauge[24] = post.GaugeMetric{"TotalAlloc", post.Gauge(m.TotalAlloc)}
		ListGauge[25] = post.GaugeMetric{"RandomValue", post.Gauge(rand.Float64())}

		ListCounter[1] = post.CounterMetric{"PollCount", post.Counter(PollCount)}

		if int(PollCount)%(Config.ReportInterval/Config.PollInterval) == 0 {
			if ok := client.PostCounter(ListGauge, ListCounter); ok != nil {
				fmt.Printf("Error Post metrics: %v\n", ok)
			}
		}

		time.Sleep(time.Duration(Config.PollInterval) * time.Second)

		PollCount++
	}

}
