package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	// Loading the file content
	content, err := os.ReadFile("E128.coli")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	text := string(content)

	start := time.Now()

	patternString := "tacccagattatcgccatcaacgg"
	n := len(text)
	m := len(patternString)

	matchLength := n - m
	match := make([]byte, matchLength+1)

	numThreads := runtime.GOMAXPROCS(0)
	block := matchLength / numThreads
	counter := Counter{val: 0}

	// Spawning Goroutines
	done := make(chan bool)
	for i := 0; i < numThreads; i++ {
		from := i * block
		to := from + block
		if i == numThreads-1 {
			to = matchLength
		}

		go StringMatchWorker(match, from, to, m, patternString, text, &counter, done)
	}

	// Waiting for all Goroutines to finish
	for i := 0; i < numThreads; i++ {
		<-done
	}

	elapsedTime := time.Since(start)
	fmt.Printf("Time: %v\n", elapsedTime)
	for i := 0; i < matchLength; i++ {
		if match[i] == '1' {
			fmt.Printf("%d ", i)
		}
	}
	fmt.Println()
	fmt.Printf("Total matches %d\n", counter.get())
}

// Counter struct with a Mutex for synchronization
type Counter struct {
	val  int
	lock sync.Mutex
}

func (c *Counter) inc() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.val++
}

func (c *Counter) get() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.val
}

func StringMatchWorker(table []byte, from int, to int, patternLength int, pattern string, text string, counter *Counter, done chan bool) {
	for i := from; i < to; i++ {
		j := 0
		for j < patternLength && pattern[j] == text[i+j] {
			j++
		}
		if j >= patternLength {
			table[i] = '1'
			counter.inc()
		}
	}

	done <- true
}
