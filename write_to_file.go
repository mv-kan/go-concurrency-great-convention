package main

import (
	"fmt"
	"os"
)

func WriteToFile(done <-chan struct{}, processedChan <-chan Person) error {
	f, err := os.OpenFile("result.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	go func() {
		defer f.Close()
		for {
			select {
			case p, ok := <-processedChan:
				if !ok {
					fmt.Println("WriteToFile done")
					return
				}
				line := fmt.Sprintf("%s %d\n", p.Name, p.Age)
				_, err := f.WriteString(line)
				if err != nil {
					fmt.Printf("WriteToFile err=%v", err)
				}
			case <-done:
				fmt.Println("WriteToFile done")
				return
			}
		}
	}()

	return nil
}
