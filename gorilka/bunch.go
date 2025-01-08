package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"internal/memo"
	"io"
	"net/http"
)

type gauge = memo.Gauge
type counter = memo.Counter
type MemStorage = memo.MemStorage

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func bunchas(rwr http.ResponseWriter, req *http.Request) {
	telo, err := io.ReadAll(req.Body)
	//	reader := bytes.NewReader(telo)
	//var err error
	//	telo := []byte
	// if req.Header.Get(`Content-Encoding`) == `gzip` { // if compressed input
	// //	unp, err := unpackFromGzip(reader) // then unpack body
	// 	unp, err := gzip.NewReader(reader)
	// 	if err != nil {
	// 		rwr.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}
	// 	_, err = unp.Read(telo)
	// 	if err != nil {
	// 		rwr.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}
	// }

	//telo, err := io.ReadAll(reader)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest)
		return
	}
	buf := bytes.NewBuffer(telo)
	memor := []Metrics{}
	err = json.NewDecoder(buf).Decode(&memor)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, j := range memor {
		fmt.Printf("%+v\n", j)
	}
	reply := struct{ Dlina int }{Dlina: len(memor)}
	json.NewEncoder(rwr).Encode(reply)
}
