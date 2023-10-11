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
	cDB := subDB.Updates()
	cWeb := subWeb.Updates()
	processedPeople := ProcessPeople(done, 18, cWeb, cDB)

	// third stage (data consumer)
	err := WriteToFile(done, processedPeople)
	if err != nil {
		// close all routines
		close(done)
	}

	time.Sleep(time.Second * 4)
	fmt.Print("\nStoping all go routines\n\n")
	close(done)
	time.Sleep(time.Second * 4)
	fmt.Printf("active go routines: %d\n", runtime.NumGoroutine())
}
