package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const localPort = ":8080"

func main() {
	//	router := http.NewServeMux()
	router := mux.NewRouter().StrictSlash(true)
	//	router.Headers("Content-Type", "text/plain")
	//router.Use(mux.CORSMethodMiddleware(r))
	fmt.Println("\x2FA")
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Fprintf(w, "Welcome %v\n", vars)
		fmt.Println("Welcome to the home page!!!!!")
	}).Methods("GET")
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Zero posterizan\n")
	}).Methods("POST")
	//router.HandleFunc("/u", func(w http.ResponseWriter, req *http.Request) {
	//w.WriteHeader(http.StatusCreated)
	router.HandleFunc("/update/{metrixType}/{metrixName}/{metrixValue}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		req.Header.Set("Content-Type", "text/plain")
		metrixType := vars["metrixType"]
		metrixName := vars["metrixName"]
		metrixValue := vars["metrixValue"]

		fmt.Fprintf(w, "Metr type %s name %s value %s\n", metrixType, metrixName, metrixValue)
		fmt.Println(metrixType, metrixName, metrixValue)
	}).Methods("POST")

	if err := http.ListenAndServe(localPort, router); err != nil {
		fmt.Println(err.Error())
	}
}
