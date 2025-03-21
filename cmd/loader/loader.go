package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"gorono/internal/memos"
	"gorono/internal/middlas"
	"gorono/internal/models"
	"gorono/internal/privacy"

	"github.com/go-resty/resty/v2"
)

var host = "localhost:8080"
var reportInterval = 100
var pollInterval = 100

var key = ""

// var key = "keyForCrypta"
var rateLimit = 4
var cunt int64

func main() {
	if err := initLoader(); err != nil {
		log.Fatal("INTERVALS error ", err)
		return
	}

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	marshalledGau, _ := json.Marshal(m1)
	postMetras(marshalledGau, "/update/")
	postMetras(marshalledGau, "/value/")
	marshalledCunt, _ := json.Marshal(m2)
	postMetras(marshalledCunt, "/update/")
	postMetras(marshalledCunt, "/value/")

	const chanCap = 4

	metroBarn := make(chan []models.Metrics, chanCap)
	go metrixIN(metroBarn)

	fenix := make(chan struct{})
	for w := 1; w <= rateLimit; w++ {
		go bolda(metroBarn, fenix)
	}
	for {
		fenix <- struct{}{}        // блокируем канал пока балда не прочитает из него при своём завершении по ошибке
		go bolda(metroBarn, fenix) // нанимаем нового
	}

}

// получает банчи метрик и складывает в barn
func metrixIN(metroBarn chan<- []models.Metrics) {
	memStorage := []models.Metrics{}
	tickerPoll := time.NewTicker(time.Duration(pollInterval) * time.Millisecond)
	tickerReport := time.NewTicker(time.Duration(reportInterval) * time.Millisecond)
	for {
		select {
		case <-tickerPoll.C:
			memStorage = *memos.GetMetrixFromOS()
			addMetrix := *memos.GetMoreMetrix()
			memStorage = append(memStorage, addMetrix...)
			atomic.AddInt64(&cunt, 1) //			cunt++

			for ind, metr := range memStorage {
				if metr.ID == "PollCount" && metr.MType == "counter" { // search for PollCount metric
					cu := atomic.LoadInt64(&cunt)
					memStorage[ind].Delta = &cu // memStorage[ind].Delta = cunt
					break
				}
			}
		case <-tickerReport.C:
			metroBarn <- memStorage
		}
	}
}

// работник отсылает банчи метрик на сервер, феникс - канал для подачи сигнала о завершении по ошибке
func bolda(metroBarn <-chan []models.Metrics, fenix <-chan struct{}) {
	for {
		bunch := <-metroBarn
		marshalledBunch, err := json.Marshal(bunch)
		if err != nil {
			<-fenix // в случае ошибки читаем из феникса, разблокируя канал и выходим
			return
		}
		var haHex string
		if key != "" {
			keyB := md5.Sum([]byte(key))

			coded, err := privacy.EncryptB2B(marshalledBunch, keyB[:])
			if err != nil {
				<-fenix
				return
			}
			ha := privacy.MakeHash(nil, coded, keyB[:])
			haHex = hex.EncodeToString(ha)
			marshalledBunch = coded
		}
		compressedBunch, err := middlas.Pack2gzip(marshalledBunch)
		if err != nil {
			<-fenix
			return
		}

		httpc := resty.New() //
		httpc.SetBaseURL("http://" + host)

		req := httpc.R().
			SetHeader("Content-Encoding", "gzip"). // сжаtо
			SetBody(compressedBunch).
			SetHeader("Accept-Encoding", "gzip")

		if key != "" {
			req.Header.Add("HashSHA256", haHex) // Хеш в заголовок, значит - зашифровано
		}

		resp, _ := req.
			SetDoNotParseResponse(false).
			Post("/updates/") // slash on the tile
		if resp.StatusCode() == http.StatusOK { // при успешной отправке метрик обнуляем cчётчик
			atomic.StoreInt64(&cunt, 0) //	cunt = 0
		} else {
			log.Printf("AGENT responce from server %+v\n", resp.StatusCode())
		}
	}
}

var m1 = models.Metrics{MType: "gauge", ID: "gaaga", Value: middlas.Ptr(777.77)}
var m2 = models.Metrics{MType: "counter", ID: "cunt", Delta: middlas.Ptr[int64](777)}

func postMetras(body []byte, urla string) {
	tickerReport := time.NewTicker(time.Duration(reportInterval) * time.Millisecond)
	go func() {
		for range tickerReport.C {
			httpc := resty.New() //
			httpc.SetBaseURL("http://" + host)
			req := httpc.R().
				SetBody(body)
			resp, err := req.
				SetDoNotParseResponse(false).
				Post(urla) // slash on the tile
			if resp.StatusCode() != http.StatusOK || err != nil {
				log.Printf("StatusCode %+v err %+v\n", resp.StatusCode(), err)
			}
		}
	}()
}
