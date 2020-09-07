package mathtest

import "fmt"

const kItemTypeOffset = 100000

func GetItemKey(itemType, itemID int32) int32 {
	return itemType*kItemTypeOffset + itemID
}

func ParseItemKey(itemKey int32) (int32, int32) {
	return itemKey / kItemTypeOffset, itemKey % kItemTypeOffset
}

func RunMathTest() {
	key := GetItemKey(9991, 99991)
	fmt.Println("key ", key)

	id, num := ParseItemKey(key)
	fmt.Println("id ", id, "num ", num)
	//stack.PrintGoroutineMemConsume()
	//channeltest.RunUnbounded()
	//go timertest.RunTimerAfter()

	//test2()

	//http.ListenAndServe("0.0.0.0:8899", nil)

}
