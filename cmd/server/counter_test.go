package main

import (
	//	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCounterPage(t *testing.T) {
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
			send: "/update/counter/test/56",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/plain",
			},
		},
		{
			name: "Error in digital test #2",
			send: "/update/counter/test/5x6",
			want: want{
				code:        400,
				response:    `{"status":"Bad Request"}`,
				contentType: "text/plain",
			},
		},
		{
			name: "Error no metric test #3",
			send: "/update/counter/test",
			want: want{
				code:        404,
				response:    `{"status":"Bad Request"}`,
				contentType: "text/plain",
			},
		},
		{
			name: "Error unknow req test #4",
			send: "/updat/counter/test/56",
			want: want{
				code:        400,
				response:    `{"status":"Bad Request"}`,
				contentType: "text/plain",
			},
		},
		{
			name: "float test #5",
			send: "/update/counter/test/56.6",
			want: want{
				code:        400,
				response:    `{"status":"Bad Request"}`,
				contentType: "text/plain",
			},
		},
	}
	//	var b bytes.Buffer

	for _, test := range tests {
		//		fmt.Printf("%s\n", test.send)
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.send, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			CounterPage(w, request)

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
