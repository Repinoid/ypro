package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrearMetrix(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		mappa map[string]string
		urla string
		want want
	}{
		{
			name: "Right case",
			mappa: map[string]string{"metricType": "gauge", "metricName": "Alloc", "metricValue": "77.77"},
			urla: "/update/gauge/Alloc/77.77",
			want: want{
				code:        http.StatusNotFound,
				response:    `{"status":"StatusNotFound"}`,
				contentType: "text/plain",
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.urla, nil)
			request = mux.SetURLVars(request, tt.mappa)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			badPost(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody), tt.urla)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
