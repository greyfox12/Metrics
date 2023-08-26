// Получаю скроку адреса сервера из переменных среды или ключа командной строки

package getparam

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ServerParam struct {
	IPAddress     string
	StoreInterval int
	FileStorePath string
	Restore       bool
	DSN           string
	OnDSN         bool
	OnFile        bool
	OnMap         bool
}

func Param(sp *ServerParam) ServerParam {
	//	var cfg string
	var ok bool
	var tStr string
	var cfg ServerParam
	var err error

	if cfg.IPAddress, ok = os.LookupEnv("ADDRESS"); !ok {
		cfg.IPAddress = sp.IPAddress
	}
	fmt.Printf("LookupEnv(ADDRESS)=%v\n", cfg.IPAddress)

	if tStr, ok = os.LookupEnv("STORE_INTERVAL"); !ok {
		cfg.StoreInterval = sp.StoreInterval
	} else {
		if cfg.StoreInterval, err = strconv.Atoi(tStr); err != nil {
			cfg.StoreInterval = sp.StoreInterval
			fmt.Printf("Error value STORE_INTERVAL=%v", tStr)
		}
	}

	if cfg.FileStorePath, ok = os.LookupEnv("FILE_STORAGE_PATH"); !ok {
		cfg.FileStorePath = sp.FileStorePath
	}

	if tStr, ok = os.LookupEnv("RESTORE"); !ok {
		cfg.Restore = sp.Restore
	} else {
		if cfg.Restore, err = strconv.ParseBool(tStr); err != nil {
			cfg.Restore = sp.Restore
			fmt.Printf("Error value RESTORE=%v", tStr)
		}
	}

	flag.StringVar(&cfg.IPAddress, "a", cfg.IPAddress, "Endpoint server IP address host:port")
	flag.StringVar(&cfg.FileStorePath, "f", cfg.FileStorePath, "File Store Path")
	flag.IntVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "Store interval")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "Restore data from file")
	flag.StringVar(&cfg.DSN, "d", "", "Restore data from file")
	flag.Parse()

	if tStr, ok = os.LookupEnv("DATABASE_DSN"); ok {
		fmt.Printf("LookupEnv(DATABASE_DSN)=%v\n", tStr)
		cfg.DSN = tStr
	}
	//	if cfg.DSN == "" {
	//		cfg.DSN = sp.DSN
	//	}

	// Определяю где храним данные
	cfg.OnDSN = cfg.DSN != ""
	cfg.OnFile = cfg.DSN == "" && cfg.FileStorePath != ""
	cfg.OnMap = cfg.DSN == "" && cfg.FileStorePath == ""

	fmt.Printf("After key (ADDRESS)=%v\n", cfg.IPAddress)
	fmt.Printf("After key (DATABASE_DSN)=%v\n", cfg.DSN)
	fmt.Printf("After key cfg.FileStorePath=%v\n", cfg.FileStorePath)
	fmt.Printf("After key cfg.StoreInterval=%v\n", cfg.StoreInterval)
	fmt.Printf("After key cfg.Restore=%v\n", cfg.Restore)

	fmt.Printf("cfg.OnDSN=%v\n", cfg.OnDSN)
	fmt.Printf("After key cfg.OnFile=%v\n", cfg.OnFile)

	//	fmt.Printf("os.Args=%v\n", os.Args)
	//	fmt.Printf("os.Environ=%v\n", os.Environ())
	return cfg
}
