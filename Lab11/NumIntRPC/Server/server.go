// NumInt.go
package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type NumInt struct {
	computedPIs ComputedPIs
}

type ComputedPIs struct {
	hash map[int]float64
	mu   sync.Mutex
}

func (c *ComputedPIs) Add(key int, value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hash[key] = value
}

func (c *ComputedPIs) Get(key int) (float64, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, exists := c.hash[key]
	return value, exists
}

func (c *ComputedPIs) Contains(key int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exists := c.hash[key]
	return exists
}

func (n *NumInt) NumInt(args *int, reply *float64) error {
	steps := *args
	if value, exists := n.computedPIs.Get(steps); exists {
		*reply = value
		return nil
	}

	numThreads := 4 // For simplicity, we use a fixed number of threads
	limit := steps / numThreads

	sum := make(chan float64)
	go numIntThreadRecursive(0, steps, 1.0/float64(steps), limit, sum)

	result := <-sum
	pi := result * (1.0 / float64(steps))
	n.computedPIs.Add(steps, pi)
	*reply = pi
	return nil
}

func numIntThreadRecursive(from, to int, step float64, limit int, sum chan float64) {
	workload := to - from
	if workload <= limit {
		mySum := 0.0
		for i := from; i < to; i++ {
			x := (float64(i) + 0.5) * step
			mySum += 4.0 / (1.0 + x*x)
		}
		sum <- mySum
	} else {
		mid := from + workload/2
		leftSum := make(chan float64)
		rightSum := make(chan float64)

		go numIntThreadRecursive(from, mid, step, limit, leftSum)
		go numIntThreadRecursive(mid, to, step, limit, rightSum)

		mySum := <-leftSum + <-rightSum
		sum <- mySum
	}
}

func main() {
	numInt := new(NumInt)
	numInt.computedPIs.hash = make(map[int]float64)
	rpc.Register(numInt)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 1234")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
