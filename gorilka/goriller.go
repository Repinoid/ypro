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
func unpackFromGzip(data2unpack io.Reader) (io.Reader, error) {
	gzipReader, err := gzip.NewReader(data2unpack)
	if err != nil {
		return nil, fmt.Errorf("gzip.NewReader %w ", err)
	}
	if err := gzipReader.Close(); err != nil {
		return nil, fmt.Errorf("zr.Close %w ", err)
	}
	return gzipReader, nil
}

func params(rwr http.ResponseWriter, req *http.Request) {
	var reader io.Reader
	var err error
	if req.Header.Get(`Content-Encoding`) == `gzip` {
		reader, err = unpackFromGzip(req.Body)
		if err != nil {
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
		}
	} else {
		reader = req.Body
	}
	telo, err := io.ReadAll(reader)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Header.Get(`Accept-Encoding`) == `gzip` {
		rwr.Header().Set("Content-Encoding", "gzip")
		telo, err = pack2gzip(telo)
	}
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	//rwr.Header().Set("Content-Type", "application/javascript")

	rwr.Write(telo)

}
