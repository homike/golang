// map_test.go 测试golang中的map的并发读写能力
package base

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 并发读写会panic
func _Test_MapWriteRead(t *testing.T) {
	c := make(map[string]int)
	go func() { //开一个协程写map
		for j := 0; j < 1000000; j++ {
			c[fmt.Sprintf("%d", j)] = j
		}
	}()
	go func() { //开一个协程读map
		for j := 0; j < 1000000; j++ {
			_ = c[fmt.Sprintf("%d", j)]
		}
	}()

	time.Sleep(time.Second * 20)
}

// 并发读不会panic
func _Test_MapRead(t *testing.T) {
	c := make(map[string]int)
	for j := 0; j < 1000000; j++ {
		c[fmt.Sprintf("%d", j)] = j
	}

	for i := 0; i < 10; i++ {
		go func() { //开一个协程读map
			for j := 0; j < 1000000; j++ {
				_ = c[fmt.Sprintf("%d", j)]
			}
		}()
	}

	time.Sleep(time.Second * 10)
}

// 并发写会panic
func _Test_MapWrite(t *testing.T) {
	c := make(map[string]int)

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 1000000; j++ {
				c[fmt.Sprintf("%d", j)] = j
			}
		}()
	}
	time.Sleep(time.Second * 20) //让执行main函数的主协成等待20s,不然不会执行上面的并发操作
}

func _Test_Map(t *testing.T) {
	// var counter = struct {
	// 	sync.RWMutex
	// 	m map[string]int
	// }{m: make(map[string]int)}
}

type myStruct1 struct {
	value int
	name  string
}

func modiftyMap(m map[int]myStruct1) {
	m[0] = myStruct1{value: 100}
	v, _ := m[1]
	v.value = 200
}
func _Test_Map_Args(t *testing.T) {
	//baseSlice = []*myStruct{{value: 1}, {value: 2}}
	m := make(map[int]myStruct1)
	m[0] = myStruct1{value: 1}
	m[1] = myStruct1{value: 2}

	//fmt.Printf("one：%p\n", &baseSlice)
	modiftyMap(m)

	fmt.Println("------base-------")
	for _, v := range m {
		fmt.Printf(" %v", v)
	}
}

/*
func _TestMapCopy(t *testing.T) {
	type mapStruct struct {
		maps map[int32]int32
		data int
	}
	map1 := &mapStruct{
		maps: make(map[int32]int32),
		data: 2,
	}
	map1.maps[1] = 1
	map1.data = 1
	map2 := &map1
	map2.maps[2] = 2
	map2.data = 2

	fmt.Println("map1 ", map1)
	fmt.Println("map2 ", map2)
}
*/

var count = 100

type MemberPoint struct {
	Value int
}

type Table struct {
	M1  *MemberPoint
	M2  *MemberPoint
	M3  *MemberPoint
	M4  *MemberPoint
	M5  *MemberPoint
	M6  *MemberPoint
	M8  *MemberPoint
	M9  *MemberPoint
	M10 *MemberPoint

	M11 *MemberPoint
	M12 *MemberPoint
	M13 *MemberPoint
	M14 *MemberPoint
	M15 *MemberPoint
	M16 *MemberPoint
	M18 *MemberPoint
	M19 *MemberPoint
	M20 *MemberPoint

	M21 *MemberPoint
	M22 *MemberPoint
	M23 *MemberPoint
	M24 *MemberPoint
	M25 *MemberPoint
	M26 *MemberPoint
	M28 *MemberPoint
	M29 *MemberPoint
	M30 *MemberPoint

	M41 *MemberPoint
	M42 *MemberPoint
	M43 *MemberPoint
	M44 *MemberPoint
	M45 *MemberPoint
	M46 *MemberPoint
	M48 *MemberPoint
	M49 *MemberPoint
	M40 *MemberPoint

	M51 *MemberPoint
	M52 *MemberPoint
	M53 *MemberPoint
	M54 *MemberPoint
	M55 *MemberPoint
	M56 *MemberPoint
	M58 *MemberPoint
	M59 *MemberPoint
	M50 *MemberPoint

	M61 *MemberPoint
	M62 *MemberPoint
	M63 *MemberPoint
	M64 *MemberPoint
	M65 *MemberPoint
	M66 *MemberPoint
	M68 *MemberPoint
	M69 *MemberPoint
	M60 *MemberPoint

	M71 *MemberPoint
	M72 *MemberPoint
	M73 *MemberPoint
	M74 *MemberPoint
	M75 *MemberPoint
	M76 *MemberPoint
	M78 *MemberPoint
	M79 *MemberPoint
	M70 *MemberPoint

	M81 *MemberPoint
	M82 *MemberPoint
	M83 *MemberPoint
	M84 *MemberPoint
	M85 *MemberPoint
	M86 *MemberPoint
	M88 *MemberPoint
	M89 *MemberPoint
	M80 *MemberPoint

	M91 *MemberPoint
	M92 *MemberPoint
	M93 *MemberPoint
	M94 *MemberPoint
	M95 *MemberPoint
	M96 *MemberPoint
	M98 *MemberPoint
	M99 *MemberPoint
	M90 *MemberPoint

	M101 *MemberPoint
	M102 *MemberPoint
	M103 *MemberPoint
	M104 *MemberPoint
	M105 *MemberPoint
	M106 *MemberPoint
	M108 *MemberPoint
	M109 *MemberPoint
	M100 *MemberPoint
}

