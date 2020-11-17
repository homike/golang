// map_test.go 测试golang中的map的并发读写能力
package base

import (
	"fmt"
	"testing"
	"time"
)

// 并发读写会panic
func _Test_MapWriteRead(t *testing.T) {
	c := make(map[string]int)
	go func() { //开一个协程写map
		for j := 0; j < 1000000; j++ {
			c[fmt.Sprintf("%d", j)] = j
		}
	}()
	go func() { //开一个协程读map
		for j := 0; j < 1000000; j++ {
			_ = c[fmt.Sprintf("%d", j)]
		}
	}()

	time.Sleep(time.Second * 20)
}

// 并发读不会panic
func Test_MapRead(t *testing.T) {
	c := make(map[string]int)
	for j := 0; j < 1000000; j++ {
		c[fmt.Sprintf("%d", j)] = j
	}

	for i := 0; i < 10; i++ {
		go func() { //开一个协程读map
			for j := 0; j < 1000000; j++ {
				_ = c[fmt.Sprintf("%d", j)]
			}
		}()
	}

	time.Sleep(time.Second * 10)
}

// 并发写会panic
func _Test_MapWrite(t *testing.T) {
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

func _Test_Map(t *testing.T) {
	// var counter = struct {
	// 	sync.RWMutex
	// 	m map[string]int
	// }{m: make(map[string]int)}
}
