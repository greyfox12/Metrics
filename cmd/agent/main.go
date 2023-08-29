package main

import (
	"fmt"
	"time"

	"github.com/greyfox12/Metrics/internal/agent/collectmetric"
	"github.com/greyfox12/Metrics/internal/agent/getparam"
	"github.com/greyfox12/Metrics/internal/agent/logmy"
	"github.com/greyfox12/Metrics/internal/agent/post"
)

func main() {

	// Инициализирую логирование
	if ok := logmy.Initialize("info"); ok != nil {
		panic(ok)
	}

	// Читаю окружение и ключи командной строки
	Config := getparam.Param()
	fmt.Printf("ServerAdr = %v\n", Config.Address)
	fmt.Printf("PollInterval = %v\n", Config.PollInterval)
	fmt.Printf("ReportInterval = %v\n", Config.ReportInterval)
	fmt.Printf("Key = %v\n", Config.Key)

	if Config.PollInterval > Config.ReportInterval {
		panic("ReportInterval должен быть больше PollInterval")
	}

	PollCount := post.Counter(0) //Счетчик циклов опроса

	//	ListGauge = make(map[int]post.GaugeMetric)
	//	ListCounter = make(map[int]post.CounterMetric)

	client := post.NewClient(Config)

	// создаем буферизованный канал для принятия задач в воркер
	jobs := make(chan map[int]post.CollectMetr, Config.RateLimit)
	// создаем буферизованный канал для отправки результатов
	results := make(chan error, Config.RateLimit)

	// Запускаю исполнителей
	for w := 1; w <= Config.RateLimit; w++ {
		go client.PostCounter(jobs, results, "updates")
	}

	for {
		if int(PollCount)%(Config.ReportInterval/Config.PollInterval) == 0 {

			go collectmetric.CollectGauge(int64(PollCount), jobs)
			go collectmetric.CollectAdd(int64(PollCount), jobs)
			time.Sleep(time.Duration(Config.PollInterval) * time.Second)
		}

		PollCount++
		fmt.Printf("PollCount: %v\n", PollCount)

		// Проверяю резудьтат выполнения
		go func(results chan error) {
			for res := range results {
				fmt.Printf("results:\n")
				if res != nil {
					logmy.OutLog(fmt.Errorf("post metrics: %w", res))
				}
			}
		}(results)

	}

}
