package interfacetest

import "fmt"

type AInterface interface {
	Func1()
	Func2()
}

type A struct {
	Data int
}

func (self *A) Func1() {
	fmt.Println("A Func1")
}

func (self *A) Func2() {
	fmt.Println("A Func2")
}

type AA struct {
	Data1 int
	A
}

type BB struct {
	AA
	Data2 int
}

func (self *BB) Func1() {
	fmt.Println("BB Func1()")
}

func RunPoly() {
	bb := &BB{
		AA: AA{
			Data1: 2,
		},
		Data2: 3,
	}

	aInterface := (interface{}(bb)).(AInterface)
	aInterface.Func1()
	aInterface.Func2()
}
