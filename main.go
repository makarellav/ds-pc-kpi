package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/makarellav/ds-pc-kpi/collatzconjecture"
)

func main() {
	implementation := flag.String("implementation", "naive",
		"Collatz conjecture implementation to be used. Options: naive, semaphore, batch, workerpool. (default: naive)")

	maxWorkers := flag.Int("maxWorkers", runtime.NumCPU(),
		"number of workers to be used (only for semaphore, batch, workerpool implementations). (default: runtime.NumCPU())")

	n := flag.Int("n", 10000000,
		"positive integer for Collatz conjecture. (default: 10000000)")

	batchSize := flag.Int("batchSize", 1000,
		"batch size for batch implementation. (default: 1000)")

	flag.Parse()

	start := time.Now()

	var result int
	switch *implementation {
	case "naive":
		result = collatzconjecture.NaiveImplementation(*n)
	case "semaphore":
		result = collatzconjecture.SemaphoreImplementation(*n, *maxWorkers)
	case "batch":
		result = collatzconjecture.SemaphoreImplementationBatch(*n, *maxWorkers, *batchSize)
	case "workerpool":
		result = collatzconjecture.WorkerPoolImplementation(*n, *maxWorkers)
	default:
		fmt.Println("Invalid implementation specified. Please choose from: naive, semaphore, batch, workerpool.")
		return
	}

	duration := time.Since(start)

	fmt.Printf("Avg. Steps: %d\n", result)
	fmt.Printf("Execution Time for %q implementation: %s\n", *implementation, duration)
}
