package stack

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

func PrintStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func PrintGoroutineMemConsume() {
	var c chan int
	//var wg sync.WaitGroup
	const goroutineNum = 1 // 1 * 10^4

	memConsumed := func() uint64 {
		runtime.GC() //GC，排除对象影响
		var memStat runtime.MemStats
		runtime.ReadMemStats(&memStat)
		return memStat.Sys
	}

	noop := func(lenChan chan int) {
		//wg.Done()
		for {
			select {
			case <-c:
				return
			case l := <-lenChan:
				val := make([]int, l)
				_ = val
			}
		}
	}

	//wg.Add(goroutineNum)
	before := memConsumed() //获取创建goroutine前内存

	chStack := make([]chan int, goroutineNum)
	fmt.Println("aaa")
	for i := 0; i < goroutineNum; i++ {
		go noop(chStack[i])
	}
	fmt.Println("bbb")

	{
		fmt.Println("b1")
		for i := 0; i < goroutineNum; i++ {
			select {
			case chStack[i] <- 10:
				fmt.Println("b1---->1")
			default:
				fmt.Println("b1---->2")
			}
		}
		fmt.Println("b2")

		after := memConsumed()
		fmt.Printf("%.3f KB\n", float64(after-before)/goroutineNum/1000)
	}

	fmt.Println("ccc")
	{
		for i := 0; i < goroutineNum; i++ {
			chStack[i] <- math.MaxInt32
		}

		after := memConsumed() //获取创建goroutine后内存
		fmt.Printf("%.3f KB\n", float64(after-before)/goroutineNum/1000)
	}

	for {
		for i := 0; i < goroutineNum; i++ {
			chStack[i] <- 10
		}
		after := memConsumed()
		fmt.Printf("%.3f KB\n", float64(after-before)/goroutineNum/1000)

		time.Sleep(1 * time.Second)
	}

	//wg.Wait()
}

func getGoroutineMemConsume() {
	var c chan int
	var wg sync.WaitGroup
	const goroutineNum = 1e4 // 1 * 10^4

	memConsumed := func() uint64 {
		runtime.GC() //GC，排除对象影响
		var memStat runtime.MemStats
		runtime.ReadMemStats(&memStat)
		return memStat.Sys
	}

	noop := func() {
		wg.Done()
		<-c //防止goroutine退出，内存被释放
	}

	wg.Add(goroutineNum)
	before := memConsumed() //获取创建goroutine前内存
	for i := 0; i < goroutineNum; i++ {
		go noop()
	}
	wg.Wait()
	after := memConsumed() //获取创建goroutine后内存

	fmt.Printf("%.3f KB\n", float64(after-before)/goroutineNum/1000)
}
