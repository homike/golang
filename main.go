package main

import (
	"GoTest/basetest"
	_ "GoTest/channeltest"
)

// 测试提交
// Main

func main() {
	//reflect.RunReflect1()
	basetest.Run3()

	//bittest.RunBit()
	//interfacetest.RunInterface()
	//thirdparty.RunRandom()
}
