package timertest

import (
	"fmt"
	"time"
)

// timer.After 在for 循环中需要谨慎
// 查看timer.After 的说明，timer到期前是不会释放的
// 如果大量请求到来，会导致大量timer无法释放
func useTimerAfterWrong(queue <-chan string) {
	for {
		select {
		case <-queue:
			//fmt.Println(queue)
		case <-time.After(3 * time.Minute):
			return
		}
	}
}

func useTimerAfterRight(queue <-chan string) {
	intel := time.NewTimer(3 * time.Minute)
	defer intel.Stop()

	for {
		select {
		case <-queue:
			//fmt.Println(queue)
			return
		case <-intel.C:
			return
		}
	}
}

func RunTimerAfter() {
	queue := make(chan string, 100)
	go useTimerAfterWrong(queue)
	for i := 0; i < 600000; i++ {
		queue <- fmt.Sprintf("queue")
	}
}
