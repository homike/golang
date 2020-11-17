package base

import (
	"fmt"
	"testing"
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

type sliceType []int

func Delete(s *sliceType) {
	*s = append((*s)[:0], (*s)[1:]...)
}

func Test_Slice(t *testing.T) {
	testS := myStruct{
		value: 0,
		name:  "test",
	}
	testS.test1()
	fmt.Printf("*T %v\n", testS)
	testS.test2()
	fmt.Printf("T %v\n", testS)

	slice1 := sliceType{1, 2, 3, 4}
	Delete(&slice1)

	for _, v := range slice1 {
		fmt.Println("value ", v)
	}
	fmt.Print("end \n ")
}
