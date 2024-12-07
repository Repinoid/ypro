package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	//"strings"
)

func postMetric() error {
	//host := "http://localhost:8080"

	//urla := "/s"

	//var jsonStr = []byte(`{"t":"B"}`)
	//req, err := http.NewRequest("POST", "http://localhost:8080/s", strings.NewReader("zalupan"))
	//req, err := http.NewRequest("POST", "http://localhost:8080/s", bytes.NewBuffer(jsonStr))
	//if err != nil {
	//	return fmt.Errorf("erra http.NewRequest %w ", err)
	//}
	//req.Header.Set("Content-Type", "application/json")

	url := "http://localhost:8080/s" //"http://" + host + "/update/" + metricType + "/" + metricName + "/" + metricValue
	//resp, err := http.Post(url, "application/json", nil)
	resp, err := http.Post(url, "application/json", strings.NewReader("zalupan"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	/*

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("client.Do http.NewRequest %w ", err)
		}
		defer res.Body.Close()
	*/

	/*	body1, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("erra io.ReadAll %w ", err)
		}*/
	var p []byte
	n, err := resp.Body.Read(p)
	if err != nil {
		return fmt.Errorf("erra io.ReadAll %w ", err)
	}
	fmt.Printf("body ...%v ... bodyread %v  bytes %d\n", "string(body1)", p, n)

	return nil
}
func main() {
	//	stat := postMetric("1", "2", "3")
	err := postMetric()
	if err != nil {
		log.Fatalf("fatalled postmetric -->\n%v\n<---\n", err)
	}
	fmt.Println("status ", err)
	//		panic(err)

}
