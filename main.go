package main

import (
	"fmt"
	"sync"
)

func say(msg string) {
	fmt.Println(msg)
}

func main() {
	say("hello start")

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			say(fmt.Sprintf("hello %d", i))
		}()
	}

	wg.Wait()

	say("hello end")
}
