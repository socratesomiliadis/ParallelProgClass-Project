package main

import (
	"fmt"
	"runtime"
	"time"
)

func vecAddWorker(a []float64, b []float64, c []float64, myfrom, myto int, done chan bool) {
	for i := myfrom; i < myto; i++ {
		a[i] = b[i] + c[i]
	}
	done <- true
}

func main() {
	size := 1000
	numThreads := runtime.GOMAXPROCS(0)

	a := make([]float64, size)
	b := make([]float64, size)
	c := make([]float64, size)

	for i := 0; i < size; i++ {
		a[i] = 0.0
		b[i] = 1.0
		c[i] = 0.5
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

		go vecAddWorker(a, b, c, from, to, done)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}

	// For debugging
	// for i := 0; i < size; i++ {
	// 	fmt.Println(a[i])
	// }

	elapsedTime := time.Since(start)
	fmt.Println("Time:", elapsedTime.Microseconds(), "μs")
}
