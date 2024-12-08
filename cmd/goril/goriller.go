package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Worker struct {
	First  int `json:"a"`
	Second int `json:"b"`
}

const localPort = "localhost:8080"

func main() {
	router := mux.NewRouter()
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
	rwr.WriteHeader(http.StatusOK)
	b, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var worker Worker
	err1 := json.Unmarshal([]byte(b), &worker)
	if err1 != nil {
		fmt.Printf("Ошибка чтения JSON-данных: err %v  struct %v readall %v", err1, worker, b)
	}
	//json.NewEncoder(rwr).Encode(string(b))
	//json.NewEncoder(rwr).Encode(b)
	fmt.Fprintf(rwr, "as str %[1]s as raw  %[1]v as encoded %[2]v\n", b, worker)
	//fmt.Printf("aaaaaaaaaaaaaaa %s", string(b))

}
