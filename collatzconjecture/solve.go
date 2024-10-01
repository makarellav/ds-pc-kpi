package collatzconjecture

import "sync/atomic"

// type result struct {
// 	steps      int
// 	hailstones []int
// }

type total struct {
	steps int64 // must be int64 to use with atomic package
}

func (t *total) add(steps int) {
	atomic.AddInt64(&t.steps, int64(steps))
}

func (t *total) get() int64 {
	return atomic.LoadInt64(&t.steps)
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
		t.add(1)
	}

	return int(t.get())
}
