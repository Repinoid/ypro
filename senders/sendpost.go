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
	"time"
	//"strings"
)

var host = "localhost:8080"

func pack2gzip(data2pack []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.ModTime = time.Now()
	_, err := zw.Write(data2pack)
	if err != nil {
		return nil, fmt.Errorf("gzip.NewWriter.Write %w ", err)
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("gzip.NewWriter.Close %w ", err)
	}
	return buf.Bytes(), nil
}
func unpackFromGzip(data2unpack io.Reader) ([]byte, error) {
	gzipReader, err := gzip.NewReader(data2unpack)
	if err != nil {
		return nil, fmt.Errorf("gzip.NewReader %w ", err)
	}
	if err := gzipReader.Close(); err != nil {
		return nil, fmt.Errorf("zr.Close %w ", err)
	}
	decompressedData := make([]byte, 1000)
	_, err = gzipReader.Read(decompressedData)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("gzipReader.Read %w ", err)
	}
	return decompressedData, nil
}

func postJSONByNewRequest(jsonStr string) ([]byte, error) {

	jsonStrMarshalled, err := json.Marshal(jsonStr)
	if err != nil {
		return nil, fmt.Errorf("marshal err %w ", err)
	}
	jsonStrPacked, err := pack2gzip(jsonStrMarshalled)
	if err != nil {
		return nil, fmt.Errorf("pack2gzip %w ", err)
	}
	requerest, err := http.NewRequest("POST", "http://"+host+"/params", bytes.NewBuffer(jsonStrPacked))
	if err != nil {
		return nil, fmt.Errorf("erra http.NewRequest %w ", err)
	}
	//fmt.Println("Header ", requerest.Header)
	//     	requerest.Header.Set("Content-Type", "application/json")
	requerest.Header.Set("Accept-Encoding", "gzip")
	requerest.Header.Set("Content-Encoding", "gzip")

	client := &http.Client{}

	responsa, err := client.Do(requerest)
	if err != nil {
		return nil, fmt.Errorf("client.Do  %w ", err)
	}
	defer responsa.Body.Close()

	unpacked, err := unpackFromGzip(responsa.Body)
	if err != nil {
		return nil, fmt.Errorf("unpackFromGzip %w ", err)
	}
	return unpacked, nil
}
func main() {
	metr := `{"a":7,"b":"1234a"}`
	//	err := postByPost(metr)
	outer, err := postJSONByNewRequest(metr)
	if err != nil {
		log.Fatalf("fatalled postmetric -->\n%v\n<---\n", err)
	}
	if err != nil && err != io.EOF {
		fmt.Println("Ошибка чтения декомпрессированных данных:", err)
		return
	}

	fmt.Printf("%v", string(outer))
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
