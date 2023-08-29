package post

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/greyfox12/Metrics/internal/agent/compress"
	"github.com/greyfox12/Metrics/internal/agent/getparam"
	"github.com/greyfox12/Metrics/internal/agent/hash"
	"github.com/greyfox12/Metrics/internal/agent/logmy"
)

type Counter int64
type Gauge float64

type CounterMetric struct {
	Name string
	Val  Counter
}

type GaugeMetric struct {
	Name string
	Val  Gauge
}

type Client struct {
	cfg getparam.TConfig
}

func NewClient(cfg getparam.TConfig) Client {
	return Client{cfg}
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type CollectMetr struct {
	ID    string  // имя метрики
	MType string  // параметр, принимающий значение gauge или counter
	Delta Counter // значение метрики в случае передачи counter
	Value Gauge   // значение метрики в случае передачи gauge
}

func (c Client) PostCounter(job chan map[int]CollectMetr, result chan<- error, updateTyp string) {
	//	fmt.Printf("Time: %v\n", time.Now().Unix())
	stArr := make([]Metrics, 1000)
	cn := 0
	var st Metrics

	adrstr := fmt.Sprintf("%s/%s", c.cfg.Address, updateTyp)

	for dt := range job {
		for _, val := range dt {
			if val.MType == "gauge" {
				st = Metrics{ID: val.ID, MType: "gauge", Value: (*float64)(&val.Value)}
			}
			if val.MType == "counter" {
				st = Metrics{ID: val.ID, MType: "counter", Delta: (*int64)(&val.Delta)}
			}

			stArr[cn] = st
			cn++

			if updateTyp == "update" {
				if ok := postMess(st, adrstr, c.cfg); ok != nil {
					fmt.Printf("Error post: %v, %v\n", st, ok)
				}
			}
		}

		if updateTyp == "update" {
			result <- nil
			return
		}

		if ok := postUpdates(stArr[0:cn], adrstr, c.cfg); ok != nil {

			fmt.Printf("Error posts:  %v\n", ok)
			if _, yes := ok.(net.Error); yes {
				result <- ok
				return
			}
		}

		result <- nil
		//		return
	}
}

// Вывод по одной записи
func postMess(st Metrics, adrstr string, cfg getparam.TConfig) error {
	jsonData, err := json.Marshal(st)
	if err != nil {
		return error(err)
	}

	err = Resend(jsonData, adrstr, cfg)
	if err != nil {
		return error(err)
	}
	//	fmt.Println("response Body:", body)
	return nil
}

// Вывод слайса
func postUpdates(stArr []Metrics, adrstr string, cfg getparam.TConfig) error {
	var err error

	jsonData, err := json.Marshal(stArr)
	if err != nil {
		return error(err)
	}

	err = Resend(jsonData, adrstr, cfg)
	if err != nil {
		return error(err)
	}
	//	fmt.Println("response Body:", body)
	return nil
}

// Повторяю при ошибках вывод
func Resend(buf []byte, adrstr string, cfg getparam.TConfig) error {
	var err error

	for i := 1; i <= 4; i++ {
		if i > 1 {
			fmt.Printf("Pause: %v sec\n", WaitSec(i-1))
			time.Sleep(time.Duration(WaitSec(i-1)) * time.Second)
		}

		if err = ActPost(buf, adrstr, cfg); err == nil {
			return nil
		}

		logmy.OutLog(fmt.Errorf("post send message: %w", err))
		if _, yes := err.(net.Error); !yes {
			return err
		}
	}
	return nil
}

// Отправить JSON
func ActPost(buf []byte, adrstr string, cfg getparam.TConfig) error {

	//hashIn := base64.URLEncoding.EncodeToString(hash.MakeHash(buf))

	jsonZip, err := compress.Compress(buf)

	if err != nil {
		return error(err)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", adrstr, bytes.NewBuffer(jsonZip))
	if err != nil {
		return error(err)
	}

	if cfg.Key != "" {
		hashIn := hex.EncodeToString(hash.MakeHash(buf))
		//		fmt.Printf("HashIn: %v\n", hashIn)
		req.Header.Set("HashSHA256", hashIn)
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return error(err)
	}

	fmt.Printf("Head response: %v\n", response.Header)
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return error(err)
	}
	fmt.Println("response Body:", string(body))

	if response.Header.Get("Content-Encoding") == "gzip" || response.Header.Get("Content-Encoding") == "flate" {
		fmt.Printf("Header gzip \n")
		body, err = compress.Decompress(body, "flate")
		if err != nil {
			return error(err)
		}
	}

	// Проверяю Подпись ответа, если есть
	if cfg.Key != "" && response.Header.Get("HashSHA256") != "" {
		hashOut := hex.EncodeToString(hash.MakeHash(body))
		if response.Header.Get("HashSHA256") == hashOut {
			fmt.Printf("Check hash Server: OK!\n")
		} else {
			fmt.Printf("Check hash Server: different! %v  %v\n", response.Header.Get("HashSHA256"), hashOut)
			//		return fmt.Errorf("check hash server: different! %v  %v\n", response.Header.Get("HashSHA256"), hashOut)
		}
	}

	logmy.OutLog(fmt.Errorf("post send response body: %v", string(body)))
	fmt.Println("response Body:", string(body))
	return nil
}

// Считаю задержку - по номеру повторения возвращаю длительность в сек
func WaitSec(period int) int {
	switch period {
	case 1:
		return 1
	case 2:
		return 3
	case 3:
		return 5
	default:
		return 0
	}
}
