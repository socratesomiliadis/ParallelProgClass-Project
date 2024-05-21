package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

type SafeSum struct {
	sum  float64
	lock sync.Mutex
}

func (s *SafeSum) add(localsum float64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.sum += localsum
}

func main() {
	numSteps := 1000000000
	numThreads := runtime.NumCPU()
	limit := numSteps / numThreads

	startTime := time.Now()

	step := 1.0 / float64(numSteps)
	wg := &sync.WaitGroup{}
	sum := SafeSum{sum: 0.0}

	wg.Add(1)
	go (&NumIntRecursive{0, numSteps, step, limit, 0}).compute(wg, &sum)

	wg.Wait()

	pi := sum.sum * step

	endTime := time.Now()
	fmt.Printf("Results with %d steps\n", numSteps)
	fmt.Printf("Computed pi = %22.20f\n", pi)
	fmt.Printf("Difference between estimated pi and Math.PI = %22.20f\n", math.Abs(pi-math.Pi))
	fmt.Printf("Time to compute = %f seconds\n", endTime.Sub(startTime).Seconds())
}

type NumIntRecursive struct {
	start  int
	stop   int
	step   float64
	limit  int
	result float64
}

func (n *NumIntRecursive) compute(wg *sync.WaitGroup, sum *SafeSum) {
	defer wg.Done()
	mySum := 0.0
	workload := n.stop - n.start
	if workload <= n.limit {
		for i := n.start; i < n.stop; i++ {
			x := (float64(i) + 0.5) * n.step
			mySum += 4.0 / (1.0 + x*x)
		}
		sum.add(mySum)
	} else {
		mid := n.start + workload/2
		wg.Add(2)
		go (&NumIntRecursive{n.start, mid, n.step, n.limit, 0}).compute(wg, sum)
		go (&NumIntRecursive{mid, n.stop, n.step, n.limit, 0}).compute(wg, sum)
	}
}
