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

			t.add(result)
		}(i)
	}

	wg.Wait()

	// var total int

	// for range n {
	// 	result := <-outputCh
	// 	total += result.steps
	// }

	return int(t.get()) / n
}

func SemaphoreImplementationBatch(n int, maxWorkers int, batchSize int) int {
	sem := semaphore{
		c: make(chan struct{}, maxWorkers),
	}

	// outputCh := make(chan int, n)

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

			t.add(batchTotal)
		}(i)
	}

	wg.Wait()

	return int(t.get()) / n
}
