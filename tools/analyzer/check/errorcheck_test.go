package check

import (
	"fmt"
	"testing"
)

func Test_ErrorcCeck(t *testing.T) {
	err := ErrorCheck("./test")
	fmt.Println("UnCheckError: ", err)
}
