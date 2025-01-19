package main

import (
	"app/internal/dbaser"
	"app/internal/memo"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

func badPost(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "text/html")
	rwr.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(rwr, `{"status":"StatusNotFound"}`)
}

func getAllMetrix(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "text/html")
	if req.URL.Path != "/" { // if GET with wrong arguments structure
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		return
	}
	if MetricBaseStruct.IsBase {
		mGauge := map[string]float64{}
		err := dbaser.TableGetAllsWrapper(dbaser.TableGetAllGauges)(&MetricBaseStruct, &mGauge)
		if err != nil {
			log.Printf("bad allgauges\n %v\n", err)
		}
		mCounter := map[string]int64{}
		err = dbaser.TableGetAllsWrapper(dbaser.TableGetAllCounters)(&MetricBaseStruct, &mCounter)
		if err != nil {
			log.Printf("bad allcounters\n %v\n", err)
		}
		for nam, val := range mGauge {
			flo := strconv.FormatFloat(float64(val), 'f', -1, 64) // -1 - to remove zeroes tail
			fmt.Fprintf(rwr, "Gauge Metric name   %20s\t\tvalue\t%s\n", nam, flo)
		}
		for nam, val := range mCounter {
			fmt.Fprintf(rwr, "Counter Metric name %20s\t\tvalue\t%d\n", nam, val)
		}
		rwr.WriteHeader(http.StatusOK)
		return
	}
	var mutter sync.RWMutex
	mutter.RLock() // <---- MUTEX
	defer mutter.RUnlock()
	for nam, val := range memStor.Gaugemetr {
		flo := strconv.FormatFloat(float64(val), 'f', -1, 64) // -1 - to remove zeroes tail
		fmt.Fprintf(rwr, "Gauge Metric name   %20s\t\tvalue\t%s\n", nam, flo)
	}
	for nam, val := range memStor.Countmetr {
		fmt.Fprintf(rwr, "Counter Metric name %20s\t\tvalue\t%d\n", nam, val)
	}
	rwr.WriteHeader(http.StatusOK)
}
func getMetric(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(req)
	metricType := vars["metricType"]
	metricName := vars["metricName"]
	switch metricType {
	case "counter":
		var cunt counter
		if memo.GetCounterValue(&memStor, MetricBaseStruct, metricName, &cunt) != nil {
			//	if memStor.GetCounterValue(metricName, &cunt) != nil {
			rwr.WriteHeader(http.StatusNotFound)
			fmt.Fprint(rwr, nil)
			return
		}
		fmt.Fprint(rwr, cunt)
	case "gauge":
		var gaaga gauge
		if memo.GetGaugeValue(&memStor, MetricBaseStruct, metricName, &gaaga) != nil {
			//	if memStor.GetGaugeValue(metricName, &gaaga) != nil {
			rwr.WriteHeader(http.StatusNotFound)
			fmt.Fprint(rwr, nil)
			return
		}
		fmt.Fprint(rwr, gaaga)
	default:
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rwr, nil)
		return
	}
	rwr.WriteHeader(http.StatusOK)
}

func treatMetric(rwr http.ResponseWriter, req *http.Request) {

	//log.Printf("%v\nisBase - %v\ncheck - %v\n\n\n", MetricBase.MetricBase, isBase, check)

	rwr.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(req)
	metricType := vars["metricType"]
	metricName := vars["metricName"]
	metricValue := vars["metricValue"]
	if metricValue == "" {
		rwr.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rwr, `{"status":"StatusNotFound"}`)
		return
	}
	switch metricType {
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			rwr.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
			return
		}
		//	memStor.AddCounter(metricName, counter(value))
		memo.AddCounter(&memStor, MetricBaseStruct, metricName, counter(value))
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			rwr.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
			return
		}
		//	memStor.AddGauge(metricName, gauge(value))
		memo.AddGauge(&memStor, MetricBaseStruct, metricName, gauge(value))
	default:
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		return
	}
	rwr.WriteHeader(http.StatusOK)
	fmt.Fprintf(rwr, `{"status":"StatusOK"}`)

	if storeInterval == 0 {
		_ = memStor.SaveMS(fileStorePath)
	}
}
