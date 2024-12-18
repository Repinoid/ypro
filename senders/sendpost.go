package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	//"strings"
)

var host = "localhost:8080"

func postJSONByNewRequest(jsonStr string) (*gzip.Reader, error, int) {

	jsonStrMarshalled, err := json.Marshal(jsonStr)
	if err != nil {
		return nil, fmt.Errorf("marshal err %w ", err), 0
	}
	requerest, err := http.NewRequest("POST", "http://"+host+"/params", bytes.NewBuffer(jsonStrMarshalled))
	if err != nil {
		return nil, fmt.Errorf("erra http.NewRequest %w ", err), 0
	}
	fmt.Println("Header ", requerest.Header)
	//     	requerest.Header.Set("Content-Type", "application/json")
	requerest.Header.Set("Accept-Encoding", "gzip")
	requerest.Header.Set("Content-Encoding", "gzip;zalupan")

	client := &http.Client{}

	responsa, err := client.Do(requerest)
	if err != nil {
		return nil, fmt.Errorf("client.Do  %w ", err), 0
	}
	defer responsa.Body.Close()

	dlina, errD := strconv.Atoi(responsa.Header.Get("Content-Length"))
	if errD != nil {
		return nil, fmt.Errorf("strconv.Atoi(responsa.Header.Get(\"Content-Length\")) %w ", errD), 0
	}

	gzipReader, err := gzip.NewReader(responsa.Body)
	if err != nil {
		return nil, fmt.Errorf("gzip.NewReader %w ", err), 0
	}

	if err := gzipReader.Close(); err != nil {
		return nil, fmt.Errorf("zr.Close %w ", err), 0
	}

	return gzipReader, nil, dlina
}
func main() {
	metr := `{"a":7,"b":"1234a"}`
	//	err := postByPost(metr)
	outer, err, dlina := postJSONByNewRequest(metr)
	if err != nil {
		log.Fatalf("fatalled postmetric -->\n%v\n<---\n", err)
	}
	decompressedData := make([]byte, dlina)
	_, err = outer.Read(decompressedData)
	if err != nil && err != io.EOF {
		fmt.Println("Ошибка чтения декомпрессированных данных:", err)
		return
	}

	fmt.Printf("%v", string(decompressedData))
	// if _, err := io.Copy(os.Stdout, outer); err != nil {
	// 	log.Fatal(err)
	// }
}

func postByPost(metr string) error {

	march, _ := json.Marshal(metr)
	respon, err := http.Post("http://"+host+"/params", "application/json", bytes.NewBuffer(march))
	if err != nil {
		return err
	}
	defer respon.Body.Close()

	//	fmt.Printf("bashka %v\n", respon.Header)

	respon.Header.Set("Accept-Encoding", "gzip")

	zr, err := gzip.NewReader(respon.Body)
	if err != nil {
		log.Fatal(err)
	}

	if written, err := io.Copy(os.Stdout, zr); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("\nWritten %d\nName: %s\nComment: %s\nModTime: %s\n\n", written, zr.Name, zr.Comment, zr.ModTime.UTC())
	}
	fmt.Printf("Header %v", respon.Header)
	// body1, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return fmt.Errorf("erra io.ReadAll %w ", err)
	// }
	// fmt.Printf("body ...%v ...\n", string(body1))
	return nil
}
