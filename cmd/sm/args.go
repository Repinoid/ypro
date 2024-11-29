package main

import (
	"fmt"
	"os"
	"strconv"
	//"io"
	//"net/http"
	//"strings"
)

/*
Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
*/
var host = "localhost:8080"
var reportInterval = 10
var pollInterval = 2

func main() {
	if useClientArguments() != 0 {
		return
	}
}
func useClientArguments() int {
	args := os.Args[1:]

	for _, a := range args {
		if len(a) < 3 {
			fmt.Printf("unknown Argument -  %s\n", a)
			return 1
		}
		flagus := a[:3]
		tail := a[3:]
		switch flagus {
		case "-a=":
			host = a[3:]
		case "-r=":
			secs, err := strconv.Atoi(tail)
			if err != nil {
				fmt.Printf("Bad Argument for reportInterval with %s\n", tail)
				return 2
			}
			reportInterval = secs
		case "-p=":
			secs, err := strconv.Atoi(tail)
			if err != nil {
				fmt.Printf("Bad Argument for pollInterval with %s\n", tail)
				return 3
			}
			reportInterval = secs
		default:
			fmt.Printf("unknown Argument -  %s\n", a)
			return 4
		}
	}
	fmt.Println(host, pollInterval, reportInterval)
	return 0
}
