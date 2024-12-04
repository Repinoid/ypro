package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type gauge float64
type counter int64
type MemStorage struct {
	gau   map[string]gauge
	count map[string]counter
}

var memStor MemStorage
var host = "localhost:8080"

func main() {
	if faa4server() != 0 {
		return
	}
	memStor = newMemStorage()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	router := mux.NewRouter()
	router.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", treatMetric).Methods("POST")
	router.HandleFunc("/value/{metricType}/{metricName}", getMetric).Methods("GET")
	router.HandleFunc("/", getAllMetrix).Methods("GET")
	router.HandleFunc("/", badPost).Methods("POST") // if POST with wrong arguments structure

	return http.ListenAndServe(host, router)
}

func badPost(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "text/plain")
	rwr.WriteHeader(http.StatusNotFound)
	//	fmt.Fprintf(rwr, "POST http.StatusNotFound with %s\n", req.URL.Path)
	fmt.Fprintf(rwr, `{"status":"StatusNotFound"}`)
}

func getAllMetrix(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "text/plain")
	if req.URL.Path != "/" { // if GET with wrong arguments structure
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		return
	}
	rwr.WriteHeader(http.StatusOK)
	for nam, val := range memStor.gau {
		flo := strconv.FormatFloat(float64(val), 'f', -1, 64) // -1 - to remove zeroes tail
		fmt.Fprintf(rwr, "Gauge Metric name   %20s\t\tvalue\t%s\n", nam, flo)
	}
	for nam, val := range memStor.count {
		fmt.Fprintf(rwr, "Counter Metric name %20s\t\tvalue\t%d\n", nam, val)
	}
}
func getMetric(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "text/plain")
	vars := mux.Vars(req)
	val := "badly"                // does not matter what initial value, could be "var val string"
	status := http.StatusNotFound // this remains if getGaugeValue or getCounterValue don't work
	metricType := vars["metricType"]
	metricName := vars["metricName"]
	if metricType == "gauge" {
		status = memStor.getGaugeValue(metricName, &val)
	}
	if metricType == "counter" {
		status = memStor.getCounterValue(metricName, &val)
	}
	if status == http.StatusOK {
		rwr.WriteHeader(http.StatusOK)
		fmt.Fprint(rwr, val)
	} else {
		rwr.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rwr, "BadRequest, No value for %s of %s type\n", metricName, metricType)
	}
}

func treatMetric(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "text/plain")
	vars := mux.Vars(req)
	rwr.Header().Set("Content-Type", "text/plain")
	metricType := vars["metricType"]
	metricName := vars["metricName"]
	metricValue := vars["metricValue"]
	if metricValue == "" {
		rwr.WriteHeader(http.StatusNotFound)
		return
	}
	if metricType != "gauge" && metricType != "counter" {
		rwr.WriteHeader(http.StatusBadRequest)
		return
	}
	if metricType == "counter" {
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			rwr.WriteHeader(http.StatusBadRequest)
			return
		}
		memStor.addCounter(metricName, counter(value))
	} else { //	if metricType == "gauge" {
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			rwr.WriteHeader(http.StatusBadRequest)
			return
		}
		memStor.addGauge(metricName, gauge(value))
	}
	rwr.WriteHeader(http.StatusOK)
}
