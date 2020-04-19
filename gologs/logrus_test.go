package gologs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
)

func BenchmarkLogger_Logrus_File(b *testing.B) {
	b.StopTimer()

	f, err := os.OpenFile("logrus.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer f.Close()

	log2 := logrus.New()
	log2.SetLevel(logrus.DebugLevel)
	log2.SetReportCaller(true)
	log2.Formatter = &logrus.TextFormatter{
		DisableColors:  true,
		FullTimestamp:  true,
		DisableSorting: true,
	}

	log2.SetOutput(f)

	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			log2.Debugf("aaa")
		}
	})
}

func BenchmarkLogger_Logrus_Discard(b *testing.B) {
	b.StopTimer()

	log2 := logrus.New()
	log2.SetLevel(logrus.DebugLevel)
	log2.SetReportCaller(true)
	log2.Formatter = &logrus.TextFormatter{
		DisableColors:  true,
		FullTimestamp:  true,
		DisableSorting: true,
	}

	log2.SetOutput(ioutil.Discard)

	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			log2.Debugf("aaa")
		}
	})
}

func BenchmarkLogger_Logrus_File_Nocaller(b *testing.B) {
	b.StopTimer()

	f, err := os.OpenFile("logrus_nocaller.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer f.Close()

	log2 := logrus.New()
	log2.SetLevel(logrus.DebugLevel)
	log2.SetReportCaller(false)
	log2.Formatter = &logrus.TextFormatter{
		DisableColors:  true,
		FullTimestamp:  true,
		DisableSorting: true,
	}

	log2.SetOutput(f)

	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			log2.Debugf("aaa")
		}
	})
}
