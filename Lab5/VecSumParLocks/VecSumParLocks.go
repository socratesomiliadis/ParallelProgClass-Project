package main

import (
	"fmt"
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

func VecSumWorker(id int, a []float64, sum *SafeSum, numThreads int, done chan bool) {
	myStart := id * (len(a) / numThreads)
	myStop := myStart + (len(a) / numThreads)

	if id == (numThreads - 1) {
		myStop = len(a)
	}

	mySum := 0.0
	for j := myStart; j < myStop; j++ {
		mySum += a[j]
	}

	sum.add(mySum)

	done <- true
}

func main() {
	size := 100000
	numThreads := runtime.NumCPU()

	a := make([]float64, size)
	for i := 0; i < size; i++ {
		a[i] = float64(i)
	}

	sum := SafeSum{sum: 0.0}

	start := time.Now()

	done := make(chan bool)

	for i := 0; i < numThreads; i++ {
		go VecSumWorker(i, a, &sum, numThreads, done)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}

	elapsedTime := time.Since(start)

	fmt.Println("sum =", sum.sum)
	fmt.Println("time in ms =", elapsedTime.Milliseconds())
}
