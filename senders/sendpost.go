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
func unpackFromGzip(data2unpack io.Reader) (io.Reader, error) {
	gzipReader, err := gzip.NewReader(data2unpack)
	if err != nil {
		return nil, fmt.Errorf("gzip.NewReader %w ", err)
	}
	if err := gzipReader.Close(); err != nil {
		return nil, fmt.Errorf("zr.Close %w ", err)
	}
	return gzipReader, nil
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
	//requerest, err := http.NewRequest("POST", "http://"+host+"/pure", bytes.NewBuffer(jsonStrMarshalled))
	requerest, err := http.NewRequest("POST", "http://"+host+"/value/", bytes.NewBuffer(jsonStrPacked))
	if err != nil {
		return nil, fmt.Errorf("erra http.NewRequest %w ", err)
	}
	//requerest.Header.Set("Accept-Encoding", "gzip")
	requerest.Header.Set("Content-Type", "application/json")

	requerest.Header.Set("Content-Encoding", "gzip") // mark that data encoded

	client := &http.Client{}

	responsa, err := client.Do(requerest)
	if err != nil {
		return nil, fmt.Errorf("client.Do  %w ", err)
	}
	defer responsa.Body.Close()

	var reader io.Reader
	if responsa.Header.Get(`Content-Encoding`) == `gzip` { // if response is encoded
		reader, err = unpackFromGzip(responsa.Body)
		if err != nil {
			return nil, fmt.Errorf("unpackFromGzip %w ", err)
		}
	} else {
		reader = responsa.Body
	}
	telo, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll(reader) %w ", err)
	}

	return telo, nil
}
func main() {
	metr := `{"type":"gauge","id":"Alloc"}`
	//	err := postByPost(metr)
	outer, err := postJSONByNewRequest(metr)
	if err != nil {
		log.Fatalf("fatalled postmetric -->\n%v\t\t<---\n", err)
	}
	if err != nil && err != io.EOF {
		fmt.Println("Ошибка чтения декомпрессированных данных:", err)
		return
	}

	fmt.Printf("telo - \t%s", string(outer))
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
