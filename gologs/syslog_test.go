package gologs

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func BenchmarkLogger_Sys_File(b *testing.B) {
	b.StopTimer()

	f, err := os.OpenFile("sys.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags|log.Lshortfile)

	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			logger.Output(2, "aaa")
		}
	})
}

func BenchmarkLogger_Sys_Discard(b *testing.B) {
	b.StopTimer()

	logger := log.New(ioutil.Discard, "", log.LstdFlags|log.Lshortfile)

	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			logger.Output(2, "aaa")
		}
	})
}

func BenchmarkLogger_Sys_File_Nocaller(b *testing.B) {
	b.StopTimer()

	f, err := os.OpenFile("sys_nocaller.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)

	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			logger.Output(2, "aaa")
		}
	})
}
