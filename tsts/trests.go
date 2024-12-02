package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	//router := http.NewServeMux()

	router := mux.NewRouter()
	//router.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", treatMetric).Methods("POST")
	//	router.HandleFunc("/value/{metricType}/{metricName}", getMetric).Methods("GET")
	//	router.HandleFunc("/", getAllMetrix).Methods("GET")
	router.HandleFunc("/update/{metricType}/{metricName}/{metricValue}{$}", badPost).Methods("POST")
	//	router.HandleFunc("POST /a{ABC}{$}", badPost) //.Methods("POST")

	return http.ListenAndServe(":8080", router)
}

func badPost(rwr http.ResponseWriter, req *http.Request) {
	req = mux.SetURLVars(req, map[string]string{"metricType": "gauge", "metricName": "Alloc", "metricValue": "77.77"})
	vars := mux.Vars(req)
	//vars := req.PathValue("ABC")
	metricType := vars["metricType"]
	metricName := vars["metricName"]
	metricValue := vars["metricValue"]
	rwr.Header().Set("Content-Type", "text/plain")
	rwr.WriteHeader(http.StatusNotFound)
	//fmt.Fprintf(rwr, `{"status":"StatusNotFound"}`)
	//	fmt.Fprintf(rwr, "value - %[1]v, type %[1]T", vars)
	fmt.Fprintln(rwr, metricType, metricName, metricValue)
}
