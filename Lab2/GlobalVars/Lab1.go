package main

import (
	"fmt"
)

var n int   //Καθολική μεταβλητή
var a []int //Καθολική μεταβλητή

func main() {
	numThreads := 4 //Τοπική μεταβλητή της main

	a = make([]int, numThreads) //numThreads όρισμα τιμής

	done := make(chan bool, 1)
	threadN := make(chan int, 1)

	for i := 0; i < numThreads; i++ {
		go counter(i, done, threadN) //i όρισμα τιμής στοιχισμενο, done και threadN ορίσματα αναφοράς
	}

	for i := 0; i < numThreads; i++ {
		<-done
		n = n + <-threadN
		fmt.Printf("a[%d] = %d\n", i, a[i])
	}

	fmt.Println("n = ", n)
}

func counter(threadID int, done chan bool, finalThreadN chan int) { // Launch goroutines with captured threadID

	threadN := n //Τοπική μεταβλητή της counter

	threadN = threadN + 1 + threadID
	a[threadID] = a[threadID] + 1 + threadID
	fmt.Printf("Thread %d n = %d a[%d] = %d\n", threadID, n, threadID, a[threadID])
	finalThreadN <- threadN
	done <- true
}
