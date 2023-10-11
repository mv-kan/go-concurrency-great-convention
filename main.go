package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Starting...")
	done := make(chan struct{})
	// first stage (data producer)
	subWeb := SubscribeWeb()
	subDB := SubscribeDB()

	// close subs when done is closed
	go func() {
		<-done
		subDB.Close()
		subWeb.Close()
	}()

	// second stage
	processedPeople := ProcessPeople(done, 18, subWeb, subDB)

	// third stage (data consumer)
	err := WriteToFile(done, processedPeople)
	if err != nil {
		// close all routines
		close(done)
	}
	// fmt.Print("insert any value here to stop program: ")
	// input := bufio.NewScanner(os.Stdin)
	// input.Scan()
	time.Sleep(time.Second * 4)
	fmt.Println("Stoping all go routines")
	done <- struct{}{}
	time.Sleep(time.Second * 4)
	fmt.Printf("active go routines: %d\n", runtime.NumGoroutine())
	time.Sleep(time.Second * 100)

}
