package memorymodel

import (
	"sync"
)

var a1 string

var once sync.Once

func setup() {
	a1 = "hello"
}

func doprint() {
	once.Do(setup)
	print(a1)
}

func Run1() {
	go doprint()
	go doprint()
}
