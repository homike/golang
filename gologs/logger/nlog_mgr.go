package logger

import (
	"os"
	"time"
)

// 初始化Logger
func InitNLog(nlogpath string, loggerlevel int) error {

	logger := NewNLog(nlogpath, loggerlevel, Rotate{Size: 500 * MB, Expired: time.Hour * 24 * 7, Interval: time.Hour})
	_nlog = logger
	return nil
}

// 初始化并返回NLogger指针
func NewNLog(fp string, level int, rotate Rotate) *NLog {
	file, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		_std_fatal.Fatalln(err.Error())
	}
	logger := &NLog{}
	logger.rotate = rotate
	logger.file = file
	logger.level = make(chan int)
	go logger.loop()
	logger.SetLevel(level)
	return logger
}

// 设置Logger级别
func SetNLogLevel(level int) {
	_nlog.SetLevel(level)
}

// 返回NLogger指针
func GetNLog() *NLog {
	return _nlog
}
