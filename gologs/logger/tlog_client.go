package logger

import (
	"math/rand"
	"net"
	"sync"
)

type TlogCommon struct {
	svrID      int32
	zoneID     int32
	tlogSwitch int32
	appID      string
}

type TlogClient struct {
	common *TlogCommon
	_conn  []net.Conn

	lock sync.RWMutex
}

func NewTlogClient(common *TlogCommon, tlogaddrs []string) *TlogClient {

	tlc := &TlogClient{}
	tlc.common = common

	err := tlc.DialTlogSvr(tlogaddrs)
	if err != nil {
		Error("new tlog client (%v, %v) err (%v).", common, tlogaddrs, err)
		return nil
	}

	return tlc
}

// connect to tlogd server
func (tlc *TlogClient) DialTlogSvr(tlogaddrs []string) error {
	tlc.lock.Lock()
	defer tlc.lock.Unlock()

	Notic("connect tlog server %v.", tlogaddrs)
	tlc._conn = make([]net.Conn, len(tlogaddrs))

	for i := 0; i < len(tlogaddrs); i++ {
		addr, err := net.ResolveUDPAddr("udp", tlogaddrs[i])
		if err != nil {
			Error("net.ResolveUDPAddr (%s) err (%v).", tlogaddrs[i], err)
			return err
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			Error("net.DialUDP (%s) err (%v).", tlogaddrs[i], err)
			return err
		}
		//conn.SetDeadline(time.Now().Add(TLOG_SEND_TIMEOUT * time.Second))
		tlc._conn[i] = conn
		Notic("tlog server (%v) connected.", tlogaddrs[i])
	}

	return nil
}

func (tlc *TlogClient) MapIDToConIdx() uint8 {
	tlc.lock.RLock()
	defer tlc.lock.RUnlock()

	conIndex := rand.Intn(len(tlc._conn))

	return uint8(conIndex)
}

// send msg to tlogd server
func (tlc *TlogClient) Send(data []byte, conindex uint8) bool {
	tlc.lock.RLock()
	defer tlc.lock.RUnlock()

	//Error("send packet to tlog server(%v, %s).", conindex, string(data))

	if int(conindex) >= len(tlc._conn) {
		Error("send (%d, %d) err (conindex out of range).", conindex, len(tlc._conn))
		return false
	}
	// send the packet
	_, err := tlc._conn[conindex].Write(data)
	if err != nil {
		Error("send packet to tlog server(%v, %s) err (%v).", conindex, string(data), err)
		return false
	}
	return true
}

// send msg to tlogd server
func (tlc *TlogClient) AsyncSend(data []byte, conindex uint8) {
	go tlc.Send(data, conindex)
}
