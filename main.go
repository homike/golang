package main

import (
	"GoTest/httptest"
	"fmt"
)

func main() {
	httptest.RunRedis()
	fmt.Printf("hello!!!")
}
