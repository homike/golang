package memorymodel

var c = make(chan int, 10)
var a string

func f() {
	a = "hello, world"
	c <- 0
}

func Run() {
	go f()
	<-c
	print(a)
}
