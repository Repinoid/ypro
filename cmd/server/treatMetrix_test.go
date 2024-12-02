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
		name  string
		urla  string
		mappa map[string]string
		want  want
	}{
		{
			name:  "No metric value",
			urla:  "/update/gauge/Alloc",
			mappa: map[string]string{"metricType": "gauge", "metricName": "Alloc"},
			want: want{
				code:        http.StatusNotFound,
				response:    `{"status":"StatusNotFound"}`,
				contentType: "text/plain",
			},
		},
		{
			name:  "wrong type",
			urla:  "/update/gaug/Alloc/77.77",
			mappa: map[string]string{"metricType": "gaug", "metricName": "Alloc", "metricValue": "77.77"},
			want: want{
				code:        http.StatusBadRequest,
				response:    `{"status":"StatusBadRequest"}`,
				contentType: "text/plain",
			},
		},
		{
			name:  "wrong value",
			urla:  "/update/gauge/Alloc/77.ee77",
			mappa: map[string]string{"metricType": "gauge", "metricName": "Alloc", "metricValue": "77.ee77"},
			want: want{
				code:        http.StatusBadRequest,
				response:    `{"status":"StatusBadRequest"}`,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.urla, nil)
			w := httptest.NewRecorder()
			request = mux.SetURLVars(request, tt.mappa)	// !!!!!
			treatMetric(w, request)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
