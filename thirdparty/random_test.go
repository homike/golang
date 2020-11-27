package thirdparty

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	arrRewards := []int32{111, 222, 333, 444, 555, 666, 777}
	l := len(arrRewards)
	ia := rand.New(rand.NewSource(time.Now().Unix())).Perm(l)[0:1]

	var rs []int32

	for i := 0; i < len(ia); i++ {
		rs = append(rs, arrRewards[ia[i]])
	}

	fmt.Println("result:", rs)
}
