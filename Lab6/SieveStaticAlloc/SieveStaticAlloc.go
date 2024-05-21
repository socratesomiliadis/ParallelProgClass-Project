package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func main() {
	size := 100000000
	numThreads := runtime.GOMAXPROCS(0)

	prime := make([]bool, size+1)

	for i := 2; i <= size; i++ {
		prime[i] = true
	}

	start := time.Now()

	limit := int(math.Sqrt(float64(size))) + 1
	block := (limit - 2) / numThreads

	done := make(chan bool)

	for i := 0; i < numThreads; i++ {
		from := i * block
		to := (i + 1) * block
		if i == 0 {
			from = 2
		}
		if i == numThreads-1 {
			to = limit
		}

		go sieveWorker(prime, from, to, size, done)
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

func sieveWorker(prime []bool, from int, to int, size int, done chan bool) {
	for p := from; p < to; p++ {
		if prime[p] {
			for i := p * p; i <= size; i += p {
				prime[i] = false
			}
		}
	}

	done <- true
}
