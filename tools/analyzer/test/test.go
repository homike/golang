package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("vim-go")

	go func() {
		defer recover()
		fmt.Println("aaa")
	}()

	n, err := strconv.Atoi("1")
	fmt.Println(n, err)
}
