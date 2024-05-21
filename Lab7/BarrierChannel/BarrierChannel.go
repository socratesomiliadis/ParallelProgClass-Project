package main

import "fmt"

func main() {
	barriers := Barrier(10)

	function := func(n int) {
		// Generate goroutines which print "Hi from i" and "Bye from i"
		for i := range barriers {
			go func(i int, ch chan bool) {
				for j := 0; j < n; j++ {
					// Print all "Hi from i" before "Bye from i"
					fmt.Println("Hi from", i)
					ch <- false
					<-ch

					fmt.Println("Bye from", i)
					ch <- true
					<-ch
				}
			}(i, barriers[i])
		}
	}

	callback := func() {
		// Wait for all "Hi from i" printed
		for i := range barriers {
			<-barriers[i]
		}

		SignalAll(barriers)

		// Wait for all "Bye from i" printed
		for i := range barriers {
			<-barriers[i]
		}

		SignalAll(barriers)
	}

	Sync(10, function, callback)
}

func Barrier(i int) []chan bool {
	// Return slice of #i bool channel(s)
	channels := make([]chan bool, i)
	for i := range channels {
		channels[i] = make(chan bool)
	}
	return channels
}

func Sync(n int, function func(int), callback func()) {
	function(n)
	for i := 0; i < n; i++ {
		callback()
	}
}

func SignalAll(barriers []chan bool) {
	// Send signal to all channel
	for i := range barriers {
		barriers[i] <- true
	}
}
