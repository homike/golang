package errortest

import (
	"log"
	"net"
	"syscall"
)

func isCaredNetError(err error) bool {
	netErr, ok := err.(net.Error)
	if !ok {
		return false
	}

	// net error
	if netErr.Timeout() {
		log.Print("time out")
		return true
	}
	if netErr.Temporary() {
		log.Print("time out")
		return true
	}

	// op error
	opErr, ok := netErr.(*net.OpError)
	if !ok {
		return false
	}

	switch t := opErr.Err.(type) {
	case *net.DNSError:
		log.Print("dns error")
		return true
	case *net.SyscallError:
		log.Print("sys call error")
		if errno, ok := t.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.ETIMEDOUT:
				log.Print("time out")
				return true
			case syscall.ECONNREFUSED:
				log.Print("connection refused")
				return true
			}
		}
		return true
	}
}
