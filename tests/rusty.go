package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

type Metrics struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	ID    string   `json:"id"`              // имя метрики
}

func main() {

	httpc := resty.New() //
	httpc.SetBaseURL("http://localhost:8080")
	req := httpc.R().
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")
	var result Metrics
	//	f64 := float64(55)
	resp, err := req.
		SetBody(&Metrics{
			ID:    "namer",
			MType: "gauge",
			//			Value: &f64,
		}).
		//		SetResult(&result).
		Post("/value/")
	req.SetResult(&result)

	umjs := Metrics{}
	if err != nil {
		log.Print("gggg")
	}
	log.Print("aaa")

	bod := resp.Body()
	err = json.Unmarshal(bod, &umjs)
	//	telo, err := io.ReadAll(resp.Body)

	fmt.Printf("%v\n%v\n%v\n", resp.StatusCode(), string(bod), err)

	// js := Metrics{MType: "gauge", ID: "Alloc"}

	// mjs, _ := json.Marshal(js)

	// fmt.Printf("%s\n%s\n%v\n%v\n", js, mjs, umjs, err)

}
