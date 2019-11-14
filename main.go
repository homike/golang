package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/homike/gotest/timertest"
)

func main() {
	go timertest.RunTimerAfter()

	http.ListenAndServe("0.0.0.0:8899", nil)

}
