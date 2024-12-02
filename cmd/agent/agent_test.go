package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostMetric(t *testing.T) {
	tsts := []struct { // добавляем слайс тестов
		name                string
		mtype, mname, value string
		want                int
	}{
		{
			name:  "simple test #1",                        // описываем каждый тест:
			mtype: "gauge", mname: "Alloc", value: "55.55", // значения, которые будет принимать функция,
			want: http.StatusOK, // и ожидаемый результат
		},
		{
			name:  "simple test #1",                       // описываем каждый тест:
			mtype: "gaug", mname: "Alloc", value: "55.55", // значения, которые будет принимать функция,
			want: http.StatusBadRequest, // и ожидаемый результат
		},
		{
			name:  "simple test #1",                          // описываем каждый тест:
			mtype: "gauge", mname: "Alloc", value: "55ff.55", // значения, которые будет принимать функция,
			want: http.StatusBadRequest, // и ожидаемый результат
		},
	}

	for _, test := range tsts {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, postMetric(test.mtype, test.mname, test.value), test.want)
		})
	}
}
