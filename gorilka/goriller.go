package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const localPort = "localhost:8080"

func main() {
	router := mux.NewRouter()
	router.Headers("Content-Type", "application/json")

	router.HandleFunc("/params", params).Methods("POST")

	if err := http.ListenAndServe(localPort, router); err != nil {
		fmt.Println(err.Error())
	}
}

func pack2gzip(data2pack []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.ModTime = time.Now()
	_, err := zw.Write(data2pack)
	if err != nil {
		return nil, fmt.Errorf("gzip.NewWriter.Write %w ", err)
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("gzip.NewWriter.Close %w ", err)
	}
	return buf.Bytes(), nil
}

func params(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Encoding", "gzip")
	rwr.Header().Set("Content-Type", "application/javascript")
	telo, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("Accept-Encoding  %v\n", req.Header.Get("Accept-Encoding"))
	//fmt.Printf("Content-Encoding  %v\n", req.Header.Get("Content-Encoding"))

	compressed, err := pack2gzip(telo)
	if err != nil {
		panic(err)
	}

	rwr.Write(compressed)

	//		fmt.Printf("Name: %s\nComment: %s\nModTime: %s\n\n", zr.Name, zr.Comment, zr.ModTime.UTC())

}
