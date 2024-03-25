package main

import (
	"fmt"
	"runtime"
	"time"
)

func matrixAddWorker(a [][]float64, b [][]float64, c [][]float64, myfrom, myto int, size int, done chan bool) {
	for i := myfrom; i < myto; i++ {
		for j := 0; j < size; j++ {
			a[i][j] = b[i][j] + c[i][j]
		}
	}
	done <- true
}

func main() {
	size := 10
	numThreads := runtime.GOMAXPROCS(0)

	a := make([][]float64, size)
	b := make([][]float64, size)
	c := make([][]float64, size)

	for i := 0; i < size; i++ {
		a[i] = make([]float64, size)
		b[i] = make([]float64, size)
		c[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			a[i][j] = 0.1
			b[i][j] = 0.3
			c[i][j] = 0.5
		}
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

		go matrixAddWorker(a, b, c, from, to, size, done)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}

	// For debugging
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			fmt.Print(a[i][j], " ")
		}
		fmt.Println()
	}

	elapsedTime := time.Since(start)
	fmt.Println("Time:", elapsedTime.Microseconds(), "μs")
}
