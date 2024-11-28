package main

import (
	"fmt"
	"net/http"
)

const localPort = ":8080"

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	router.HandleFunc("POST /", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Zero posterizan\n")
	})
	router.HandleFunc("POST /update/{metrixType}/{metrixName}/{metrixValue}", func(w http.ResponseWriter, req *http.Request) {
		req.Header.Set("Content-Type", "text/plain")
		//	req.Header.Set("Content-Type", "application/json")
		metrixType := req.PathValue("metrixType")
		metrixName := req.PathValue("metrixName")
		metrixValue := req.PathValue("metrixValue")

		fmt.Fprintf(w, "Metr type %s name %s value %s\n", metrixType, metrixName, metrixValue)
		fmt.Println(metrixType, metrixName, metrixValue)
	})

	if err := http.ListenAndServe(localPort, router); err != nil {
		fmt.Println(err.Error())
	}
}
