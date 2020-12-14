package main

import (
	"net/http"
	"runtime"
	"time"

	_ "net/http/pprof"
)

type Node struct {
	next *Node
	data [1024]byte
}

func HeapLeak() {
	go func() {
		for i := 0; i < 100000; i++ {
			go func() {
				nodes := make(map[int]*Node)
				j := 0
				for {
					j++
					nodes[j] = new(Node)
					time.Sleep(1 * time.Second)
				}
			}()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			println("[mem] ", memStats.HeapInuse)
		}
	}()

	http.ListenAndServe(":80", nil)
}
