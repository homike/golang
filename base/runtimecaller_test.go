package base

import (
	"log"
	"runtime"
	"testing"
)

func Test_Caller(t *testing.T) {
	test()
}

func test() {
	pc, file, line, _ := runtime.Caller(3)
	f := runtime.FuncForPC(pc)
	log.Println("pc3 ", pc, "file ", file, "line", line, "name", f.Name())

	pc, file, line, _ = runtime.Caller(2)
	f = runtime.FuncForPC(pc)
	log.Println("pc2 ", pc, "file ", file, "line", line, "name", f.Name())

	pc, file, line, _ = runtime.Caller(1)
	f = runtime.FuncForPC(pc)
	log.Println("pc1 ", pc, "file ", file, "line", line, "name", f.Name())

	pc, file, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.Println("pc0 ", pc, "file ", file, "line", line, "name", f.Name())
}
