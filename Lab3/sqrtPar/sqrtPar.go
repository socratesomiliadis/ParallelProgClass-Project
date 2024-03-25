package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func sqrtWorker(table []float64, myfrom, myto int, done chan bool) {
	for i := myfrom; i < myto; i++ {
		table[i] = math.Sqrt(table[i])
	}
	done <- true
}

func main() {
	size := 1000
	numThreads := runtime.GOMAXPROCS(0)

	a := make([]float64, size)
	for i := 0; i < size; i++ {
		a[i] = float64(i)
		// a[i] = rand.Float64() // For random
	}

	start := time.Now()

	block := size / numThreads

	done := make(chan bool, 1)

	for i := 0; i < numThreads; i++ {
		from := i * block
		to := i*block + block
		if i == (numThreads - 1) {
			to = size
		}

		go sqrtWorker(a, from, to, done)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}

	elapsedTime := time.Since(start)
	fmt.Println("Time:", elapsedTime.Microseconds(), "μs")
}
