package main

import (
	"fmt"
	"strconv"
	"sync"
)

<<<<<<< HEAD
func newMemStorage() *MemStorage {
	return &MemStorage{
		mutter: &sync.RWMutex{},
		gau : make(map[string]gauge),
		count : make(map[string]counter),
	}
//	return &ret
=======
func newMemStorage() (MemStorage, error) {
	memStor := new(MemStorage)
	memStor.gau = make(map[string]gauge)
	memStor.count = make(map[string]counter)
	return *memStor, nil
>>>>>>> 3a4833708c655630ca9e1b2b221add7df5364d30
}
func (ms *MemStorage) addGauge(name string, value gauge) error {
	ms.mutter.Lock()
	defer ms.mutter.Unlock()
	ms.gau[name] = value
	return nil
}
func (ms *MemStorage) addCounter(name string, value counter) error {
	ms.mutter.Lock()
	defer ms.mutter.Unlock()
	if _, ok := ms.count[name]; ok {
		ms.count[name] += value
		return nil
	}
	ms.count[name] = value
	return nil
}

/*
	func (ms *MemStorage) resetPollCount() error {
		ms.mutter.Lock()
		defer ms.mutter.Unlock()
		if _, ok := ms.count["PollCount"]; ok {
			ms.count["PollCount"] = 0
		} else {
			return fmt.Errorf("no PollCount")
		}
		return nil
	}
*/
func (ms *MemStorage) getCounterValue(name string, value *string) error {
	if _, ok := ms.count[name]; ok {
		*value = strconv.FormatInt(int64(ms.count[name]), 10)
		return nil
	}
	return fmt.Errorf("no %s key", name)
}
func (ms *MemStorage) getGaugeValue(name string, value *string) error {
	if _, ok := ms.gau[name]; ok {
		*value = strconv.FormatFloat(float64(ms.gau[name]), 'f', -1, 64)
		return nil
	}
	return fmt.Errorf("no %s key", name)
}
