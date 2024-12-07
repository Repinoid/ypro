package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

const localPort = "localhost:8080"

func main() {
	router := mux.NewRouter()
	router.Headers("Content-Type", "application/json")
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Fprintf(w, "Welcome %v\n", vars)
		fmt.Println("Welcome to the home page!!!!!")
	}).Methods("GET")
	router.HandleFunc("/s", treat).Methods("POST")
	router.HandleFunc("/", treat).Methods("POST")

	if err := http.ListenAndServe(localPort, router); err != nil {
		fmt.Println(err.Error())
	}
}

func treat(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")
	rwr.WriteHeader(http.StatusOK)
	b, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(rwr, "body %v\n", string(b))
	fmt.Printf("%s", b)

}