func Benchmark_FindByMap(b *testing.B) {
	//s := []int{}
	m := make(map[int]*MemberPoint)
	for i := 0; i < count; i++ {
		//s = append(s, i)
		m[i] = &MemberPoint{Value: i}
	}
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	// 开始计时器
	b.StartTimer()

	index := rand.Intn(count)
	fmt.Println(m[index])
}

func Benchmark_FindByPointer(b *testing.B) {
	m := &Table{
		M100: &MemberPoint{Value: 100},
	}
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	// 开始计时器
	b.StartTimer()

	fmt.Println(m.M100.Value)
}

func _Benchmark_Slice(b *testing.B) {
	s := []int{}
	//m := make(map[int]int)
	for i := 0; i < count; i++ {
		s = append(s, i)
		//m[i] = i
	}
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	// 开始计时器
	b.StartTimer()

	index := rand.Intn(count)
	_ = s[index]
}

func _Benchmark_Map(b *testing.B) {
	//s := []int{}
	m := make(map[int]int)
	for i := 0; i < count; i++ {
		//s = append(s, i)
		m[i] = i
	}
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	// 开始计时器
	b.StartTimer()

	index := rand.Intn(count)
	_ = m[index]
}

func _Test_InsertMap(t *testing.T) {
	m1 := map[int32][]int32{}
	m2 := map[int32][]int32{}

	m1[1] = append(m1[1], 20)

	_, ok := m2[1]
	if !ok {
		m2[1] = []int32{20}
	} else {
		m2[1] = append(m2[1], 20)
	}

	for _, v := range m1 {
		fmt.Println("m1: ", v)
	}

	for _, v := range m2 {
		fmt.Println("m2: ", v)
	}
}

func _Test_InsertMap2(t *testing.T) {
	m1 := map[int32][]int32{}
	m2 := map[int32][]int32{}

	arr := []int32{10, 20, 30}

loop:
	for _, v := range arr {
		if v == 20 {
			continue loop
		}
		m1[1] = append(m1[1], v)
	}

	for _, v := range arr {
		if v == 20 {
			continue
		}
		_, ok := m2[1]
		if !ok {
			m2[1] = []int32{v}
		} else {
			m2[1] = append(m2[1], v)
		}
	}

	for _, v := range m1 {
		fmt.Println("m1: ", v)
	}

	for _, v := range m2 {
		fmt.Println("m2: ", v)
	}
}

func Test_InsertMap2(t *testing.T) {
	m := make(map[int][]int)
	m[1] = []int{2}

	_, ok := m[1]
	fmt.Println("ok1: ", ok, "len: ", len(m[1]))
	m[1] = append(m[1][:0], m[1][1:]...)

	_, ok = m[1]
	fmt.Println("ok1: ", ok, "len: ", len(m[1]))
}
