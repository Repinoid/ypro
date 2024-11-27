package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type gauge float64
type counter int64
type MemStorage struct {
	gau   map[string]gauge
	count map[string]counter
}

var memStor *MemStorage
var PollCount int

func getMetrix(memStor *MemStorage) error {
	var mS runtime.MemStats
	runtime.ReadMemStats(&mS)
	PollCount++
	memStor.gau = map[string]gauge{
		"Alloc":         gauge(mS.Alloc),
		"BuckHashSys":   gauge(mS.BuckHashSys),
		"Frees":         gauge(mS.Frees),
		"GCCPUFraction": gauge(mS.GCCPUFraction),
		"GCSys":         gauge(mS.GCSys),
		"HeapAlloc":     gauge(mS.HeapAlloc),
		"HeapIdle":      gauge(mS.HeapIdle),
		"HeapInuse":     gauge(mS.HeapInuse),
		"HeapObjects":   gauge(mS.HeapObjects),
		"HeapReleased":  gauge(mS.HeapReleased),
		"HeapSys":       gauge(mS.HeapSys),
		"LastGC":        gauge(mS.LastGC),
		"Lookups":       gauge(mS.Lookups),
		"MCacheInuse":   gauge(mS.MCacheInuse),
		"MCacheSys":     gauge(mS.MCacheSys),
		"MSpanInuse":    gauge(mS.MSpanInuse),
		"MSpanSys":      gauge(mS.MSpanSys),
		"Mallocs":       gauge(mS.Mallocs),
		"NextGC":        gauge(mS.NextGC),
		"NumForcedGC":   gauge(mS.NumForcedGC),
		"NumGC":         gauge(mS.NumGC),
		"OtherSys":      gauge(mS.OtherSys),
		"PauseTotalNs":  gauge(mS.PauseTotalNs),
		"StackInuse":    gauge(mS.StackInuse),
		"StackSys":      gauge(mS.StackSys),
		"Sys":           gauge(mS.Sys),
		"TotalAlloc":    gauge(mS.TotalAlloc),
		"RandomValue":   gauge(rand.Float64()),
	}
	memStor.count = map[string]counter{
		"PollCount": counter(PollCount),
	}
	return nil
}
func postMetric(metricType, metricName, metricValue string) int {
	host := "http://localhost:8080"
	url := host + "/update/" + metricType + "/" + metricName + "/" + metricValue
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return resp.StatusCode
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	memStor = new(MemStorage)

	for {
		for i := 0; i < 5; i++ {
			err := getMetrix(memStor)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(2 * time.Second)
		}
		for name, value := range memStor.gau {
			status := postMetric("gauge", name, strconv.FormatFloat(float64(value), 'f', 4, 64))
			if status != http.StatusOK {
				log.Println(status)
			}
			fmt.Printf("Metrix Name: %[1]s  value of %[2]f\n", name, value)
		}
		for name, value := range memStor.count {
			status := postMetric("counter", name, strconv.FormatInt(int64(value), 10))
			if status != http.StatusOK {
				log.Println(status)
			}
			fmt.Printf("Metrix Name: %[1]s  value of %[2]d\n", name, value)
		}
	}
//	return nil
}
