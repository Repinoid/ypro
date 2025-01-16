package main

import (
	"app/internal/dbaser"
	"app/internal/memo"
	"app/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
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
	switch metra.MType {
	case "counter":
		var cunt counter
		if memo.GetCounterValue(&memStor, MetricBaseStruct, metra.ID, &cunt) != nil {
			rwr.WriteHeader(http.StatusNotFound)
			fmt.Fprint(rwr, nil)
			return
		}
		delt := int64(cunt)
		metra.Delta = &delt
		json.NewEncoder(rwr).Encode(metra)
	case "gauge":
		var gaaga gauge
		if memo.GetGaugeValue(&memStor, MetricBaseStruct, metra.ID, &gaaga) != nil {
			rwr.WriteHeader(http.StatusNotFound)
			fmt.Fprint(rwr, nil)
			return
		}
		flo := float64(gaaga)
		metra.Value = &flo
		json.NewEncoder(rwr).Encode(metra)
	default:
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rwr, nil)
		return
	}
	rwr.WriteHeader(http.StatusOK)

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
	err = json.Unmarshal([]byte(telo), &metra)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		return
	}

	metricType := metra.MType
	metricName := metra.ID
	metricValue := metra.Value
	metricDelta := metra.Delta

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
		rwr.WriteHeader(http.StatusOK)
		memo.AddCounter(&memStor, MetricBaseStruct, metricName, counter(*metricDelta))
		// get new value from memstorage
		var cunt counter
		if memo.GetCounterValue(&memStor, MetricBaseStruct, metra.ID, &cunt) != nil {
			rwr.WriteHeader(http.StatusNotFound)
			fmt.Fprint(rwr, nil)
			return
		}
		*metra.Delta = int64(cunt)
		json.NewEncoder(rwr).Encode(metra)
	case "gauge":
		if metricValue == nil {
			rwr.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
			return
		}
		rwr.WriteHeader(http.StatusOK)

		memo.AddGauge(&memStor, MetricBaseStruct, metricName, gauge(*metricValue))
		// get new value from memstorage
		var gaaga gauge
		if memo.GetGaugeValue(&memStor, MetricBaseStruct, metra.ID, &gaaga) != nil {
			rwr.WriteHeader(http.StatusNotFound)
			fmt.Fprint(rwr, nil)
			return
		}
		*metra.Value = float64(gaaga)
		json.NewEncoder(rwr).Encode(metra)
	default:
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		return
	}

	if storeInterval == 0 {
		_ = memStor.SaveMS(fileStorePath)
	}
}
func buncheras(rwr http.ResponseWriter, req *http.Request) {
	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest)
		return
	}
	buf := bytes.NewBuffer(telo)
	memor := []models.Metrics{}
	err = json.NewDecoder(buf).Decode(&memor)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest)
		return
	}
	if MetricBaseStruct.IsBase {
		err = dbaser.TableBuncherWrapper(dbaser.TableBuncher)(&MetricBaseStruct, memor)
		if err != nil {
			log.Printf("%-v", err)
		}
	}
	var mutter sync.RWMutex
	mutter.Lock()
	for _, m := range memor {
		switch m.MType {
		case "gauge":
			memStor.Gaugemetr[m.ID] = gauge(*m.Value)
		case "counter":
			if _, ok := memStor.Countmetr[m.ID]; ok {
				memStor.Countmetr[m.ID] += counter(*m.Delta)
				continue
			}
			memStor.Countmetr[m.ID] = counter(*m.Delta)

		default:
			log.Printf("wrong metric type %s\n", m.MType)
		}
	}
	mutter.Unlock()

	json.NewEncoder(rwr).Encode(&memor)

}
