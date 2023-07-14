package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

//const (
//	ServerAdr      = "http://localhost:8080"
//	pollInterval   = 2
//	reportInterval = 10
//)

type counter int64
type gauge float64
type CounterMetric struct {
	name string
	Val  counter
}
type GaugeMetric struct {
	name string
	Val  gauge
}

var ListCounter map[int]CounterMetric
var ListGauge map[int]GaugeMetric

// ///////////
type NetAddress string

func (o *NetAddress) Set(flagValue string) error {
	fmt.Printf("flagValue=%s\n", flagValue)

	if !strings.HasPrefix(flagValue, "http://") {
		*o = NetAddress("http://" + flagValue)
	}
	return nil

}

func (o *NetAddress) String() string {
	fmt.Printf("flag\n")
	if *o == "" {
		*o = "http://localhost:8080"
	}
	return string(*o)
}

//////////////

func main() {
	ServerAdr := new(NetAddress) // {"http://localhost:8080"}
	_ = flag.Value(ServerAdr)

	// проверка реализации
	flag.Var(ServerAdr, "a", "Net address host:port")

	//	ServerAdr := flag.String("a", "http://localhost:8080", "Endpoint server IP address host:port")
	pollInterval := flag.Int("p", 2, "Pool interval sec.")
	reportInterval := flag.Int("r", 10, "Report interval sec.")
	flag.Parse()

	fmt.Printf("ServerAdr = %v\n", *ServerAdr)

	var m runtime.MemStats
	PollCount := counter(0) //Счетчик циклов опроса
	RandomValue := gauge(0)
	ListGauge = make(map[int]GaugeMetric)
	ListCounter = make(map[int]CounterMetric)

	client := NewClient(string(*ServerAdr))

	for {
		runtime.ReadMemStats(&m)
		RandomValue = gauge(rand.Float64())
		ListGauge[1] = GaugeMetric{"Alloc", gauge(m.Alloc)}
		ListGauge[2] = GaugeMetric{"BuckHashSys", gauge(m.BuckHashSys)}
		ListGauge[3] = GaugeMetric{"Frees", gauge(m.Frees)}
		ListGauge[4] = GaugeMetric{"GCCPUFraction", gauge(m.GCCPUFraction)}
		ListGauge[5] = GaugeMetric{"GCSys", gauge(m.GCSys)}
		ListGauge[6] = GaugeMetric{"HeapAlloc", gauge(m.HeapAlloc)}
		ListGauge[7] = GaugeMetric{"HeapIdle", gauge(m.HeapIdle)}
		ListGauge[8] = GaugeMetric{"HeapObjects", gauge(m.HeapObjects)}
		ListGauge[9] = GaugeMetric{"HeapReleased", gauge(m.HeapReleased)}
		ListGauge[10] = GaugeMetric{"HeapSys", gauge(m.HeapSys)}
		ListGauge[11] = GaugeMetric{"LastGC", gauge(m.LastGC)}
		ListGauge[12] = GaugeMetric{"Lookups", gauge(m.Lookups)}
		ListGauge[13] = GaugeMetric{"MCacheInuse", gauge(m.MCacheInuse)}
		ListGauge[14] = GaugeMetric{"MCacheSys", gauge(m.MCacheSys)}
		ListGauge[15] = GaugeMetric{"Mallocs", gauge(m.Mallocs)}
		ListGauge[16] = GaugeMetric{"NextGC", gauge(m.NextGC)}
		ListGauge[17] = GaugeMetric{"NumForcedGC", gauge(m.NumForcedGC)}
		ListGauge[18] = GaugeMetric{"NumGC", gauge(m.NumGC)}
		ListGauge[19] = GaugeMetric{"OtherSys", gauge(m.OtherSys)}
		ListGauge[20] = GaugeMetric{"PauseTotalNs", gauge(m.PauseTotalNs)}
		ListGauge[21] = GaugeMetric{"StackInuse", gauge(m.StackInuse)}
		ListGauge[22] = GaugeMetric{"StackSys", gauge(m.StackSys)}
		ListGauge[23] = GaugeMetric{"Sys", gauge(m.Sys)}
		ListGauge[24] = GaugeMetric{"TotalAlloc", gauge(m.TotalAlloc)}
		ListGauge[25] = GaugeMetric{"RandomValue", gauge(RandomValue)}

		ListCounter[1] = CounterMetric{"PollCount", counter(PollCount)}

		if int(PollCount)%(*reportInterval / *pollInterval) == 0 {
			_ = client.PostCounter(ListGauge, ListCounter)
		}

		time.Sleep(time.Duration(*pollInterval) * time.Second)

		PollCount++
	}

}
