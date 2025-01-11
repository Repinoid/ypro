package main

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Metrics struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	ID    string   `json:"id"`              // имя метрики
}
var delays = []int{1,3,5}
// func main1() {
func main() {

	httpc := resty.New() //
	httpc.SetBaseURL("http://localhost:8080")
	httpc.SetRetryCount(len(delays))
	delays = delays[1:]
	httpc.SetRetryWaitTime(1 * time.Second)    // начальное время повтора
	httpc.SetRetryMaxWaitTime(9 * time.Second) // 1+3+5
	tn := time.Now()                           // -------------
	httpc.SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
		rwt := client.RetryWaitTime
		fmt.Printf("waittime \t%+v\t time %+v  count %d\n", rwt, time.Since(tn), client.RetryCount) // -------
		client.SetRetryWaitTime(time.Duration(delays[0])*time.Second)
		if len(delays) > 1 {
			delays = delays[1:]
		}
	//	client.SetRetryWaitTime(rwt + 2*time.Second)
		tn = time.Now() // ----------------
		return client.RetryWaitTime, nil
	})
	req := httpc.R().
		SetHeader("Accept", "text/html").
		SetHeader("Content-Type", "text/html").
		SetHeader("Accept-Encoding", "gzip")

		//	resp, err :=
	req.
		SetDoNotParseResponse(false).
		Get("/")

	//	fmt.Printf("response - %+v\nerror %v\n", resp.Header(), err)
}
