package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP("127.0.0.1"), 9100, ""})
	if err != nil {
		return
	}

	conn, err := listen.Accept()
	if err != nil {
		return
	}

	Run(conn)
}

func Run(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	for {
		buf := make([]byte, 10)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error")
			return
		}

		_, err = writer.Write(buf)
		if err != nil {
			fmt.Println("write error")
			return
		}
		writer.Flush()
	}
}
