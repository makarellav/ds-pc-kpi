package collatzconjecture

import "sync"

func solve(wg *sync.WaitGroup, jobCh <-chan int, t *total) {
	defer wg.Done()

	for n := range jobCh {
		result := Solve(n)

		t.add(result)
	}
}

// ||||

func WorkerPoolImplementation(n int, workers int) int {
	jobCh := make(chan int)
	// outputCh := make(chan result)

	var wg sync.WaitGroup
	var t total

	for range workers {
		wg.Add(1)
		go solve(&wg, jobCh, &t)
	}

	go func() {
		for i := 1; i <= n; i++ {
			jobCh <- i
		}

		close(jobCh)
	}()

	wg.Wait()

	return int(t.get()) / n
}
