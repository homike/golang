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
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP("10.10.6.203"), 9100, ""})
	if err != nil {
		return
	}

	conn, err := listen.Accept()
	if err != nil {
		return
	}
	conn.(*net.TCPConn).SetNoDelay(false)

	Run(conn)
}

func Run(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	for {
		data, err := unPackStream(conn)
		if err != nil {
			fmt.Printf("[client] %v, %v \n", err, data)
		}

		//time.Sleep(1 * time.Millisecond)
		_, err = writer.Write(PackStream(data))
		if err != nil {
			fmt.Println("write error")
			return
		}
		writer.Flush()
	}
}
