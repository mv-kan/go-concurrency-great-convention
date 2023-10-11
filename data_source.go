/*
https://go.dev/blog/pipelines
Pipeline has three stages
1. First one, getting data from *inbound* channels (producer stage)
2. Second one, perform some processing on that data
3. Third one, sending data to *outbound* channels (consumer stage)

This file is FIRST STAGE i.e PRODUCER
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Subscription interface {
	Updates() <-chan Person
	Close() error
}

func SubscribeDB() Subscription {
	s := subDB{
		done:    make(chan struct{}),
		perChan: make(chan Person),
	}
	go s.loop()
	return &s
}

func SubscribeWeb() Subscription {
	s := subWeb{
		done:    make(chan struct{}),
		perChan: make(chan Person),
	}
	go s.loop()
	return &s
}

/*
This receives data from "db" and returns result
We will integrate this into our concurrent code
*/
func ReceiveDataFromDB() Person {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	age := r1.Intn(20)
	time.Sleep(time.Millisecond * time.Duration(age) * 50)

	return Person{
		Name: "PersonDB",
		Age:  age,
	}
}

func ReceiveDataFromWeb() Person {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	age := r1.Intn(100)
	time.Sleep(time.Millisecond * 10000) // REALLY SLOW CONNECTION

	return Person{
		Name: "PersonWeb",
		Age:  age,
	}
}

type subDB struct {
	perChan chan Person
	done    chan struct{}
}

func (s *subDB) loop() {
	for {
		select {
		case s.perChan <- ReceiveDataFromDB():
		case <-s.done: // a receive operation on a closed channel can always proceed immediately, yielding the element type’s zero value.
			fmt.Printf("SubscriptionDB stop\n\n")
			close(s.perChan)
			return
		}
	}
}

func (s *subDB) Updates() <-chan Person {
	return s.perChan
}

func (s *subDB) Close() error {
	close(s.done)
	return nil
}

type subWeb struct {
	perChan chan Person
	done    chan struct{}
}

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
		case person := <-receiveDone:
			s.perChan <- person
			go receive()
		case <-s.done: // a receive operation on a closed channel can always proceed immediately, yielding the element type’s zero value.
			fmt.Printf("SubscriptionWeb stop\n")
			close(s.perChan)
			return
		}
	}
}

func (s *subWeb) Updates() <-chan Person {
	return s.perChan
}

func (s *subWeb) Close() error {
	close(s.done)
	return nil
}
