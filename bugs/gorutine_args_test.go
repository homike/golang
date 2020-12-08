package gobugs

import (
	"fmt"
)

// gorutine的参数传递, 同步参数列表传递，不要直接使用局部变量，会导致意想不到的问题
func Test_GorutineArgs() {
	for i := 17; i <= 21; i++ {
		/*
			go func() {
				_ := fmt.Sprintf("%d", i)
			}()
		*/
		fmt.Println("start1")
		go func(v int) {
			fmt.Println("end")
			fmt.Printf("%d \n", v)
		}(i)
	}
}
