/*
metricstest -test.v -test.run="^TestIteration10[AB]*$" ^
-binary-path=cmd/server/server.exe -source-path=cmd/server/ ^
-agent-binary-path=cmd/agent/agent.exe ^
-server-port=8080 -file-storage-path=goshran.txt ^
-database-dsn=postgres://postgres:passwordas@localhost:5432/postgres


curl localhost:8080/update/ -H "Content-Type":"application/json" -d "{\"type\":\"counter\",\"id\":\"PollCount\",\"value\":77}"
curl localhost:8080/value/ -H "Content-Type":"application/json" -d "{\"type\":\"counter\",\"id\":\"PollCount\"}"
*/

package main

import (
	"context"
	"gorono/internal/handlera"
	"gorono/internal/memos"
	"gorono/internal/middlas"
	"gorono/internal/models"
	"log"

	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/gorilla/mux"
) // подключаем пакет pprof

// type Metrics = memos.Metrics
type MemStorage = memos.MemoryStorageStruct

var host = "localhost:8080"

var ctx context.Context

// var models.Inter models.Inter // 	= memStor OR dbStorage

func main() {

	if err := InitServer(); err != nil {
		log.Println(err, " no success for foa4Server() ")
		return
	}

	if models.ReStore {
		_ = models.Inter.LoadMS(models.FileStorePath)
	}

	if models.StoreInterval > 0 {
		go models.Inter.Saver(models.FileStorePath, models.StoreInterval)
	}

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	router := mux.NewRouter()
	router.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", handlera.PutMetric).Methods("POST")
	router.HandleFunc("/value/{metricType}/{metricName}", handlera.GetMetric).Methods("GET")

	router.HandleFunc("/update/", handlera.PutJSONMetric).Methods("POST")
	router.HandleFunc("/value/", handlera.GetJSONMetric).Methods("POST")

	router.HandleFunc("/updates/", handlera.Buncheras).Methods("POST")	// for AGENT
	router.HandleFunc("/", handlera.GetAllMetrix).Methods("GET")	// Вывод всех метрик на экран
	router.HandleFunc("/", handlera.BadPost).Methods("POST") // if POST with wrong arguments structure
	router.HandleFunc("/ping", handlera.DBPinger).Methods("GET")

	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	router.Use(middlas.GzipHandleEncoder)
	router.Use(middlas.GzipHandleDecoder)
	router.Use(middlas.WithLogging)
	router.Use(handlera.CryptoHandleDecoder)

	//	return http.ListenAndServe(":8080", router)
	return http.ListenAndServe(host, router)
}

/*
metricstest -test.v -test.run="^TestIteration11[AB]*$" ^
-binary-path=cmd/server/server.exe -source-path=cmd/server/ ^
-agent-binary-path=cmd/agent/agent.exe ^
-server-port=8080 -file-storage-path=goshran.txt ^
-database-dsn=postgres://postgres:passwordas@localhost:5432/postgres


metricstest -test.v -test.run="^TestIteration1[AB]*$" -binary-path=cmd/server/server.exe -source-path=cmd/server/

go run . -d=postgres://postgres:passwordas@localhost:5432/postgres

*/
