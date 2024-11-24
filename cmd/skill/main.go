package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type gauge float64
type counter int64
type MemStorage struct {
	gauge   map[string]gauge
	counter map[string]counter
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(webhook))
}

func webhook(rwr http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // разрешаем только POST-запросы
		rwr.WriteHeader(http.StatusMethodNotAllowed)
		rwr.Write([]byte("Only POST method is allowed"))
		return
	}
	rwr.Header().Set("Content-Type", "text/plain")

	body := r.URL.String()
	rwr.Write([]byte(body))

	splisli := strings.Split(body, "/")
	if len(splisli) < 5 {
		rwr.WriteHeader(http.StatusNotFound)
		rwr.Write([]byte("StatusNotFound, man\n"))
		return
	}
	if splisli[1] != "update" || (splisli[2] != "gauge" && splisli[2] != "counter") {
		rwr.WriteHeader(http.StatusBadRequest)
		rwr.Write([]byte("Bad Request, no \"update\"\n"))
		return
	}
	if splisli[2] == "counter" {
		_, err := strconv.ParseInt(splisli[4], 10, 64)
		if err != nil {
			rwr.WriteHeader(http.StatusBadRequest)
			rwr.Write([]byte("Bad Request counter \n"))
			return
		}
	}
	if splisli[2] == "gauge" {
		_, err := strconv.ParseFloat(splisli[4], 64)
		if err != nil {
			rwr.WriteHeader(http.StatusBadRequest)
			rwr.Write([]byte("Bad Request gauge \n"))
			return
		}
	}

	outer := fmt.Sprintf("len %d\n", len(splisli))
	for i, v := range splisli {
		outer += fmt.Sprintf("%d %s\n", i, v)
	}
	rwr.Write([]byte(outer))

}
