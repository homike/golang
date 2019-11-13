package main

import (
	"fmt"
	"net/http"
	"time"

	_ "net/http/pprof"
)

func fn() int {
	return 1
}

func test() {
	ch := make(chan int)

	go func(ch1 chan int) {
		result := fn()
		fmt.Println("func start")
		time.Sleep(12 * time.Second)
		ch1 <- result
		fmt.Println("func end")
	}(ch)

	go func(ch1 chan int){
		result := fu()
		fmt.Println("func start")
		time.Sleep(12 * time.Second)
		ch1 <- result
		fmt.Println("func end")
	}(ch)

	select ch {
	case ch <- :
		fmt.Println("signl success")
	case time.After(10 * time.Second) :
		fmt.Println("time out")
	}

	fmt.Println("end")
}

type Pool struct {
	IsSaving bool
}

func (pool *Pool) addTeammateToSavePool(uid uint) {
	if !pool.IsSaving {
		pool.IsSaving = true
		go pool.keepSavingTeammateToRedis()
	}

	select {
	case pool.SavePool <- uid:
	default:
	}
}

func (pool *Pool) keepSavingTeammateToRedis() {
	for {
		select {
		case uid := <-pool.SavePool:
			fmt.Println(1)
			time.Sleep(time.Second / 10)
		case <-time.After(time.Hour * 24):
			pool.IsSaving = false
			break
		}
	}
}

func main() {

	go func() {
		fmt.Println(http.ListenAndServe(":6060", nil))
	}()

	test()

	for {
	}
}
