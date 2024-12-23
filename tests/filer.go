package main

import (
	"encoding/json"
	"fmt"
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
	Count map[string]counter// `json:"count"`
}

type Metrica struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

var memStor MemStorage

func (m *MemStorage) MarshalMS() ([]byte, error) {
	ret, err := json.Marshal(MStorJSON{
		Gau:   m.gau,
		Count: m.count,
	})
	return ret, err
}
func (m *MemStorage) UnmarshalMS(data []byte) error {
	ms := MStorJSON{}
	err := json.Unmarshal(data, &ms)
	m.gau = ms.Gau
	m.count = ms.Count
	return err
}

func main() {

	g := map[string]gauge{"gs1": gauge(77.77), "gs2": gauge(88.88)}
	c := map[string]counter{"cs1": counter(77), "cs2": counter(88)}
	memStor = MemStorage{gau: g, count: c}

	ma, _ := memStor.MarshalMS()
	se := MemStorage{}
	se.UnmarshalMS(ma)

	fmt.Printf("%s\n %v %v", ma, se.count, se.gau)

}
