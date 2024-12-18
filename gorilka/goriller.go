package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const localPort = "localhost:8080"

func main() {
	router := mux.NewRouter()
	router.Headers("Content-Type", "application/json")

	router.HandleFunc("/params", params).Methods("POST")
	//	router.HandleFunc("/", treat).Methods("POST")

	if err := http.ListenAndServe(localPort, router); err != nil {
		fmt.Println(err.Error())
	}
}

func params(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Encoding", "gzip")
	rwr.Header().Set("Content-Type", "application/javascript")
	telo, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Accept-Encoding  %v\n", req.Header.Get("Accept-Encoding"))
	fmt.Printf("Content-Encoding  %v\n", req.Header.Get("Content-Encoding"))
	// hdr := req.Header
	// fmt.Fprintf(rwr, "headers %[1]v type %[1]T\n", hdr)
	// hst := req.Host
	// fmt.Fprintf(rwr, "Host %30[1]v type %[1]T\n", hst)
	// frm := req.Form
	// fmt.Fprintf(rwr, "Form %30[1]v type %[1]T\n", frm)
	// resp := req.Response
	// fmt.Fprintf(rwr, "Response %30[1]v type %[1]T\n", resp)
	// meth := req.Method
	// fmt.Fprintf(rwr, "Method %30[1]v type %[1]T\n", meth)
	// requ := req.RequestURI
	// fmt.Fprintf(rwr, "RequestURI %30[1]v type %[1]T\n", requ)
	// url := req.URL
	// fmt.Fprintf(rwr, "URL %30[1]v type %[1]T\n", url)
	// zippers := req.Header.Get("Accept-Encoding")
	// fmt.Fprintf(rwr, "Zippers %30[1]v type %[1]T\n", zippers)
	// fmt.Fprintf(rwr, "Compressed %[1]v type %[1]T\n", buf)

	//	gzip.NewWriter(w io.Writer)
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Name = "zw name"
	zw.Comment = "zw comment"
	zw.ModTime = time.Now()
	_, err = zw.Write(telo)
	if err != nil {
		log.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
	rwr.Write(buf.Bytes())

	//		fmt.Printf("Name: %s\nComment: %s\nModTime: %s\n\n", zr.Name, zr.Comment, zr.ModTime.UTC())

}
