package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"sync"
	"time"
)

type gauge float64
type counter int64
type MemStorage struct {
	gau    map[string]gauge
	count  map[string]counter
	mutter sync.RWMutex
}
type MStorJSON struct {
	Gau   map[string]gauge
	Count map[string]counter // `json:"count"`
}

type Metrica struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

var memStor MemStorage
var fnam = "ms.txt"

func (memorial *MemStorage) MarshalMS() ([]byte, error) {
	buf := new(bytes.Buffer)
	memorial.mutter.RLock()
	err := json.NewEncoder(buf).Encode(MStorJSON{
		Gau:   memorial.gau,
		Count: memorial.count,
	})
	memorial.mutter.RUnlock()
	return append(buf.Bytes(), '\n'), err
}
func (memorial *MemStorage) UnmarshalMS(data []byte) error {
	ms := MStorJSON{}
	memorial.mutter.Lock()
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&ms)
	memorial.gau = ms.Gau
	memorial.count = ms.Count
	memorial.mutter.Unlock()
	return err
}

func (memorial *MemStorage) SaveMS(fnam string) error {
	phil, err := os.OpenFile(fnam, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("file %s Open error %v", fnam, err)
	}
	march, err := memorial.MarshalMS()
	if err != nil {
		return fmt.Errorf(" Memstorage Marshal error %v", err)
	}
	_, err = phil.Write(march)
	if err != nil {
		return fmt.Errorf("file %s Write error %v", fnam, err)
	}
	return nil
}

func (memorial *MemStorage) LoadMS(fnam string) error {
	phil, err := os.OpenFile(fnam, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("file %s Open error %v", fnam, err)
	}
	reader := bufio.NewReader(phil)
	data, err := reader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("file %s Read error %v", fnam, err)
	}
	err = memorial.UnmarshalMS(data)
	if err != nil {
		return fmt.Errorf(" Memstorage UnMarshal error %v", err)
	}
	return nil
}

func saver(memStor *MemStorage, fnam string) error {

	for {
		time.Sleep(time.Duration(1) * time.Second)
		err := memStor.SaveMS(fnam)
		if err != nil {
			return fmt.Errorf("save err %v", err)
		}
	}
	//return nil

}

func maina() {

	g := map[string]gauge{"gs1": gauge(99.77), "gs2": gauge(88.88)}
	c := map[string]counter{"cs1": counter(rand.IntN(100)), "cs2": counter(rand.IntN(100))}
	memStor = MemStorage{gau: g, count: c}
	//memStor = MemStorage{}

	err := memStor.LoadMS(fnam)
	fmt.Println(err)

	go saver(&memStor, "out.out")

	var inp string
	fmt.Scan(&inp)

	// ma, _ := memStor.MarshalMS()
	// se := MemStorage{}
	// se.UnmarshalMS(ma)

	fmt.Printf(" %v %v", memStor.count, memStor.gau)

}
