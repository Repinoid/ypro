package main

import (
	"io"
	"log"
	"strings"
)

func main4() {
	r := strings.NewReader("some io.Reader stream to be read\t")
	buf := []byte("                                                              ")
	//	tee := io.Reader(&buf)

	n, err := r.Read(buf)

//	r.

	r1, err1 := io.ReadAll(r)
//	r2, err2 := io.ReadAll(r)

	log.Printf("%s %v\n", r1, err1)
	log.Printf("%d %s %v\n", n, buf, err)
}
