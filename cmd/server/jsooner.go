package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getJSONMetric(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")

	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		return
	}
	metra := Metrics{}
	err = json.Unmarshal([]byte(telo), &metra)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		return
	}
	//	metricType := metra.MType
	//	metricName := metra.ID
	//	metricValue := metra.Value
	//	metricDelta := metra.Delta

	switch metra.MType {
	case "counter":
		//err = memStor.getCounterValue(metricName, &val)
		memStor.mutter.RLock() // <---- MUTEX
		*metra.Delta = int64(memStor.count[metra.ID])
		memStor.mutter.RUnlock()
		resp, err := json.Marshal(metra)
		if err != nil {
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
			return
		}
		rwr.Write(resp)
	case "gauge":
		//		err = memStor.getGaugeValue(metricName, &val)
		memStor.mutter.RLock() // <---- MUTEX
		*metra.Value = float64(memStor.gau[metra.ID])
		memStor.mutter.RUnlock()
		resp, err := json.Marshal(metra)
		if err != nil {
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
			return
		}
		rwr.Write(resp)
	default:
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rwr, nil)
		return
	}
	if err != nil {
		rwr.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rwr, nil)
		return
	}
}

func treatJSONMetric(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")

	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		return
	}
	metra := Metrics{}
	json.Unmarshal([]byte(telo), &metra)
	metricType := metra.MType
	metricName := metra.ID
	metricValue := metra.Value
	metricDelta := metra.Delta

	//	fmt.Fprint(rwr, string(telo), metra)
	if metricValue == nil && metricDelta == nil {
		rwr.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rwr, `{"status":"StatusNotFound"}`)
		return
	}
	switch metricType {
	case "counter":
		if metricDelta == nil {
			rwr.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
			return
		}
		memStor.addCounter(metricName, counter(*metricDelta))
		memStor.mutter.RLock() // <---- MUTEX
		*metra.Delta = int64(memStor.count[metra.ID])
		memStor.mutter.RUnlock()
		resp, err := json.Marshal(metra)
		if err != nil {
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
			return
		}
		rwr.Write(resp)
	case "gauge":
		if metricValue == nil {
			rwr.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
			return
		}
		memStor.addGauge(metricName, gauge(*metricValue))
		memStor.mutter.RLock() // <---- MUTEX
		*metra.Value = float64(memStor.gau[metra.ID])
		memStor.mutter.RUnlock()
		resp, err := json.Marshal(metra)
		if err != nil {
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
			return
		}
		rwr.Write(resp)
	default:
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		return
	}
}
