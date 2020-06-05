package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ByteSize float64

const (
	_           = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

const (
	_ = iota
	LOG_LEVEL_FATAL
	LOG_LEVEL_ERROR
	LOG_LEVEL_WARN
	LOG_LEVEL_NOTIC
	LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG
)

const (
	CHECK_INTERVAL = 2 * time.Minute
	CHECK_EXPIRED  = 2 * time.Hour
	//CHECK_EXPIRED = 2 * time.Second
)

type NLog struct {
	errCount int32
	rotate   Rotate
	level    chan int

	rwm                                  sync.RWMutex
	file                                 *os.File
	debug, info, notic, warn, err, fatal *log.Logger
	raw                                  *log.Logger
}

type Rotate struct {
	Size              ByteSize
	Expired, Interval time.Duration
}

func (logger *NLog) SetLevel(level int) {
	logger.level <- level
}

// 设置Logger的级别
func (logger *NLog) setLevel(f *os.File, level int) {
	switch {
	case level >= LOG_LEVEL_DEBUG:
		logger.debug = log.New(f, "\033[0;36mDEBUG:\033[0m ", log.LstdFlags|log.Lshortfile)
		fallthrough
	case level >= LOG_LEVEL_INFO:
		logger.info = log.New(f, "INFO : ", log.LstdFlags|log.Lshortfile)
		fallthrough
	case level >= LOG_LEVEL_NOTIC:
		logger.notic = log.New(f, "\033[0;32mNOTIC:\033[0m ", log.LstdFlags|log.Lshortfile)
		fallthrough
	case level >= LOG_LEVEL_WARN:
		logger.warn = log.New(f, "\033[0;35mWARN :\033[0m ", log.LstdFlags|log.Lshortfile)
		fallthrough
	case level >= LOG_LEVEL_ERROR:
		logger.err = log.New(f, "\033[0;31mERROR:\033[0m ", log.LstdFlags|log.Lshortfile)
		fallthrough
	case level >= LOG_LEVEL_FATAL:
		logger.fatal = log.New(f, "\033[0;33mFATAL:\033[0m ", log.LstdFlags|log.Lshortfile)
	}
	switch {
	case level < LOG_LEVEL_FATAL:
		logger.fatal = nil
		fallthrough
	case level < LOG_LEVEL_ERROR:
		logger.err = nil
		fallthrough
	case level < LOG_LEVEL_WARN:
		logger.warn = nil
		fallthrough
	case level < LOG_LEVEL_NOTIC:
		logger.notic = nil
		fallthrough
	case level < LOG_LEVEL_INFO:
		logger.info = nil
		fallthrough
	case level < LOG_LEVEL_DEBUG:
		logger.debug = nil
	}

	logger.raw = log.New(f, "", 0)
}

// 获取Logger文件大小
func (logger *NLog) getFileSize() ByteSize {
	logger.rwm.RLock()
	defer logger.rwm.RUnlock()
	fi, err := logger.file.Stat()
	if err != nil {
		logger.Warn("get log file size failed, no trunc %s", err.Error())
		return 0.0
	}
	return ByteSize(fi.Size())
}

// 超过Interval的Logger文件重命名
func (logger *NLog) trunc(fp, ext string, level int) {
	logger.rwm.Lock()
	defer logger.rwm.Unlock()
	err := logger.file.Close()
	if err != nil {
		_std_warn.Println("fail to close log file", err.Error())
		return
	}
	err = os.Rename(fp, fp+ext)
	if err != nil {
		_std_warn.Println("fail to rename log file, no trunc", err.Error())
	}
	file, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Error("create log file failed %s", err.Error())
		return
	}
	logger.setLevel(file, level)
	logger.file = file
}

// 生成Logger日期后缀
func suffix(t time.Time) string {
	year, month, day := t.Date()
	return "-" + fmt.Sprintf("%04d%02d%02d%02d", year, month, day, t.Hour())
}

// 截断获得下一次指定时间段的时间
func toNextBound(duration time.Duration) time.Duration {
	return time.Now().Truncate(duration).Add(duration).Sub(time.Now())
}

