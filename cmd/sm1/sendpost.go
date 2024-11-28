package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	//"strings"
)

func postMetric(metricType, metricName, metricValue string) error {
	host := "http://localhost:8080"

	client := &http.Client{}
	payloadStr := "/update/" + metricType + "/" + metricName + "/" + metricValue
	//	payload := strings.NewReader(payloadStr)
	reader := bytes.NewReader([]byte(payloadStr))
	req, err := http.NewRequest(http.MethodPost, host, reader)
	//req, err := http.NewRequest(http.MethodPost, host+payloadStr, nil)
	fmt.Printf("plstr %[1]v type %[1]T\n", payloadStr)

	if err != nil {
		fmt.Println("NewRequest ", err)
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Add("Accept", "text/html")
	req.Header.Add("cache-control", "no-cache")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do ", err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("io.ReadAll ", err)
		return err
	}
	fmt.Print("body ...", string(body))

	/*	resp, err := http.Post(url, "text/plain", nil) //strings.NewReader(url))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		var p []byte
		resp.Body.Read(p)
		fmt.Printf("body\t%v\nheader\t%v\n", p, resp.Header)
		return err */
	return nil
}
func main() {
	stat := postMetric("gauge", "Alloc", "55.66")
	fmt.Println("status ", stat)
	//		panic(err)

}
