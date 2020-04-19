package gologs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func BenchmarkLogger_Zeo_File(b *testing.B) {
	b.StopTimer()

	f, err := os.OpenFile("zeo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer f.Close()

	logger := zerolog.New(f).With().Timestamp().Logger()
	logger = logger.With().Caller().Logger() // 打印代码行
	//logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // 带颜色输出

	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			logger.Info().Msg("aaa")
		}
	})
}

func BenchmarkLogger_Zeo_Discard(b *testing.B) {
	b.StopTimer()

	logger := zerolog.New(ioutil.Discard).With().Timestamp().Logger()
	logger = logger.With().Caller().Logger() // 打印代码行
	//logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // 带颜色输出

	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			logger.Info().Msg("aaa")
		}
	})
}

func BenchmarkLogger_Zeo_File_Nocaller(b *testing.B) {
	b.StopTimer()

	f, err := os.OpenFile("zeo_nocaller.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer f.Close()

	logger := zerolog.New(f).With().Timestamp().Logger()
	//logger = logger.With().Caller().Logger() // 打印代码行
	//logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // 带颜色输出

	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			logger.Info().Msg("aaa")
		}
	})
}
