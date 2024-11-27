package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func postMetric(metricType, metricName, metricValue string) error {
	host := "http://localhost:8080"
	//	url := "/update/" + metricType + "/" + metricName + "/" + metricValue

	client := &http.Client{}
	payloadStr := "/update/" + metricType + "/" + metricName + "/" + metricValue
	payload := strings.NewReader(payloadStr)
	req, err := http.NewRequest("POST", host, payload)
	if err != nil {
		fmt.Println("NewRequest ",err)
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do ",err)
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
