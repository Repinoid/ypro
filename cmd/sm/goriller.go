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
	// router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	vars := mux.Vars(req)
	// 	fmt.Fprintf(w, "Welcome %v\n", vars)
	// 	fmt.Println("Welcome to the home page!!!!!")
	// }).Methods("GET")
	router.Headers("Content-Type", "text/plain")
	router.HandleFunc("/params", params).Methods("POST")
	//	router.Headers("Content-Type", "application/json")
	//	router.HandleFunc("/", treat).Methods("POST")

	if err := http.ListenAndServe(localPort, router); err != nil {
		fmt.Println(err.Error())
	}
}

func params(rwr http.ResponseWriter, req *http.Request) {
	telo, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(rwr, "string(body) --- >  %v\n", string(telo))
	hdr := req.Header
	fmt.Fprintf(rwr, "headers %[1]v type %[1]T\n", hdr)
	hst := req.Host
	fmt.Fprintf(rwr, "Host %[1]v type %[1]T\n", hst)
	frm := req.Form
	fmt.Fprintf(rwr, "Form %[1]v type %[1]T\n", frm)
	resp := req.Response
	fmt.Fprintf(rwr, "Response %[1]v type %[1]T\n", resp)
	meth := req.Method
	fmt.Fprintf(rwr, "Method %[1]v type %[1]T\n", meth)
	requ := req.RequestURI
	fmt.Fprintf(rwr, "RequestURI %[1]v type %[1]T\n", requ)
	url := req.URL
	fmt.Fprintf(rwr, "URL %[1]v type %[1]T\n", url)

	//fmt.Fprintf(rwr, "%v", ae)
	//	rwr.Write([]byte(ae))

}
