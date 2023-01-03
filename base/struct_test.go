package base

import (
	"fmt"
	"testing"
	"unsafe"
)

type emptyStruct struct{}

// 空结构体的变量的内存地址都是一样的
func TestStruct(t *testing.T) {
	a := struct{}{}
	b := struct{}{}
	c := emptyStruct{}

	fmt.Printf("a: %p, b: %p, c: %p", &a, &b, &c)
}

// 结构体 和 结构体指针区别
type A struct {
	i int
}

func (a *A) print() {
	fmt.Println("print a: ", a.i)
}

func (a A) print2() {
	fmt.Println("print2 a: ", a.i)
}

type B1 struct {
	A
}

type B2 struct {
	*A
}

func TestStructOrStructPointer(t *testing.T) {
	b1 := &B1{}
	b2 := &B2{}

	// 1. 结构体指针默认是nil, 直接调用会panic,
	// 结构体默认不为nil, 可以直接调用
	b1.print()
	//b2.print()

	b2.A = &A{}
	b1.print2()
	b2.print2()

	// 2. 内存占用是一样的
	fmt.Println("b1: ", unsafe.Sizeof(b1))
	fmt.Println("b2: ", unsafe.Sizeof(b2))

}
