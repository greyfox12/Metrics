package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
)

// Compress сжимает слайс байт.
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	// создаём переменную w — в неё будут записываться входящие данные,
	// которые будут сжиматься и сохраняться в bytes.Buffer
	w, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %w", err)
	}
	// запись данных
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %w", err)
	}
	// обязательно нужно вызвать метод Close() — в противном случае часть данных
	// может не записаться в буфер b; если нужно выгрузить все упакованные данные
	// в какой-то момент сжатия, используйте метод Flush()
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %w", err)
	}
	// переменная b содержит сжатые данные
	return b.Bytes(), nil
}

// Decompress распаковывает слайс байт.
func Decompress(data []byte, typCompress string) ([]byte, error) {
	var r io.Reader
	var err error
	switch typCompress {
	case "flate", "gzip":
		// переменная r будет читать входящие данные и распаковывать их
		//r = flate.NewReader(bytes.NewReader(data))
		//		defer r.Close()
		//	case "gzip":
		r, err = gzip.NewReader(bytes.NewReader(data))
		//	defer r.Close()
		if err != nil {
			return nil, fmt.Errorf("failed gzip.NewReader: %v", err)
		}
	default:
		return nil, fmt.Errorf("failed decompress type: %v", typCompress)
	}

	var b bytes.Buffer
	// в переменную b записываются распакованные данные
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %w", err)
	}

	return b.Bytes(), nil
}
