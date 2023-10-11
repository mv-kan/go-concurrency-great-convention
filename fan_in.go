package main

import (
	"fmt"
	"sync"
)

// merges multiple chans into one
// out chan after closing rapidly returns value, ok
// you have to check for ok value if it is true
func Merge[T any](done <-chan struct{}, cs ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	output := func(c <-chan T) {
		defer wg.Done()
		defer func() { fmt.Println("Merge routine done") }()

		for v := range c { // finishes when c (chan) is closed
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}

	for _, c := range cs {
		wg.Add(1)
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
