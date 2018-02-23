package main

import (
	"GoTest/astar"
	_ "GoTest/basetest"
	_ "GoTest/channeltest"
)

// 测试提交
// Main

func main() {
	astar.RunAstar()
	// basetest.Run2()

	// reflect.RunReflect1()
	// bittest.RunBit()
	// interfacetest.RunInterface()
	// thirdparty.RunRandom()
}
