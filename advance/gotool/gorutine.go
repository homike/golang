package main

import (
	"fmt"
	"net/http"
	"time"

	_ "net/http/pprof"
)

func GorutineLeak() {
	time.AfterFunc(15*time.Second, func() {
		fmt.Println("Time Out")
	})

	go func() {
		for i := 0; i < 100000; i++ {
			go func() {
				for {
					//runtime.GC()
					time.Sleep(1 * time.Second)
				}
			}()
			time.Sleep(2 * time.Second)
		}
	}()

	fmt.Println("serve listen")

	http.ListenAndServe(":80", nil)
}
