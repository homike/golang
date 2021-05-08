package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"
)

type Stats struct {
	//mutex sync.Mutex

	//counters map[string]int

	counters sync.Map
}

func (s *Stats) Snapshot() {
	//s.mutex.Lock()
	//defer s.mutex.Unlock()

	//return s.counters

	s.counters.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}

func (s *Stats) Add(name string, num int) {
	//s.mutex.Lock()
	//defer s.mutex.Unlock()

	//s.counters[name] = num
	s.counters.Store(name, num)
}

func _TestXX(t *testing.T) {
	s := &Stats{
		//counters: make(map[string]int),
	}

	s.Add("aa", 1)
	s.Snapshot()

	for i := 1; i < 100; i++ {
	}
}

func _TestXX2(t *testing.T) {
	var x int
	inc := func() int {
		x++
		return x
	}
	fmt.Println(func() (a, b int) {
		return inc(), inc()
	}())
}

func _TestXX3(t *testing.T) {
	t1 := struct {
		time.Time
		N int
	}{
		time.Date(2020, 12, 20, 0, 0, 0, 0, time.UTC),
		5000,
	}

	m, _ := json.Marshal(t1)
	fmt.Printf("%s", m)
}

func _TestXX4(t *testing.T) {
	chClose := make(chan int)
	go func() {
		select {
		case <-chClose:
			return
		}
	}()

	close(chClose)
	chClose <- 1
}

func TestXX5(t *testing.T) {
	zoneID, svrID := 0, 123
	key1 := (int64(zoneID) << 32) | int64(svrID)
	key2 := int64(svrID)

	if key1 == key2 {
		fmt.Println("true", key1, " = ", key2)
	} else {
		fmt.Println("false", key1, " != ", key2)
	}
}
