package memos

import (
	"math/rand/v2"
	"runtime"

	"gorono/internal/models"
)

type gauge = models.Gauge
type counter = models.Counter

func GetMetrixFromOS() *[]models.Metrics {
	var mS runtime.MemStats
	runtime.ReadMemStats(&mS)
	gaugeMap := map[string]gauge{
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
		"RandomValue":   gauge(rand.Float64()), // self-defined
	}
	counterMap := map[string]counter{
		"PollCount": counter(0), // self-defined
	}
	metrArray := make([]models.Metrics, 0, len(gaugeMap)+len(counterMap))

	for metrName, metrValue := range counterMap {
		mval := int64(metrValue) + int64(rand.IntN(10))
		metr := models.Metrics{ID: metrName, MType: "counter", Delta: &mval}
		metrArray = append(metrArray, metr)
	}
	for metrName, metrValue := range gaugeMap {
		mval := float64(metrValue) + float64(rand.IntN(10))
		metr := models.Metrics{ID: metrName, MType: "gauge", Value: &mval}
		metrArray = append(metrArray, metr)
	}
	return &metrArray
}
