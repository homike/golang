package basetest

import (
	"log"
	"runtime"
)

func test() {
	test2()
}

func test2() {
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
