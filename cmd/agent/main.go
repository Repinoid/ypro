package main

import (
	"fmt"
	"math/rand/v2"
	"runtime"
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

func main() {
	memStor = new(MemStorage)

	getMetrix(memStor)
	for name, value := range memStor.gau {
		fmt.Printf("Metrix Name: %[1]s  value of %[2]f\n", name, value)
	}
	for name, value := range memStor.count {
		fmt.Printf("Metrix Name: %[1]s  value of %[2]d\n", name, value)
	}
}
