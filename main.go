package main

import (
	_ "GoTest/astar"
	_ "GoTest/basetest"
	_ "GoTest/channeltest"
	"fmt"

	_ "gotest/bittest"
	// _ "github.com/coreos/etcd/clientv3"
	// _ "google.golang.org/grpc"

	"gotest/memorymodel"
)

// 测试提交
// Main

func main() {

	rewards := []int{1, 2, 3, 2}
	probs := []int{100, 200, 300, 200}

	//for i := len(rewards) - 1; i >= 0; i-- {
	for i := 0; i < len(rewards); i++ {
		if rewards[i] == 2 || rewards[i] == 3 {
			rewards = append(rewards[:i], rewards[i+1:]...)
			probs = append(probs[:i], probs[i+1:]...)
			fmt.Printf("index: %v, len: %v \n", i, len(rewards))
		}
	}
	fmt.Printf("rewads: %v, probs: %v \n", rewards, probs)

	// astar.RunAstar()
	// basetest.Run2()

	// reflect.RunReflect1()
	// bittest.RunBit()
	// interfacetest.RunInterface()
	// thirdparty.RunRandom()

	memorymodel.Run1()
}
