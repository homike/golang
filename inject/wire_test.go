package inject

import (
	"testing"
)

func TestWire(t *testing.T) {
	i := InitializeAllInstance()
	i.Bar.Test()

	t.Log("hello world")
}
