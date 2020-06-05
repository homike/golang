package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	_nlog *NLog
	_tlog *TLog
)

var (
	_std_fatal = log.New(os.Stderr, "\033[0;33mFATAL :\033[0m ", log.LstdFlags|log.Lshortfile)
	_std_error = log.New(os.Stderr, "\033[0;31mERROR :\033[0m ", log.LstdFlags|log.Lshortfile)
	_std_warn  = log.New(os.Stdout, "\033[0;35mWARN :\033[0m ", log.LstdFlags|log.Lshortfile)
	_std_notic = log.New(os.Stdout, "\033[0;35mNOTIC :\033[0m ", log.LstdFlags|log.Lshortfile)
	_std_info  = log.New(os.Stdout, "INFO : ", log.LstdFlags|log.Lshortfile)
	_std_debug = log.New(os.Stdout, "DEBUG : ", log.LstdFlags|log.Lshortfile)
	_std_raw   = log.New(os.Stdout, "", 0)
)

// Logger不经初始化可直接使用.
func Error(format string, v ...interface{}) {
	if _nlog != nil {
		_nlog.Error(format, v...)
	} else {
		_std_error.Output(2, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

func Fatal(format string, v ...interface{}) {
	if _nlog != nil {
		_nlog.Fatal(format, v...)
	} else {
		_std_fatal.Output(2, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
	os.Exit(1)
}

func Warn(format string, v ...interface{}) {
	if _nlog != nil {
		_nlog.Warn(format, v...)
	} else {
		_std_warn.Output(2, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

func Notic(format string, v ...interface{}) {
	if _nlog != nil {
		_nlog.Notic(format, v...)
	} else {
		_std_notic.Output(2, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

func Info(format string, v ...interface{}) {
	if _nlog != nil {
		_nlog.Info(format, v...)
	} else {
		_std_info.Output(2, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

func Debug(format string, v ...interface{}) {
	if _nlog != nil {
		_nlog.Debug(format, v...)
	} else {
		_std_debug.Output(2, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

func Raw(format string, v ...interface{}) {
	if _nlog != nil {
		_nlog.Raw(format, v...)
	} else {
		_std_raw.Output(2, fmt.Sprintln(fmt.Sprintf(format, v...)))
	}
}

func AsyncSendTlog(logheader string, content string, conindex uint8) string {
	if _tlog != nil {
		return _tlog.AsyncSendTlog(logheader, content, conindex)
	}
	return ""
}
