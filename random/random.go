package random

import "math/rand"

func RandNum(base int, n int) int {
	if n <= 0 {
		return base
	}
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//return base + r.Intn(n)
	return base + rand.Intn(n)
}

func GetRandomWeightID(ids, weights []int) int {
	if len(ids) != len(weights) {
		return 0
	}

	total := 0
	randRange := make([]int, len(ids))
	for i := 0; i < len(ids); i++ {
		total += weights[i]
		randRange[i] = total
	}

	randVal := RandNum(0, total)
	for i := 0; i < len(randRange); i++ {
		if randRange[i] > randVal {
			return ids[i]
		}
	}
	return ids[0]
}
