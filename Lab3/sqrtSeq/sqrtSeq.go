package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	size := 10000

	a := make([]float64, size)
	for i := 0; i < size; i++ {
		a[i] = float64(i)
		// a[i] = rand.Float64() // For random
	}

	start := time.Now()

	for i := 0; i < size; i++ {
		a[i] = math.Sqrt(a[i])
	}

	elapsedTime := time.Since(start)
	fmt.Println("Time:", elapsedTime.Microseconds(), "μs")
}
