package logger

// 初始化Logger
func InitTLog(tlogsvrs string, tlogswitch int32, zoneid int32, gamesvrid int32, appid string) error {
	logger := NewTLog(tlogsvrs, tlogswitch, zoneid, gamesvrid, appid)
	_tlog = logger
	return nil
}

// 初始化并返回NLogger指针
func NewTLog(tlogsvrs string, tlogswitch int32, zoneid int32, gamesvrid int32, appid string) *TLog {
	logger := &TLog{}
	logger.Reload(tlogsvrs, tlogswitch, zoneid, gamesvrid, appid)
	return logger
}

// 返回NLogger指针
func GetTLog() *TLog {
	return _tlog
}
