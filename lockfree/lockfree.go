package lockfree

import (
	"os"
	"runtime/trace"
)

type Test struct {
	A int
	B map[int]int
}

func LockFree() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	v := &Test{
		A: 1,
		B: make(map[int]int),
	}
	v.B[1] = 1

	for i := 0; i < 10; i++ {
		go func(arg int) {
			for {
				v1 := &Test{
					A: arg,
					B: make(map[int]int),
				}
				v1.B[arg] = arg

				v = v1
				//v.B[arg] = arg
				//fmt.Println(arg)
			}
		}(i)
	}

	for {
	}
}
