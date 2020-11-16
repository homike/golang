package interfacetest

import (
	"fmt"
)

////////////////////////////////////
// 将方法 ==》转换成接口
///////////////////////////////////

type Handler interface {
	ServeHTTP(int, int)
}

type Entry struct {
	h Handler
}

// 方法一
func RegisterStruct(handler Handler) *Entry {
	return &Entry{
		h: handler,
	}
}

type UserInt struct {
}

func (u *UserInt) ServeHTTP(v1 int, v2 int) {
	fmt.Println("result2 : ", v1, v2)
}

// 方法二
type HandlerFunc func(int, int)

func (handler HandlerFunc) ServeHTTP(v1 int, v2 int) {
	handler(v1, v2)
}

func RegisterFunc(handler Handler) *Entry {
	return &Entry{
		h: handler,
	}
}

func UserFunc(v1 int, v2 int) {
	fmt.Println("result1 : ", v1, v2)
}

// 方法二 相对于 方法一, 不需要再额外定义一个UserInt 结构体，再用结构体实现ServeHTTP接口
// 但是为什么不在Entry 中，直接将 handler 替换为函数指针了？
//// ==> 函数指针不够灵活, 接口可以只想一个函数，也可以直接实现了ServeHTTP接口的结构体
//// 参考net/http

func Run3() {
	entery := RegisterStruct(&UserInt{})
	entery.h.ServeHTTP(1, 2)

	entery1 := RegisterFunc(HandlerFunc(UserFunc))
	entery1.h.ServeHTTP(3, 4)
}
