package handler

import (
	//	"fmt"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/greyfox12/Metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type tMetrics struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func TestGaugePage(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		send string
		data tMetrics
		want want
	}{
		{
			name: "positive test #1",
			send: "/update",
			data: tMetrics{
				ID:    "test",
				MType: "gauge",
				Value: 56},
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},

		{
			name: "Error no metric test #3",
			send: "/update",
			data: tMetrics{
				ID:    "test",
				MType: "gauge",
			},
			want: want{
				code:        400,
				response:    `{"status":"Bad Request"}`,
				contentType: "application/json",
			},
		},
		{
			name: "Error unknow req test #4",
			send: "/updat",
			data: tMetrics{
				ID:    "test",
				MType: "gauge",
				Value: 56},
			want: want{
				code:        400,
				response:    `{"status":"Bad Request"}`,
				contentType: "application/json",
			},
		},
		{
			name: "float test #5",
			send: "/update",
			data: tMetrics{
				ID:    "test",
				MType: "gauge",
				Value: 56.6},
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
	}

	gauge := new(storage.GaugeCounter)
	gauge.Init(100)
	metric := new(storage.MetricCounter)
	metric.Init(100)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(test.data)
			fmt.Printf("Nest: %v\n", test.name)
			//			fmt.Printf("JSON: %v\n", jsonData)

			request := httptest.NewRequest(http.MethodPost, test.send, bytes.NewBuffer(jsonData))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			//			request.Body = test.json
			PostPage(gauge, metric, 100).ServeHTTP(w, request)

			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			var err error
			_, err = io.ReadAll(res.Body)

			require.NoError(t, err)
		})
	}
}
