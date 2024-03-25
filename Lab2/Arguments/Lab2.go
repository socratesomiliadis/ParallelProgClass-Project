package main

import (
	"fmt"
)

func main() {
	numThreads := 4 //Τοπική μεταβλητή της main

	n := 0                       //Τοπική μεταβλητή της main
	a := make([]int, numThreads) //a τοπική μεταβλητή της main, numThreads όρισμα τιμής

	done := make(chan bool, 1)

	for i := 0; i < numThreads; i++ {
		go counter(i, n, a, done) //i όρισμα τιμής στοιχισμενο, n όρισμα τιμής, a και done ορίσματα αναφοράς
	}

	for i := 0; i < numThreads; i++ {
		<-done
		fmt.Printf("a[%d] = %d\n", i, a[i])
	}

	fmt.Println("n = ", n)

	for i := 0; i < numThreads; i++ {
		fmt.Printf("a[%d] = %d\n", i, a[i])
	}
}

func counter(threadID int, n int, a []int, done chan bool) {
	n = n + threadID
	a[threadID] = a[threadID] + n
	fmt.Printf("Thread %d my a = %d\n", threadID, a[threadID])
	done <- true
}
