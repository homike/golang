package basetest

import (
	"fmt"
)

////////////////////////////////////
// 可以理解为方法多态，一个函数指针可以调用不同的方法
///////////////////////////////////
type FuncHandler func(int, int)

func (handler FuncHandler) ServerHTTP(v1 int, v2 int) {
	handler(v1, v2)
}

func NoramlFunc(v1 int, v2 int) {
	fmt.Println("result : ", v1, v2)
}

func TestRun(handler FuncHandler) {
	handler.ServerHTTP(1, 2)
}

func Run3() {
	handler := (FuncHandler)(NoramlFunc)
	TestRun(handler)
}
