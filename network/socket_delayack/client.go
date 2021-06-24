package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func PackStream(data []byte) []byte {
	buff := make([]byte, 2)

	dataLen := len(data)
	binary.BigEndian.PutUint16(buff, uint16(dataLen))

	buff = append(buff, data...)
	return buff
}

func unPackStream(r io.Reader) ([]byte, error) {
	bufLen := make([]byte, 2, 2)
	io.ReadFull(r, bufLen)

	dataLen := binary.BigEndian.Uint16(bufLen)

	if dataLen > 0 {
		data := make([]byte, dataLen)
		if _, err := io.ReadFull(r, data); err != nil {
			fmt.Println(err)
			return nil, err
		}
		return data, nil
	}
	return nil, fmt.Errorf("nil data")
}

func main() {
	SendAuto()
}

func SendAuto() {
	conn, err := net.Dial("tcp", "10.10.6.203:9100")
	if err != nil {
		return
	}
	// set nodely
	conn.(*net.TCPConn).SetNoDelay(false)
	// set quick ack
	/*
		f, err := conn.File()
		if err != nil {
			err = unix.SetsockoptInt(int(f.Fd()), unix.IPPROTO_TCP, unix.TCP_QUICKACK, 1)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	*/

	fillData := func(l int) []byte {
		data := make([]byte, l, l)
		for i := 0; i < l; i++ {
			val := i % 100
			data[i] = byte(val)
		}
		return data
	}
	data1K := fillData(100)

	writer := bufio.NewWriter(conn)
	for i := 0; ; i++ {
		sendData := PackStream(data1K)
		if i < 10 {
			//Write
			_, err = writer.Write(sendData)
			if err != nil {
				fmt.Println("write string error1, ", err)
				return
			}
			writer.Flush()
		} else {
			sub1, sub2 := sendData[:10], sendData[10:]
			//Write
			_, err = writer.Write(sub1)
			if err != nil {
				fmt.Println("write string error1, ", err)
				return
			}
			writer.Flush()

			//Write
			_, err = writer.Write(sub2)
			if err != nil {
				fmt.Println("write string error1, ", err)
				return
			}
			writer.Flush()
		}

		// Read
		_, err := unPackStream(conn)
		if err != nil {
			fmt.Println("read string error, ", err)
			//return
		}
		//fmt.Printf("[server] %v, \n", string(buf))
	}
}
