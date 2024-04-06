package main

import (
	"fmt"
	"sync"
)

func main() {
	end := 1000
	array := make([]int, end)
	var arrLock sync.Mutex
	numThreads := 4

	done := make(chan bool, 1)

	for i := 0; i < numThreads; i++ {
		go counterWorker(array, end, done, &arrLock)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}

	checkArray(array, end, numThreads)
}

func checkArray(array []int, end, numThreads int) {
	errors := 0

	fmt.Println("Checking...")

	for i := 0; i < end; i++ {
		if array[i] != numThreads*i {
			errors++
			fmt.Printf("%d: %d should be %d\n", i, array[i], numThreads*i)
		}
	}

	fmt.Println(errors, "errors.")
}

func counterWorker(array []int, end int, done chan bool, arrLock *sync.Mutex) {
	for i := 0; i < end; i++ {
		for j := 0; j < i; j++ {
			arrLock.Lock()
			array[i]++
			arrLock.Unlock()
		}
	}

	done <- true
}
