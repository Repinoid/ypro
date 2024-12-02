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
	//	router.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", treatMetric).Methods("POST")
	//	router.HandleFunc("/value/{metricType}/{metricName}", getMetric).Methods("GET")
	//	router.HandleFunc("/", getAllMetrix).Methods("GET")
	router.HandleFunc("/{ABC}{$}", badPost).Methods("POST")

	return http.ListenAndServe(":8080", router)
}

func badPost(rwr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	//vars := req.PathValue("ABC")
	//metricType := vars["abc"]
	rwr.Header().Set("Content-Type", "text/plain")
	rwr.WriteHeader(http.StatusOK)
	//	fmt.Fprintf(rwr, `{"status":"StatusNotFound"}`)
	//	fmt.Fprintf(rwr, "value - %[1]v, type %[1]T", vars)
	fmt.Fprintf(rwr, "len %d val %s", len(vars), vars["ABC"])

}
