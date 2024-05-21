package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

var (
	totalTasks    int
	tasksAssigned = -1
	lock          sync.Mutex
)

func main() {
	size := 1000000000
	numThreads := runtime.NumCPU()
	prime := make([]bool, size+1)

	for i := 2; i <= size; i++ {
		prime[i] = true
	}

	start := time.Now()
	limit := int(math.Sqrt(float64(size))) + 1
	totalTasks = limit

	done := make(chan bool)

	for i := 0; i < numThreads; i++ {
		go sieveWorker(prime, size, done)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}

	elapsedTime := time.Since(start)

	count := 0
	for i := 2; i <= size; i++ {
		if prime[i] {
			count++
		}
	}

	fmt.Printf("Number of primes: %d\n", count)
	fmt.Printf("Time: %v\n", elapsedTime)
}

func getTask() int {
	lock.Lock()
	defer lock.Unlock()
	tasksAssigned++
	if tasksAssigned < totalTasks {
		return tasksAssigned
	}
	return -1
}

func sieveWorker(prime []bool, size int, done chan bool) {
	for {
		p := getTask()
		if p < 0 {
			break
		}
		if prime[p] {
			for i := p * p; i <= size; i += p {
				prime[i] = false
			}
		}
	}

	done <- true
}
