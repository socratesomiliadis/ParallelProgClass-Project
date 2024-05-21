package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Buffer struct {
	contents []int
	size     int
	front    int
	back     int
	counter  int
	lock     sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
}

func NewBuffer(size int) *Buffer {
	buff := &Buffer{
		size:     size,
		contents: make([]int, size),
		front:    0,
		back:     size - 1,
		counter:  0,
	}
	buff.notFull = sync.NewCond(&buff.lock)
	buff.notEmpty = sync.NewCond(&buff.lock)
	return buff
}

func (buff *Buffer) Put(data int) {
	buff.lock.Lock()
	defer buff.lock.Unlock()

	for buff.counter == buff.size {
		fmt.Println("The buffer is full")
		buff.notFull.Wait()
	}

	buff.back = (buff.back + 1) % buff.size
	buff.contents[buff.back] = data
	buff.counter++
	fmt.Printf("Prod %d No %d Loc %d Count = %d\n", getGID(), data, buff.back, buff.counter)

	if buff.counter == 1 {
		buff.notEmpty.Signal()
	}
}

func (buff *Buffer) Get() int {
	buff.lock.Lock()
	defer buff.lock.Unlock()

	for buff.counter == 0 {
		fmt.Println("The buffer is empty")
		buff.notEmpty.Wait()
	}

	data := buff.contents[buff.front]
	fmt.Printf("  Cons %d No %d Loc %d Count = %d\n", getGID(), data, buff.front, buff.counter-1)
	buff.front = (buff.front + 1) % buff.size
	buff.counter--

	if buff.counter == buff.size-1 {
		buff.notFull.Signal()
	}

	return data
}

func producer(buff *Buffer, reps, scale int, done chan bool) {
	for i := 0; i < reps; i++ {
		buff.Put(i)
		time.Sleep(time.Duration(rand.Intn(scale)) * time.Millisecond)
	}

	done <- true
}

func consumer(buff *Buffer, scale int, done chan bool) {
	defer func() { done <- true }()
	for {
		time.Sleep(time.Duration(rand.Intn(scale)) * time.Millisecond)
		buff.Get()
	}
}

func main() {
	bufferSize := 5
	noIterations := 20
	producerDelay := 100
	consumerDelay := 1
	noProds := 3
	noCons := 2

	buff := NewBuffer(bufferSize)

	done := make(chan bool)

	for i := 0; i < noProds; i++ {
		go producer(buff, noIterations, producerDelay, done)
	}

	for j := 0; j < noCons; j++ {
		go consumer(buff, consumerDelay, done)
	}

	for i := 0; i < noProds+noCons; i++ {
		<-done
	}
}

func getGID() uint64 {
	b := make([]byte, 64)
	n := runtime.Stack(b, false)
	idField := strings.Fields(strings.TrimPrefix(string(b[:n]), "goroutine "))[0]
	id, err := strconv.ParseUint(idField, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
