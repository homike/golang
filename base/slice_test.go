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

func _Test_Slice(t *testing.T) {
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

func modiftySlice(s []*myStruct) {
	//s = append(s, &myStruct{value: 2})
	s[0] = &myStruct{value: 100}
	s[1].value = 200
}

// slice陷阱: slice的传递是值拷贝，传递的是不同的slice，
// 只是他们指向的内容是相同的
func _Test_SliceTrap(t *testing.T) {
	baseSlice := []*myStruct{{value: 1}, {value: 2}}
	modiftySlice(baseSlice)

	fmt.Println("------base-------")
	for _, v := range baseSlice {
		fmt.Printf(" %v", v)
	}
}

var baseSlice []*myStruct

func returnSlice() []*myStruct {
	baseSlice = []*myStruct{{value: 1}, {value: 2}}
	retSlice := make([]*myStruct, len(baseSlice), len(baseSlice))
	copy(retSlice, baseSlice)
	//return append([]*myStruct{}, baseSlice...)
	return retSlice
}

func _Test_SliceTrap_ReturnValue(t *testing.T) {
	slice := returnSlice()
	slice[0] = &myStruct{value: 100}
	slice[1].value = 200

	fmt.Println("------base-------")
	for _, v := range baseSlice {
		fmt.Printf(" %v", v)
	}
}

func _Test_SliceTrap_Args(t *testing.T) {
	baseSlice = []*myStruct{{value: 1}, {value: 2}}
	modiftySlice(baseSlice)

	fmt.Println("------base-------")
	for _, v := range baseSlice {
		fmt.Printf(" %v", v)
	}
}

func init() {
	for i := 0; i < 800; i++ {
		baseSlice = append(baseSlice, &myStruct{value: i})
	}
}

func Benchmark_Copy(b *testing.B) {
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	//dst := []*myStruct{}
	// 开始计时器
	b.StartTimer()

	dst := make([]*myStruct, 800, 1000)
	copy(dst, baseSlice)

	fmt.Println("len(dst): ", len(dst), len(baseSlice))
}

func Benchmark_Copy2(b *testing.B) {
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	//dst := []*myStruct{}
	// 开始计时器
	b.StartTimer()

	dst := make([]*myStruct, 0, 1000)
	for _, v := range baseSlice {
		dst = append(dst, &myStruct{value: v.value})
	}

	fmt.Println("len(dst): ", len(dst), len(baseSlice))
}

/*
func Benchmark_Append(b *testing.B) {
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	//for i := 0; i < 1000; i++ {
	//	baseSlice = append(baseSlice, &myStruct{value: i})
	//}
	// 开始计时器
	b.StartTimer()

	//dst := make([]*myStruct, 0, 1000)
	dst := []*myStruct{}
	for _, v := range baseSlice {
		dst = append(dst, v)
	}

	fmt.Println("len(dst): ", len(dst), len(baseSlice))
}
*/
