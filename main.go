package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

func paperSign() {
	h := md5.New()
	tp := time.Now().Unix()
	//tp := 1638349105
	io.WriteString(h, "oPrGGM8pbZDg"+strconv.Itoa(int(tp)))
	fmt.Printf("timestamp:%v, md5: %x", tp, h.Sum(nil))
}

func main() {
}
