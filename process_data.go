/*
SECOND STAGE
*/
package main

import (
	"fmt"
)

/*
Run processing on persons, if age if higher than minimal age then
Write the person to file that is allow list to club
*/
func ProcessPeople(done <-chan struct{}, minimalAge int, subs ...Subscription) <-chan Person {
	processedChan := make(chan Person)
	go func() {
		cs := []<-chan Person{}
		for _, sub := range subs {
			cs = append(cs, sub.Updates())
		}
		c := Merge(done, cs...)

		for {
			select {
			case person, ok := <-c:
				if !ok {
					fmt.Println("ProcessPeople done")
					close(processedChan)
					return
				}
				if person.Age >= minimalAge {
					fmt.Println(person, "goes in")
					processedChan <- person
				} else {
					fmt.Println(person, "goes out")
				}
			case <-done:
				fmt.Println("ProcessPeople done")
				close(processedChan)
				return
			}
		}
	}()
	return processedChan
}
