package collatzconjecture

import "sync"

// type result struct {
// 	steps      int
// 	hailstones []int
// }

type total struct {
	mu    sync.Mutex
	steps int
}

func Solve(n int) int {
	// var steps int
	// hailstones := []int{n}
	var t total

	for n != 1 {
		if n%2 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}

		// steps++
		// hailstones = append(hailstones, n)
		t.steps++
	}

	return t.steps
}
