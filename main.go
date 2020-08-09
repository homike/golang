package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"runtime"
	"time"
)

const kItemTypeOffset = 100000

func GetItemKey(itemType, itemID int32) int32 {
	return itemType*kItemTypeOffset + itemID
}

func ParseItemKey(itemKey int32) (int32, int32) {
	return itemKey / kItemTypeOffset, itemKey % kItemTypeOffset
}

func main() {
	key := GetItemKey(9991, 99991)
	fmt.Println("key ", key)

	id, num := ParseItemKey(key)
	fmt.Println("id ", id, "num ", num)
	//stack.PrintGoroutineMemConsume()
	//channeltest.RunUnbounded()
	//go timertest.RunTimerAfter()

	//test2()

	//http.ListenAndServe("0.0.0.0:8899", nil)

}

func test() {
	test2()
}

func test2() {
	pc, file, line, _ := runtime.Caller(3)
	f := runtime.FuncForPC(pc)
	log.Println("pc3 ", pc, "file ", file, "line", line, "name", f.Name())

	pc, file, line, _ = runtime.Caller(2)
	f = runtime.FuncForPC(pc)
	log.Println("pc2 ", pc, "file ", file, "line", line, "name", f.Name())

	pc, file, line, _ = runtime.Caller(1)
	f = runtime.FuncForPC(pc)
	log.Println("pc1 ", pc, "file ", file, "line", line, "name", f.Name())

	pc, file, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.Println("pc0 ", pc, "file ", file, "line", line, "name", f.Name())

	CloseMsgSyn := make(chan struct{})

	go func() {
		for {
			select {
			case <-CloseMsgSyn:
				fmt.Println("close1")
				return
			default:
				fmt.Println("aa")
			}
		}
	}()

	go func() {
		for {
			select {
			case <-CloseMsgSyn:
				fmt.Println("close2")
				return
			default:
				fmt.Println("aa")
			}
		}
	}()

	time.Sleep(10)
	close(CloseMsgSyn)
}
