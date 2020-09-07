package actor

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// 测试Go, Call是否正常, 是否阻塞非阻塞
func _TestActorBase(t *testing.T) {
	client := NewClient(0)

	t.Logf("%s, start ", time.Now().Format("15:04:05.000"))

	client.Go("key", "value", nil)
	t.Logf("%s, Go1 end", time.Now().Format("15:04:05.000"))

	reply := 0
	err := client.Call("key", "value", &reply)
	t.Logf("%s, Call end, reply: %v, err: %v", time.Now().Format("15:04:05.000"), reply, err)
}

// 测试Actor模式下, 能否正常退出
// 1. actor close 之后, 是否有消息阻塞不能正常退出
// 2. actor close 之后, 消息队列是否成功清空
func _TestActorClose(t *testing.T) {
	client := NewClient(0)

	t.Logf("%s, start ", time.Now().Format("15:04:05.000"))

	go func() {
		time.Sleep(2 * time.Second)
		client.Close()
	}()

	var waitGroup sync.WaitGroup
	for i := 0; i < 10; i++ {
		waitGroup.Add(1)
		go func(goIndex int) {
			defer waitGroup.Done()

			curIndex := 0
			for {
				key := fmt.Sprintf("%v_%v", goIndex, curIndex)
				reply := 0
				err := client.Call(key, "value", &reply)
				if err != nil {
					fmt.Printf("%s, Call Exit, goroutine index: %v, err: %v \n", time.Now().Format("15:04:05.000"), goIndex, err)
					return
				}
				//t.Logf("%s, Call end, reply: %v, err: %v", time.Now().Format("15:04:05.000"), reply, err)
				curIndex++
			}
		}(i)
	}

	waitGroup.Wait()
	t.Logf("Close Success")
}

type Requset struct {
	A int
}

type Reply struct {
	B int
	C string
}

//type Handle struct{}

func MessageHandle(req *Requset) *Reply {
	fmt.Println("MessageHandle: ", req.A)

	return &Reply{
		B: req.A + 1,
		C: "Success",
	}
}

const CMD_REQUEST_TEST = 1

func TestMessageTrans(t *testing.T) {

	client := NewClient(0)

	// register
	client.RegisterCall("handle", MessageHandle)

	// exec
	req := &Requset{
		A: 1,
	}
	ret := client.ExecFunc("handle", req)

	reply := ret[0].Interface().(*Reply)
	//reply := ret[0].(*Reply)
	//reply, _ := interface{}(ret[0]).(*Reply)
	fmt.Println("[Result]", reply.B, "  ", reply.C)
	/*
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			reply := 0
			err := client.Call(CMD_REQUEST_TEST, "value", &reply)
			if err != nil {
				fmt.Printf("%s, Call Exit, err: %v \n", time.Now().Format("15:04:05.000"), err)
				return
			}
		}()

		waitGroup.Wait()
	*/
}
