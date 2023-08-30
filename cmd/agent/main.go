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
	config := getparam.Param()
	fmt.Printf("ServerAdr = %v\n", config.Address)
	fmt.Printf("PollInterval = %v\n", config.PollInterval)
	fmt.Printf("ReportInterval = %v\n", config.ReportInterval)
	fmt.Printf("Key = %v\n", config.Key)

	if config.PollInterval > config.ReportInterval {
		panic("ReportInterval должен быть больше PollInterval")
	}

	PollCount := post.Counter(0) //Счетчик циклов опроса

	client := post.NewClient(config)

	// создаем буферизованный канал для принятия задач в воркер
	jobs := make(chan map[int]post.CollectMetr, config.RateLimit)
	// создаем буферизованный канал для отправки результатов
	results := make(chan error, config.RateLimit)

	// Запускаю исполнителей
	for w := 1; w <= config.RateLimit; w++ {
		go client.PostCounter(jobs, results, "updates")
	}

	for {
		if int(PollCount)%(config.ReportInterval/config.PollInterval) == 0 {

			go collectmetric.CollectGauge(int64(PollCount), jobs)
			go collectmetric.CollectAdd(int64(PollCount), jobs)
			time.Sleep(time.Duration(config.PollInterval) * time.Second)
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
