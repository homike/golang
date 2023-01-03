package main

import "fmt"

func f1() error {
	return nil
}

func f2() error {
	return nil
}

func main() {
	{
		fmt.Println("test")
	}

	//err := f1()
	//{
	err := f2()
	if err != nil {
		//fmt.Println("number: ", n)
		return
	}
	//}
	//fmt.Println("error: ", err)
}
