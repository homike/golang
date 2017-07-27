package basetest

import (
	"fmt"
)

type myStruct struct {
	value int
	name  string
}

func (m *myStruct) test1() {
	m.value = 1
}
func (m myStruct) test2() {
	m.value = 2
}
func Run2() {
	testS := myStruct{
		value: 0,
		name:  "test",
	}
	testS.test1()
	fmt.Printf("*T %v\n", testS)
	testS.test2()
	fmt.Printf("T %v\n", testS)
}
