package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	end := 1000
	array := make([]int64, end)
	numThreads := 4

	done := make(chan bool, 1)

	for i := 0; i < numThreads; i++ {
		go counterWorker(array, end, done)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}

	checkArray(array, end, numThreads)
}

func checkArray(array []int64, end, numThreads int) {
	errors := 0

	fmt.Println("Checking...")

	for i := 0; i < end; i++ {
		if array[i] != int64(numThreads*i) {
			errors++
			fmt.Printf("%d: %d should be %d\n", i, array[i], numThreads*i)
		}
	}

	fmt.Println(errors, "errors.")
}

func counterWorker(array []int64, end int, done chan bool) {
	for i := 0; i < end; i++ {
		for j := 0; j < i; j++ {
			atomic.AddInt64(&array[i], 1)
		}
	}

	done <- true
}
