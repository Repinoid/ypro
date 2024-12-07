package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

const localPort = "localhost:8080"

func main() {
	router := mux.NewRouter()
	/* router.Headers("Content-Type", "application/json")
	router.Headers("Access-Control-Allow-Origin", "*")
	router.Headers("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	router.Headers("Access-Control-Allow-Headers", "Content-Type")
	*/
	router.HandleFunc("/s", treat).Methods("GET", "OPTIONS")
	router.HandleFunc("/s", treat).Methods("POST")

	if err := http.ListenAndServe(localPort, router); err != nil {
		fmt.Println("ёпта ...", err.Error())
	}
}

func treat(rwr http.ResponseWriter, req *http.Request) {

	rwr.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rwr.Header().Set("Content-Type", "application/json")
	rwr.Header().Set("Access-Control-Allow-Origin", "*")
	//rwr.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	//rwr.Header().Set("Access-Control-Allow-Credentials", "true")
	rwr.WriteHeader(http.StatusOK)
	json.NewEncoder(rwr).Encode("OKOK")
	b, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(rwr, string(b))
	//fmt.Printf("aaaaaaaaaaaaaaa %s", string(b))

}
