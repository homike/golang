package base

import (
	"fmt"
	"testing"
	"time"
)

func testDefer() {
	defer fmt.Println(1)
	defer fmt.Println(2)
	defer fmt.Println(3)

}

func testDeferRecover() {
	/*
	 */
	exit := make(chan int)
	go func() {
		/*
			defer func() {
				if x := recover(); x != nil {
					fmt.Println("recover: ", x)
				}
			}()
		*/
		defer func() {
			fmt.Println("exit")
			exit <- 1
		}()
		for i := 0; i < 100; i++ {
			fmt.Println("i: ", i)
			time.Sleep(100 * time.Millisecond)
			if i == 10 {
				arr := []int{}
				fmt.Println(arr[2])
			}
		}
	}()
	<-exit
	fmt.Println("====END====")
	for i := 0; i < 10; i++ {
		fmt.Println("j: ", i)
	}
}

func Test_Defer(t *testing.T) {
	testDeferRecover()
}
