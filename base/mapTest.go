package basetest

import (
	"fmt"
	"time"
)

func mapWriteTest() {
	c := make(map[string]int)
	go func() { //开一个协程写map
		for j := 0; j < 1000000; j++ {
			c[fmt.Sprintf("%d", j)] = j
		}
	}()
	go func() { //开一个协程读map
		for j := 0; j < 1000000; j++ {
			fmt.Println(c[fmt.Sprintf("%d", j)])
		}
	}()

	time.Sleep(time.Second * 20)
}

func mapWriteReadTest() {
	c := make(map[string]int)

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 1000000; j++ {
				c[fmt.Sprintf("%d", j)] = j
			}
		}()
	}
	time.Sleep(time.Second * 20) //让执行main函数的主协成等待20s,不然不会执行上面的并发操作
}

func mapReadTest() {
	c := make(map[string]int)
	for j := 0; j < 1000000; j++ {
		c[fmt.Sprintf("%d", j)] = j
	}

	for i := 0; i < 10; i++ {
		go func() { //开一个协程读map
			for j := 0; j < 1000000; j++ {
				fmt.Println(c[fmt.Sprintf("%d", j)])
			}
		}()
	}

	time.Sleep(time.Second * 20)
}

func RunMapTest() {
	// var counter = struct {
	// 	sync.RWMutex
	// 	m map[string]int
	// }{m: make(map[string]int)}
}