// Logger处理函数
func (logger *NLog) loop() error {
	interval := time.After(toNextBound(logger.rotate.Interval))
	expired := time.After(CHECK_EXPIRED)

	var sizeExt, level int = 1, 4

	fn := filepath.Base(logger.file.Name())

	fp, err := filepath.Abs(logger.file.Name())

	if err != nil {
		_std_fatal.Fatalln("get log filepath failed %s", err.Error())
	}

	for {
		var size <-chan time.Time
		if toNextBound(logger.rotate.Interval) != CHECK_INTERVAL {
			size = time.After(CHECK_INTERVAL)
		}
		select {
		case level = <-logger.level:
			// 变更Logger等级
			logger.rwm.Lock()
			logger.setLevel(logger.file, level)
			logger.rwm.Unlock()
			logger.Notic("log level change to %d", level)
		case t := <-interval:
			// 自定义生成新的Logger文件
			interval = time.After(logger.rotate.Interval)
			logger.trunc(fp, suffix(t), level)
			sizeExt = 1
			logger.Notic("log truncated by time interval")
		case <-expired:
			// 删除过期的Logger文件
			expired = time.After(CHECK_EXPIRED)
			err := filepath.Walk(filepath.Dir(fp),
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return nil
					}
					isLog := strings.Contains(info.Name(), fn)

					//log.Println("strings.Contains(", info.Name(), " log') isLog = ", isLog)
					if time.Since(info.ModTime()) > logger.rotate.Expired && isLog && info.IsDir() == false {
						if err := os.Remove(path); err != nil {
							return err
						}
						logger.Notic("remove expired log files %s", filepath.Base(path))
					}
					return nil
				})
			if err != nil {
				logger.Warn("remove expired logs failed %s", err.Error())
			}
		case t := <-size:
			// 文件大小超过上限
			if logger.getFileSize() < logger.rotate.Size {
				break
			}
			curTm := t.Add(logger.rotate.Interval)
			logger.trunc(fp, suffix(curTm)+"."+strconv.Itoa(sizeExt), level)
			sizeExt++
			logger.Notic("log over size, truncated")
		}
	}
}

// Debug log debug message with cyan color.
func (logger *NLog) Debug(format string, v ...interface{}) {
	logger.rwm.RLock()
	defer logger.rwm.RUnlock()

	if logger.debug != nil {
		logger.debug.Output(3, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

// Info log normal message.
func (logger *NLog) Info(format string, v ...interface{}) {
	logger.rwm.RLock()
	defer logger.rwm.RUnlock()

	if logger.info != nil {
		logger.info.Output(3, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

// Notice log notice message with blue color.
func (logger *NLog) Notic(format string, v ...interface{}) {
	logger.rwm.RLock()
	defer logger.rwm.RUnlock()

	if logger.notic != nil {
		logger.notic.Output(3, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

// Error log error message with red color.
func (logger *NLog) Error(format string, v ...interface{}) {
	atomic.AddInt32(&logger.errCount, 1)
	_std_error.Output(3, fmt.Sprintln(fmt.Sprintf(format, v...)))

	logger.rwm.RLock()
	defer logger.rwm.RUnlock()

	if logger.err != nil {
		logger.err.Output(3, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

func (logger *NLog) ErrCount() int32 {
	ec := atomic.LoadInt32(&logger.errCount)
	if ec < 0 {
		logger.Warn("error count overflow")
		return -1
	}
	return ec
}

func (logger *NLog) Fatal(format string, v ...interface{}) {
	_std_fatal.Output(3, fmt.Sprintln(fmt.Sprintf(format, v...)))

	logger.rwm.RLock()
	defer logger.rwm.RUnlock()

	if logger.fatal != nil {
		logger.fatal.Output(3, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
	os.Exit(1)
}

func (logger *NLog) Warn(format string, v ...interface{}) {
	_std_warn.Output(3, fmt.Sprintln(fmt.Sprintf(format, v...)))

	logger.rwm.RLock()
	defer logger.rwm.RUnlock()

	if logger.warn != nil {
		logger.warn.Output(3, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

// Raw log raw message.
func (logger *NLog) Raw(format string, v ...interface{}) {
	logger.rwm.RLock()
	defer logger.rwm.RUnlock()

	logger.raw.Output(3, fmt.Sprintln(fmt.Sprintf(format, v...)))
}
