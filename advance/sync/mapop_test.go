// https://tonybai.com/2020/11/10/understand-sync-map-inside-through-examples/
// map 总体来说比较适合读多写少的情况
// dirtyLocked: 当Store新元素时, read 又没有dirty的完整数据时, 会将read->dirty
// missLocked:	当miss次数大于dirty长度, 将dirty提升为read, 清空dirty
// 数据删除:	如果只在read中删除的元素，下一次missLoced就会删除
//				如果在dirty中的元素，马上删除
//				如果在read 和 dirty中都有的元素, 会在两次missLoced后删除
// 不管是dirtyLocked, 还是missLocked, 都会有额外的消耗, 因此对于读到写少的场景才比较适用

package sync_test

import (
	"fmt"
	xsync "gotest/advance/sync"
	"testing"
)

func _Test_Init(t *testing.T) {
	var m xsync.Map
	fmt.Println("sync.Map init status:")
	m.Dump()
}

type val struct {
	s string
}

// 当load时, miss次数大于等于dirty长度时, 将dirty promoted为 read
func _Test_StoreLoad(t *testing.T) {
	var m xsync.Map
	fmt.Println("sync.Map init status:")
	m.Dump()

	// store
	m.Store("key1", &val{"val1"})
	fmt.Println("\nafter store key1:")
	m.Dump()

	// load
	m.Load("key2") //这里我们尝试load key="key2"
	fmt.Println("\nafter load key2:")
	m.Dump()
}

func _Test_StoreLoad2(t *testing.T) {
	var m xsync.Map

	m.Store("key1", &val{"val1"})
	fmt.Println("\nafter store key1:")
	m.Dump()

	m.Load("key2")
	fmt.Println("\nafter load key2:")
	m.Dump()

	m.Store("key2", &val{"val2"})
	fmt.Println("\nafter store key2:")
	m.Dump()

	m.Store("key3", &val{"val3"})
	fmt.Println("\nafter store key3:")
	m.Dump()

	m.Load("key1")
	fmt.Println("\nafter load key1:")
	m.Dump()

	m.Load("key2")
	fmt.Println("\nafter load key2:")
	m.Dump()

	m.Load("key2")
	fmt.Println("\nafter load key2 2nd:")
	m.Dump()

	m.Load("key2")
	fmt.Println("\nafter load key2 3rd:")
	m.Dump()
}

func _Test_StoreLoad3(t *testing.T) {
	var m xsync.Map

	m.Store("key1", &val{"val1"})
	fmt.Println("\nafter store key1:")
	m.Dump()

	m.Load("key2")
	fmt.Println("\nafter load key2:")
	m.Dump()

	m.Store("key2", &val{"val2"})
	fmt.Println("\nafter store key2:")
	m.Dump()

	m.Load("key2")
	m.Load("key2")
	m.Dump()

	m.Store("key3", &val{"val3"})
	fmt.Println("\nafter store key3:")
	m.Dump()

	//m.Load("key2")
	//m.Load("key2")
	//m.Load("key2")
	//fmt.Println("\nafter load key2 3rd:")
	//m.Dump()
}

func _Test_Update(t *testing.T) {
	var m xsync.Map

	val1 := &val{"val1"}
	m.Store("key1", val1)
	fmt.Println("\nafter store key1:")
	m.Dump()

	val2 := &val{"val2"}
	m.Store("key2", val2)
	fmt.Println("\nafter store key2:")
	m.Dump()

	val2_1 := &val{"val2_1"}
	m.Store("key2", val2_1)
	fmt.Println("\nafter update key2(in read, not in dirty):")
	m.Dump()
}

// 当deleteInRead, deleteInAll时, 只是将value设置为了nil, 并没有立刻删除整个key-value对
// 只有在load时, 当miss > len(dirty) 才会将dirty==nil, 内存才会释放
// 导致当使用指针最为key时, 可能有内存泄漏
func _Test_Delete(t *testing.T) {
	//deleteInRead()

	//deleteInDirty()

	deleteInAll()
}

// 将read 和 dirty 中的元素设置为nil
func deleteInRead() {
	var m xsync.Map

	val1 := &val{"val1"}
	m.Store("key1", val1)

	val2 := &val{"val2"}
	m.Store("key2", val2)

	val3 := &val{"val3"}
	m.Store("key3", val3)

	m.Load("key1")
	m.Load("key1")
	m.Load("key1")

	m.Delete("key2")
	fmt.Println("\nafter delete key2:")
	m.Dump()
}

// 直接删除dirty中的元素
func deleteInDirty() {
	var m xsync.Map

	val1 := &val{"val1"}
	m.Store("key1", val1)

	val2 := &val{"val2"}
	m.Store("key2", val2)

	val3 := &val{"val3"}
	m.Store("key3", val3)

	m.Delete("key2")
	fmt.Println("\nafter delete key2:")
	m.Dump()
}

// 将read 和 dirty 中的元素都设置为nil
func deleteInAll() {
	var m xsync.Map

	m.Store("key1", &val{"val1"})
	m.Store("key2", &val{"val2"})

	m.Load("key1")
	m.Load("key1")

	m.Store("key3", &val{"val3"})
	m.Dump()

	m.Delete("key1")
	fmt.Println("\nafter delete key2:")
	m.Dump()

	m.Load("key3")
	m.Load("key3")
	m.Load("key3")
	m.Dump()

	m.Load("key1")
	m.Dump()
}
