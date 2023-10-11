# My cheat sheet for go concurrent programming 

## THE PIPELINE 
Pipeline has three stages
1. First one, getting data from *inbound* channels (producer stage) (this is data_source.go)
2. Second one, perform some processing on that data (this is process_data.go)
3. Third one, sending data to *outbound* channels (consumer stage) (this is write_to_file.go)

##  Really slow connection and responsiveness of system 

```
func (s *subWeb) loopNotResponsive() {
	for {
		select {
		case s.perChan <- ReceiveDataFromWeb():
		case <-s.done: // a receive operation on a closed channel can always proceed immediately, yielding the element type’s zero value.
			fmt.Printf("SubscriptionWeb stop\n")
			close(s.perChan)
			return
		}
	}
}

func (s *subWeb) loop() {
	receiveDone := make(chan Person)
	receive := func() {
		// we have got to wait for receive data from web
		// this hurts responsiveness of system, so the solution is to
		// move ReceiveDataFromWeb() to its own go routine
		receiveDone <- ReceiveDataFromWeb()
	}
	go receive()

	for {
		select {
		case s.perChan <- <-receiveDone:
			go receive()
		case <-s.done: // a receive operation on a closed channel can always proceed immediately, yielding the element type’s zero value.
			fmt.Printf("SubscriptionWeb stop\n")
			close(s.perChan)
			return
		}
	}
}
```

## Combine multiple chans into one 

```
// merges multiple chans into one
// out chan after closing rapidly returns value, ok
// you have to check for ok value if it is true
func Merge[T any](done <-chan struct{}, cs ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	output := func(c <-chan T) {
		defer wg.Done()
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
```

