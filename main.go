package main

import (
	_ "GoTest/astar"
	_ "GoTest/basetest"
	_ "GoTest/channeltest"

	_ "gotest/bittest"
	// _ "github.com/coreos/etcd/clientv3"
	// _ "google.golang.org/grpc"
)

// 测试提交
// Main
func match(input int) int {
	arr := []int{1, 3, 5}
	for _, v := range arr {
		if v == input {
			return v
		}
	}
	return 0
}

func main() {
	/*
		ids := []int{1, 2, 3, 4}
		weights := []int{0, 0, 0}

		times := []int{0, 0, 0}
		for i := 0; i < 10000; i++ {
			randomID := random.GetRandomWeightID(ids, weights)
			//fmt.Printf("%v, ", randomID)

			for j := 0; j < len(ids); j++ {
				if ids[j] == randomID {
					times[j]++
				}
			}
		}

		fmt.Printf("times: $v \n", times)
	*/

	/*
		testValue := 1
		switch testValue {
		case 1:
			fmt.Println("1")
		case 2:
			fmt.Println("2")
		case match(testValue):
			fmt.Println("3")
		}
	*/

	//	memorymodel.Run1()

	//gorutineMap := make(map[int]string)
	//for i := 0; i < 100; i++ {
	//	gorutineMap[i] = fmt.Sprintf("gorutine%v", i)
	//}

	//var ws sync.WaitGroup
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))

	//for i := 0; i < 10000; i++ {
	//	go func() {
	//		ws.Add(1)
	//		index := r.Intn(99)
	//		for {
	//			str := fmt.Sprintf("gorutine%v", index)
	//			if str != gorutineMap[index] {
	//				fmt.Println("data ", index, gorutineMap[index])
	//			}
	//		}
	//	}()
	//}
	//ws.Wait()

	//fmt.Println("End")
}
