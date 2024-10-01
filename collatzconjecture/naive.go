package collatzconjecture

func NaiveImplementation(n int) int {
	var t total

	for i := 1; i <= n; i++ {
		result := Solve(i)

		t.add(result)
	}

	return int(t.get()) / n
}
