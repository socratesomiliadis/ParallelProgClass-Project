package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type linearBarrier struct {
	totalThreads int
	arrived      int
	waiting      bool
	lock         sync.Mutex
	cond         *sync.Cond
}

func newLinearBarrier(total int) *linearBarrier {
	b := &linearBarrier{
		totalThreads: total,
		waiting:      true,
	}
	b.cond = sync.NewCond(&b.lock)
	return b
}

func (b *linearBarrier) barrier() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.arrived++
	if b.arrived == b.totalThreads {
		b.waiting = false
		b.cond.Broadcast()
	}

	for b.waiting {
		b.cond.Wait()
	}

	b.arrived--
	if b.arrived == 0 {
		b.waiting = true
		b.cond.Broadcast()
	}
}

func testThread(id int, b *linearBarrier, done chan bool) {
	defer func() { done <- true }()
	for {
		fmt.Printf("Thread%d started\n", id)
		b.barrier()
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		fmt.Printf("Thread%d reached barrier\n", id)
		b.barrier()
		fmt.Printf("Thread%d passed barrier\n", id)
		b.barrier()
	}
}

func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Println("Usage: go run barrierMain.go <number of threads>")
	// 	os.Exit(1)
	// }

	numThreads := 256
	// if err != nil {
	// 	fmt.Println("Integer argument expected")
	// 	os.Exit(1)
	// }

	testBarrier := newLinearBarrier(numThreads)

	done := make(chan bool)
	for i := 0; i < numThreads; i++ {

		go testThread(i, testBarrier, done)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}
}
