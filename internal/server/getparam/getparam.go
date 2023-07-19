// Получаю скроку адреса сервера из переменных среды или ключа командной строки

package getparam

import (
	"flag"
	"os"
)

func Param(defServerAdr string) string {
	var cfg string
	var ok bool

	if cfg, ok = os.LookupEnv("ADDRESS"); !ok {
		cfg = defServerAdr
	}

	IPAddress := flag.String("a", cfg, "Endpoint server IP address host:port")
	flag.Parse()
	return *IPAddress
}
