package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	//"strings"
)

var host = "localhost:8080"

func postByNewRequest() error {

	jsonStr := `{"a":7}`
	jsonStrM, _ := json.Marshal(jsonStr)
	//req, err := http.NewRequest("POST", "http://localhost:8080/s", nil)
	//	req, err := http.NewRequest("POST", "http://localhost:8080/s", strings.NewReader("zalupan"))
	req, err := http.NewRequest("POST", "http://localhost:8080/s", bytes.NewBuffer(jsonStrM))
	if err != nil {
		return fmt.Errorf("erra http.NewRequest %w ", err)
	}
	//req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Type", "application/json")
	//	req.Header.Set("Content-Encoding", "gzip")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do http.NewRequest %w ", err)
	}
	defer res.Body.Close()

	body1, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("erra io.ReadAll %w ", err)
	}
	fmt.Printf("body ...%v ...\n", string(body1))

	return nil
}
func main() {
	//	stat := postMetric("1", "2", "3")
	err := postByPost()
	if err != nil {
		log.Fatalf("fatalled postmetric -->\n%v\n<---\n", err)
	}
	fmt.Println("status ", err)
	//		panic(err)

}

func postByPost() error {
	metr := `{"a":7,"b":"1234a"}`
	march, _ := json.Marshal(metr)
	resp, err := http.Post("http://"+host+"/params", "application/json", bytes.NewBuffer(march))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body1, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erra io.ReadAll %w ", err)
	}
	fmt.Printf("body ...%v ...\n", string(body1))
	return nil
}
