package main

import (
	"fmt"
	"time"
)

func main() {
	numOfProds := 100
	numOfCons := 50

	ch := make(chan string)
	done := make(chan bool)

	// Start multiple producers
	for i := 0; i < numOfProds; i++ {
		go producer(i, ch, done)
	}

	// Start multiple consumers
	for i := 0; i < numOfCons; i++ {
		go consumer(i, ch, done)
	}

	// Wait for all producers and consumers to finish
	for i := 0; i < numOfProds+numOfCons; i++ {
		<-done
	}

	close(ch)
}

func producer(index int, ch chan string, done chan bool) {
	for i := 0; i < 5; i++ {
		ch <- fmt.Sprintf("Producer %v send %v", index, i)
	}

	done <- true
}

func consumer(index int, ch chan string, threadDone chan bool) {
	done := false
	for !done {
		select {
		// Receive message from channel
		case msg, ok := <-ch:
			if !ok {
				done = true
			}
			fmt.Printf("Consumer %v Received: %s\n", index, msg)

		// Timeout after 1 second
		case <-time.After(1 * time.Second):
			threadDone <- true
		}

	}

	threadDone <- true
}
