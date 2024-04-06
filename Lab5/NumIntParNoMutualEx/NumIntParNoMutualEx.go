package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func main() {
	numSteps := 100000000
	numThreads := runtime.NumCPU()

	a := make([]float64, numThreads)

	startTime := time.Now()

	done := make(chan bool, numThreads)

	for i := 0; i < numThreads; i++ {
		go NumIntWorker(i, numSteps, a, numThreads, done)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}

	sum := 0.0
	for _, val := range a {
		sum += val
	}

	step := 1.0 / float64(numSteps)
	pi := sum * step

	endTime := time.Now()

	fmt.Printf("sequential program results with %d steps\n", numSteps)
	fmt.Printf("computed pi = %22.20f\n", pi)
	fmt.Printf("difference between estimated pi and math.Pi = %22.20f\n", math.Abs(pi-math.Pi))
	fmt.Printf("time to compute = %f seconds\n", endTime.Sub(startTime).Seconds())
}

func NumIntWorker(id int, numSteps int, a []float64, numThreads int, done chan bool) {
	myStart := id * (numSteps / numThreads)
	myStop := myStart + (numSteps / numThreads)

	if id == (numThreads - 1) {
		myStop = numSteps
	}

	mySum := 0.0
	myStep := 1.0 / float64(numSteps)

	for j := myStart; j < myStop; j++ {
		x := (float64(j) + 0.5) * myStep
		mySum += 4.0 / (1.0 + x*x)
	}

	a[id] = mySum

	done <- true
}
