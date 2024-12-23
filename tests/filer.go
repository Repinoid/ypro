package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
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
	memorial.mutter.RLock()
	ret, err := json.Marshal(MStorJSON{
		Gau:   memorial.gau,
		Count: memorial.count,
	})
	memorial.mutter.RUnlock()
	ret = append(ret, '\n')
	return ret, err
}
func (memorial *MemStorage) UnmarshalMS(data []byte) error {
	memorial.mutter.Lock()
	ms := MStorJSON{}
	err := json.Unmarshal(data, &ms)
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

func main() {

	g := map[string]gauge{"gs1": gauge(77.77), "gs2": gauge(88.88)}
	c := map[string]counter{"cs1": counter(44), "cs2": counter(88)}
	memStor = MemStorage{gau: g, count: c}

	err := memStor.LoadMS(fnam)
	fmt.Println(err)

	// ma, _ := memStor.MarshalMS()
	// se := MemStorage{}
	// se.UnmarshalMS(ma)

	fmt.Printf(" %v %v",  memStor.count, memStor.gau)

}
