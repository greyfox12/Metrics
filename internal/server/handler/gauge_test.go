package handler

import (
	//	"fmt"
	"github.com/greyfox12/Metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGaugePage(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		send string
		want want
	}{
		{
			name: "positive test #1",
			send: "/update/gauge/test/56",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/plain",
			},
		},
		{
			name: "Error in digital test #2",
			send: "/update/gauge/test/5x6",
			want: want{
				code:        400,
				response:    `{"status":"Bad Request"}`,
				contentType: "text/plain",
			},
		},
		{
			name: "Error no metric test #3",
			send: "/update/gauge/test",
			want: want{
				code:        400,
				response:    `{"status":"Bad Request"}`,
				contentType: "text/plain",
			},
		},
		{
			name: "Error unknow req test #4",
			send: "/updat/gauge/test/56",
			want: want{
				code:        400,
				response:    `{"status":"Bad Request"}`,
				contentType: "text/plain",
			},
		},
		{
			name: "float test #5",
			send: "/update/gauge/test/56.6",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/plain",
			},
		},
	}

	gauge := new(storage.GaugeCounter)
	gauge.Init(100)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.send, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			GaugePage(*gauge, 100).ServeHTTP(w, request)

			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
		})
	}
}
