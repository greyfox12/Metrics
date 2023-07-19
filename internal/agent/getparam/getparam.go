package getparam

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DefServerAdr      = "http://localhost:8080"
	DefPollInterval   = 2
	DefReportInterval = 10
)

type TConfig struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

type NetAddress string

func (o *NetAddress) Set(flagValue string) error {
	fmt.Printf("flagValue=%s\n", flagValue)

	if !strings.HasPrefix(flagValue, "http://") {
		*o = NetAddress("http://" + flagValue)
	}
	return nil
}

func (o *NetAddress) String() string {

	return string(*o)
}

func Param() TConfig {

	var cfg TConfig
	if res, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.Address = res
	}

	if tmp, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if res, err := strconv.Atoi(tmp); err == nil {
			cfg.ReportInterval = res
		} else {
			panic(fmt.Sprintf("Неверное значение переменной окружения REPORT_INTERVAL = %v", res))
		}
	}

	if tmp, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if res, err := strconv.Atoi(tmp); err == nil {
			cfg.PollInterval = res
		} else {
			panic(fmt.Sprintf("Неверное значение переменной окружения REPORT_INTERVAL = %v", res))
		}
	}

	//	fmt.Printf("cfg.Address=%s", cfg.Address)

	if cfg.PollInterval == 0 {
		cfg.PollInterval = DefPollInterval
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = DefReportInterval
	}
	if cfg.Address != "" && !strings.HasPrefix(cfg.Address, "http://") {
		cfg.Address = "http://" + cfg.Address
	}
	if cfg.Address == "" {
		cfg.Address = DefServerAdr
	}

	// Ключи командной строки
	ServerAdr := new(NetAddress) // {"http://localhost:8080"}
	_ = flag.Value(ServerAdr)

	// проверка реализации
	flag.Var(ServerAdr, "a", "Net address host:port")

	flag.IntVar(&cfg.PollInterval, "p", cfg.PollInterval, "Pool interval sec.")
	flag.IntVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "Report interval sec.")
	flag.Parse()

	if *ServerAdr != "" {
		cfg.Address = string(*ServerAdr)
	}

	return cfg
}