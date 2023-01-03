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
	io.WriteString(h, "oPrGGM8pbZDg"+strconv.Itoa(int(tp)))
	fmt.Printf("timestamp:%v, md5: %x", tp, h.Sum(nil))
}
