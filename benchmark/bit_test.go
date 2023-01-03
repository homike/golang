package benchmark

import "testing"

func IsType1(i int32) bool {
	return i == 1 ||
		i == 2 ||
		i == 3 ||
		i == 4 ||
		i == 5 ||
		i == 6 ||
		i == 7 ||
		i == 8 ||
		i == 9 ||
		i == 10
}

var type2 = 100

func Type2() uint64 {
	alltype := uint64(0)
	alltype |= 1 << 1
	alltype |= 1 << 2
	alltype |= 1 << 3
	alltype |= 1 << 4
	alltype |= 1 << 5
	alltype |= 1 << 6
	alltype |= 1 << 7
	alltype |= 1 << 8
	alltype |= 1 << 9
	alltype |= 1 << 10
	return alltype
}

func IsType2(i int32) bool {
	return (type2>>i)&1 == 1
}

func Benchmark_IsType1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsType1(int32(i))
	}
}

func Benchmark_IsType2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsType2(int32(i))
	}
}
