package base

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/mohae/deepcopy"
)

// slice 底层结构
// go/src/runtime/slice.go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

// 当slice 的 cap没有超过array 指向底层数组的cap时, slice中的array不会发生改变
// 此时创建出的slice, 仅仅是通过len 来控制对array的访问范围

// ------------ Append --------------
// 1. append 当没有超过array的长度时, 只是改变了slice的 len
func _Test_Slice_Append(t *testing.T) {
	{
		//							  array
		a := make([]int, 1, 2) //a: | 0	| 0 |, len: 1, cap: 2
		a[0] = 1               //a: | 0	| 0 |, len: 1, cap: 2
		b := append(a, 1)      //b: | 0	| 1 |, len: 2, cap: 2
		c := append(a, 2)      //c: | 0	| 2 |, len: 2, cap: 2
		//d := append(a, []int{2, 3}...)
		fmt.Println(a, b, c)
	}
}

// ------------ Function Args --------------
// slice作为参数传递时，是值拷贝, len, cap都是独立的
// 但是array是指针，即使值拷贝，指向的仍然是相同的底层数组
// 当数据没有超过底层数据的cap时, slice中的array不会发生改变

// 2. 对于slice的append操作, 创建了一个新的切片, 修改的还是底层数组
func _Test_Slice_FuncArgs(t *testing.T) {
	fn := func(in []int) {
		in = append(in, 5)
		// "5"写入了slice对应的底层数组, 只是因为slice传递的是值拷贝,
		// 因此 slice的 len 未改变，导致slice无法查看到"5"
	}

	slice := make([]int, 0, 10)
	slice = append(slice, 1)
	fmt.Println(slice, len(slice), cap(slice))

	fn(slice)

	fmt.Println(slice, len(slice), cap(slice))
	s1 := slice[0:9] //数组截取, 此时会打印出5
	fmt.Println(s1, len(s1), cap(s1))
}

// 3. 对 in 进行下标操作, 相当于操作了slice.array指向的数组
// slice被更改
func Test_Slice_FuncArgs2(t *testing.T) {
	fn := func(in []int) {
		in[0] = 100
	}
	s := make([]int, 0, 1)
	s = append(s, 1)

	fn(s)

	fmt.Println(s, len(s), cap(s))
}

type S struct {
	value int
}

var ss []*S

func init() {
	ss = make([]*S, 0, 1000)
	for i := 0; i < 1000; i++ {
		ss = append(ss, &S{value: i})
	}
}

// ------------ benchmark --------------
func Benchmark_CopyByCopy(b *testing.B) {
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	//dst := []*myStruct{}
	// 开始计时器
	b.StartTimer()

	dst := make([]*S, 1000, 1000)
	copy(dst, ss)

	fmt.Println("len(dst): ", len(dst), len(ss))
}

func Benchmark_CopyByAppend(b *testing.B) {
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	//dst := []*myStruct{}
	// 开始计时器
	b.StartTimer()

	dst := make([]*S, 1000, 1000)
	for _, v := range ss {
		dst = append(dst, &S{value: v.value})
	}

	fmt.Println("len(dst): ", len(dst), len(ss))
}

func Benchmark_CopyByDeepCopy(b *testing.B) {
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	//dst := []*myStruct{}
	// 开始计时器
	b.StartTimer()

	//dst := make([]*S, 1000, 1000)
	dst, _ := deepcopy.Copy(ss).([]*S)
	fmt.Println("len(dst): ", len(dst), len(ss))
}
