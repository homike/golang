package channeltest

import (
	"fmt"
	"sync"
	"time"
)

type TestR struct {
	Value         int
	NextStepID    int
	NextStartTime int64
	Step          int
}

var reqQueue chan struct{}
var cStartTime, nStartTime int64

func producer(count int) chan *TestR {
	cStartTime = time.Now().UnixNano() / 1000000
	nStartTime = cStartTime + (int64)(1000*2)

	out := make(chan *TestR, 2000)
	go func() {
		for n := 0; n < count; n++ {
			test := &TestR{Value: n,
				NextStepID:    0,
				NextStartTime: time.Now().UnixNano()/1000000 + (int64)(maRand.New(maRand.NewSource(time.Now().UnixNano())).Intn(1000*2)),
				Step:          0,
			}
			out <- test
		}
		//close(out)
	}()
	return out
}

func consumer(robots chan *TestR) {

	fmt.Println("Start BatchDo")
	for r := range robots {
		r := r
		if r.NextStartTime <= (time.Now().UnixNano() / 1000000) {
			reqQueue <- struct{}{}
			//fmt.Println("Enter")
			go func() {
				r.NextStartTime = nStartTime + (int64)(maRand.New(maRand.NewSource(time.Now().UnixNano())).Intn(1000*2))

				fmt.Println("exec ", r.Value, ":", r.Step)
				r.Step = r.Step + 1
				r.Step = r.Step % 5

				<-reqQueue
				//fmt.Println("Exist")

				robots <- r
			}()
		} else {
			robots <- r
		}

		if (time.Now().UnixNano() / 1000000) > nStartTime {
			lock := &sync.Mutex{}
			lock.Lock()
			cStartTime = nStartTime
			nStartTime = cStartTime + (int64)(1000*2)
			lock.Unlock()
		}

	}

	fmt.Println("End BatchDo")
}

func main() {
	reqQueue = make(chan struct{}, 2)

	robots := producer(5)
	consumer(robots)
}
