package base

import (
	"fmt"
	"testing"
)

/*

import (
	"fmt"
	"testing"
)

type ArrStruct struct {
	value int
}

var baseSlice [2]*ArrStruct

func returnArr() [2]*ArrStruct {
	baseSlice = [2]*ArrStruct{{value: 1}, {value: 2}}
	//return append([]*ArrStruct{}, baseSlice...)
	return baseSlice
}

func modiftyArr(s []*ArrStruct) {
	//s = append(s, &ArrStruct{value: 2})
	s[0] = &ArrStruct{value: 100}
	s[1].value = 200
}

func Test_SliceTrap_ReturnValue(t *testing.T) {
	slice := returnArr()
	slice[0] = &ArrStruct{value: 100}
	slice[1] = &ArrStruct{value: 200}

	fmt.Println("------base-------")
	for _, v := range baseSlice {
		fmt.Printf(" %v", v)
	}
}

func Test_SliceTrap_Args(t *testing.T) {
	baseSlice = [2]*ArrStruct{{value: 1}, {value: 2}}
	modiftyArr(baseSlice[:])
	modiftyArr(baseSlice)

	fmt.Println("------base-------")
	for _, v := range baseSlice {
		fmt.Printf(" %v", v)
	}
}
*/
var arr [2]int

func returnTest() []int {
	return arr[:]
}

func TestArrayToSlice(t *testing.T) {
	arr = [2]int{1, 2}
	sli := returnTest()
	sli[0] = 100

	for _, v := range arr {
		fmt.Print(v, ",")
	}
}
