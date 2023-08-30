package hash

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/greyfox12/Metrics/internal/agent/hash"
	"github.com/greyfox12/Metrics/internal/server/getparam"
)

type hashWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w hashWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за генерацию ХЕШ, поэтому пишем в него
	return w.Writer.Write(b)
}

func MakeHash(inStr []byte) []byte {
	h := sha256.New()
	h.Write(inStr)
	return h.Sum(nil)
}

func HashHandle(next http.Handler, cfg getparam.ServerParam) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 10000)

		// Проверяю ставить ли подпись
		if cfg.Key != "" && r.Method == http.MethodPost && r.Header.Get("HashSHA256") != "" {
			fmt.Printf("Hash Enter: \n")

			n, err := r.Body.Read(body)
			if err != nil && n <= 0 {
				fmt.Printf("Error req.Body.Read(body):%v: \n", err)
				//	fmt.Printf("n =%v, Body: %v \n", n, body)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			bodyS := body[0:n]
			r.Body = io.NopCloser(bytes.NewBuffer(bodyS))

			hash := hex.EncodeToString(hash.MakeHash(bodyS))
			if r.Header.Get("HashSHA256") != hash {
				fmt.Printf("Error compare HashSHA256: header=%v  hash=%v\n body:%v", r.Header.Get("HashSHA256"), hash, bodyS)
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				fmt.Printf("Hash OK! \n")
			}

			// передаём обработчику страницы переменную типа hashWriter для вывода данных
			next.ServeHTTP(w, r)
			return
		}

		// если hash не поддерживается, передаём управление
		// дальше без изменений
		next.ServeHTTP(w, r)
	})
}

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (mrw *MyResponseWriter) Write(p []byte) (int, error) {
	return mrw.buf.Write(p)
}

func HashWriteHandle(next http.Handler, cfg getparam.ServerParam) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		body := make([]byte, 10000)

		// Проверяю ставить ли подпись
		if cfg.Key != "" {
			fmt.Printf("Hash Write Enter: \n")
			mrw := &MyResponseWriter{
				ResponseWriter: w,
				buf:            &bytes.Buffer{},
			}

			next.ServeHTTP(mrw, r)

			w.Header().Set("HashSHA256", hex.EncodeToString(MakeHash(mrw.buf.Bytes())))
			fmt.Printf("Hash Write OUT buf: %v\n", mrw.buf)
			if _, err := io.Copy(w, mrw.buf); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			//		fmt.Printf("Hash Write OUT: %v\n Hash:%s\n", mrw.buf, hex.EncodeToString(MakeHash(mrw.buf.Bytes())))
			//	w.WriteHeader(http.StatusOK)
			return
		}

		// если ХЕШ не поддерживается, передаём управление
		// дальше без изменений
		next.ServeHTTP(w, r)
	})
}
