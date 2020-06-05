package logger

import (
	"strings"
	"sync"
)

const (
	TLOG_SEND_TIMEOUT = 10 // udp sending timeout
)

type TLog struct {
	tlogClient *TlogClient

	lock sync.RWMutex
}

func (t *TLog) Reload(tlogsvrs string, tlogswitch int32, zoneid int32, svrid int32, appid string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// init the common params
	tlogCommon := &TlogCommon{}
	tlogCommon.zoneID = zoneid
	tlogCommon.svrID = svrid
	tlogCommon.appID = appid
	tlogCommon.tlogSwitch = tlogswitch

	if tlogsvrs != "" {
		tlogAddrs := strings.Split(tlogsvrs, ",")
		t.tlogClient = NewTlogClient(tlogCommon, tlogAddrs)
	}
	if t.tlogClient == nil {
		panic("t.tlogClient is nil")
	}
}

func (t *TLog) AsyncSendTlog(logheader string, content string, conindex uint8) string {
	t.lock.RLock()
	defer t.lock.RUnlock()

	var loginfo string = logheader + "|" + content

	if t.tlogClient.common.tlogSwitch == 1 {
		t.tlogClient.AsyncSend([]byte(loginfo+"\n"), conindex)
	}

	return loginfo
}

func (t *TLog) MapIDToConIdx() uint8 {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.tlogClient.MapIDToConIdx()
}
