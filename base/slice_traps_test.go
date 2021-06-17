package base

import (
	"fmt"
	"testing"
)

type StructA struct {
	Name  string
	Slice []*myStruct
}

var baseStruct *StructA

func init() {
	baseStruct = &StructA{
		Name:  "base",
		Slice: make([]*myStruct, 0, 10),
	}
	baseStruct.Slice = append(baseStruct.Slice, &myStruct{value: 1})
	baseStruct.Slice = append(baseStruct.Slice, &myStruct{value: 2})
}

func returnSlice() []*myStruct {
	retSlice := make([]*myStruct, len(baseSlice), len(baseSlice))
	copy(retSlice, baseSlice)
	//return append([]*myStruct{}, baseSlice...)
	return retSlice
}

func returnSliceNoCopy() []*myStruct {
	return baseStruct.Slice
}

// 1. slice copy之后返回，下标访问是否会影响原slice
// 2. slice 作为参数传递，下标访问是否会影响原slice
func _Test_SliceTrap_ReturnValue(t *testing.T) {
	slice := returnSliceNoCopy()
	slice[0] = &myStruct{value: 100}
	//slice[1] = &myStruct{value: 200}
	slice[1].value = 200

	fmt.Println("------base-------")
	for _, v := range baseStruct.Slice {
		fmt.Printf(" %v", v)
	}
}

func _Test_SliceTrap_Args(t *testing.T) {
	//baseSlice = []*myStruct{{value: 1}, {value: 2}}
	testSlice := []*myStruct{{value: 1}, {value: 2}}
	//fmt.Printf("one：%p\n", &baseSlice)
	//modiftySlice(baseSlice)
	modiftySlice(testSlice)

	fmt.Println("------base-------")
	for _, v := range testSlice {
		fmt.Printf(" %v", v)
	}
}

func SliceReturn(a []int) []int {
	return a
}
