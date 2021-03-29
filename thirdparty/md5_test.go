package thirdparty

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"
)

func TestMd5(t *testing.T) {
	h := md5.New()
	tp := time.Now().Unix()
	io.WriteString(h, "xxxx123"+strconv.Itoa(int(tp)))
	fmt.Printf("md5: %x", h.Sum(nil))
}
