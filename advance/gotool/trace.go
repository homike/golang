package main

import (
	"math/rand"
	"os"
	"runtime/trace"
	"sync"
)

//func TestTrace(t *testing.T) {
func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	// Your program here

	var total int
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 1000; j++ {
				total += readNumber()
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

//go:noinline
func readNumber() int {
	return rand.Intn(10)
}
