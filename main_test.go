package test

import (
	"fmt"
	"sync"
	"testing"
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

func TestXX(t *testing.T) {
	s := &Stats{
		//counters: make(map[string]int),
	}

	s.Add("aa", 1)
	s.Snapshot()

	for i := 1; i < 100; i++ {
	}
}
