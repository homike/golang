package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/homike/gotest/timertest"
)

func main() {
	go timertest.RunTimerAfter()

	test2()

	http.ListenAndServe("0.0.0.0:8899", nil)

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
