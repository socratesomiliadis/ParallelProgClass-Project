package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func main() {
	size := 1000000000
	numThreads := runtime.GOMAXPROCS(0)

	prime := make([]bool, size+1)
	for i := 2; i <= size; i++ {
		prime[i] = true
	}

	start := time.Now()

	limit := int(math.Sqrt(float64(size))) + 1
	done := make(chan bool)

	for i := 0; i < numThreads; i++ {
		go sieveWorker(prime, i, numThreads, limit, size, done)
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

func sieveWorker(prime []bool, rank, numThreads, limit, size int, done chan bool) {
	for p := rank + 2; p < limit; p += numThreads {
		if prime[p] {
			for i := p * p; i <= size; i += p {
				prime[i] = false
			}
		}
	}

	done <- true
}
