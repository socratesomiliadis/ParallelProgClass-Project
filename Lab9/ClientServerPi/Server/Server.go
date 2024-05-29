package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"
)

const (
	PORT = "1234"
)

type Request struct {
	Steps int `json:"offset,omitempty"`
}

type Response struct {
	Message string `json:"message"`
}

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
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port:", PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Received request from", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	in := bufio.NewReader(conn)
	out := bufio.NewWriter(conn)

	inmsg, _ := in.ReadString('\n')
	inmsg = strings.TrimSpace(inmsg)

	var req Request
	json.Unmarshal([]byte(inmsg), &req)

	outmsg := processRequest(req)
	resData, _ := json.Marshal(outmsg)
	out.WriteString(string(resData) + "\n")
	out.Flush()

	conn.Close()
	fmt.Println("Data socket closed")
}

func processRequest(req Request) Response {
	fmt.Println("Received message from client:", req)

	numSteps := req.Steps
	numThreads := runtime.NumCPU()
	limit := numSteps / numThreads

	step := 1.0 / float64(numSteps)
	wg := &sync.WaitGroup{}
	sum := SafeSum{sum: 0.0}

	wg.Add(1)
	go (&NumIntRecursive{0, numSteps, step, limit, 0}).compute(wg, &sum)

	wg.Wait()

	pi := sum.sum * step

	res := Response{Message: fmt.Sprintf("Pi is approximately %.16f", pi)}

	fmt.Println("Send message to client:", res.Message)
	return res
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
