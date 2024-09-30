package collatzconjecture

import "sync"

type semaphore struct {
	c chan struct{}
}

func (s *semaphore) acquire() {
	s.c <- struct{}{}
}

func (s *semaphore) release() {
	<-s.c
}

func SemaphoreImplementation(n int, maxWorkers int) int {
	sem := semaphore{
		c: make(chan struct{}, maxWorkers),
	}

	// outputCh := make(chan result, n)

	var t total

	var wg sync.WaitGroup

	for i := 1; i <= n; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			sem.acquire()
			defer sem.release()

			result := Solve(n)

			t.mu.Lock()
			t.steps += result
			t.mu.Unlock()
		}(i)
	}

	wg.Wait()

	// var total int

	// for range n {
	// 	result := <-outputCh
	// 	total += result.steps
	// }

	return t.steps / n
}

func SemaphoreImplementationBatch(n int, maxWorkers int, batchSize int) int {
	sem := semaphore{
		c: make(chan struct{}, maxWorkers),
	}

	var t total
	var wg sync.WaitGroup

	for i := 1; i <= n; i += batchSize {
		wg.Add(1)

		go func(start int) {
			defer wg.Done()

			sem.acquire()
			defer sem.release()

			batchTotal := 0
			for j := start; j < start+batchSize && j <= n; j++ {
				steps := Solve(j)
				batchTotal += steps
			}

			t.mu.Lock()
			t.steps += batchTotal
			t.mu.Unlock()
		}(i)
	}

	wg.Wait()

	return t.steps / n
}
