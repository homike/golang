package gologs

import (
	"fmt"
	"testing"

	"gitee.com/gricks/logrus"
)

func BenchmarkLoggerFile_LogrusAstone(b *testing.B) {
	b.StopTimer()

	logger := logrus.New(logrus.WithFile("logrusastone.log"))
	defer logger.Close()
	entry := logger.GetEntry().SetPrefix("prefix")

	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			entry.Info("aaa")
		}
	})
}

func BenchmarkLoggerFile_LogrusAstone_buff4096(b *testing.B) {
	logger := logrus.New(
		logrus.WithFile("logrus_bufio_4096"),
		logrus.WithColor(true),
		logrus.WithCaller(true),
		logrus.WithRotater(new(logrus.DailyRotater)),
		logrus.WithBufferSize(4096),
	)
	defer logger.Close()

	entry := logger.GetEntry().SetPrefix(fmt.Sprintf("UID:%d|GID:%d|Command:%s", 10000, 1, "GET_USER_INFO"))

	b.StopTimer()
	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			entry.Info("aaa")
		}
	})
}
