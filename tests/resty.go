package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Metrics struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	ID    string   `json:"id"`              // имя метрики
}

//func main1() {
func main2() {

	httpc := resty.New() //
	httpc.SetBaseURL("http://localhost:8080")
	//	req := httpc.R().
	//	SetHeader("Accept-Encoding", "gzip").
	//		SetHeader("Content-Type", "application/json")
	//	var result Metrics

	//f64 := float64(55)

	req := httpc.R().
		SetHeader("Accept", "text/html").
		SetHeader("Content-Type", "text/html").
		SetHeader("Accept-Encoding", "gzip")

	resp, err := req.
		SetDoNotParseResponse(true).
		Get("/")

	fmt.Printf("%v\n%v\n%v\n", resp.StatusCode(), resp.Header().Get("Content-Encoding"), err)
	/*
		resp, err0 := req.
			SetBody(&Metrics{
				ID:    "Alloc",
				MType: "gauge",
				//Value: &f64,
			}).
			SetResult(&result).
			Post("value/")

		umjs := Metrics{}
		if err0 != nil {
			log.Print(err0)
		}
		log.Print("aaa")

		bod := resp.Body()
		_ = json.Unmarshal(bod, &umjs)
		//	telo, err := io.ReadAll(resp.Body)

		fmt.Printf("%v\n%v\n%v\n", resp.StatusCode(), result, err0)
	*/
	// js := Metrics{MType: "gauge", ID: "Alloc"}
	// mjs, _ := json.Marshal(js)
	// fmt.Printf("%s\n%s\n%v\n%v\n", js, mjs, umjs, err)

}
