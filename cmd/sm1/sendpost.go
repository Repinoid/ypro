package main

import (
	"fmt"
	"io"
	"net/http"
	//"strings"
)

func postMetric(metricType, metricName, metricValue string) error {
	host := "http://localhost:8080"

	client := &http.Client{}
	payloadStr := "/update/" + metricType + "/" + metricName + "/" + metricValue
	//	payloadStr = "/u"
	//	payload := strings.NewReader(payloadStr)
	//	req, err := http.NewRequest(http.MethodPost, host, payload)
	req, err := http.NewRequest(http.MethodPost, host+payloadStr, nil)
	fmt.Printf("payload string ---  \"%[1]v\" type %[1]T\n", payloadStr)
	fmt.Printf("req ---  \"%[1]v\" type %[1]T\n", req)
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

	body1, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("io.ReadAll ", err)
		return err
	}
	fmt.Printf("body ...%v\n", string(body1))

	return nil
}
func main() {
	//	stat := postMetric("1", "2", "3")
	stat := postMetric("gauge", "Alloc", "55.66")
	fmt.Println("status ", stat)
	//		panic(err)

}
