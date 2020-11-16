package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9100")
	if err != nil {
		return
	}

	inputReader := bufio.NewReader(os.Stdin)
	for {
		//Write
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("read input error")
			return
		}
		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Println("write string error1, ", err)
			return
		}

		// Read
		buf := make([]byte, 10)
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println("read string error, ", err)
			//return
		}
		fmt.Print("[server]", string(buf))
	}
}
