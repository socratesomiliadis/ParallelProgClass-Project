package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	PORT       = "1234"
	numWorkers = 4
	numSteps   = 1000000000
)

type Sum struct {
	result float64
	mu     sync.Mutex
}

func (s *Sum) AddTo(value float64) {
	s.mu.Lock()
	s.result += value
	s.mu.Unlock()
}

func (s *Sum) GetResult() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.result
}

func main() {
	n := &Sum{}
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}
		wg.Add(1)
		go func(id int, conn net.Conn) {
			defer wg.Done()
			handleConnection(conn, id, n, numWorkers, numSteps)
		}(i, conn)
	}
	fmt.Println("All Started")
	wg.Wait()

	step := 1.0 / float64(numSteps)
	pi := n.GetResult() * step

	fmt.Printf("Master-worker program results with %d steps\n", numSteps)
	fmt.Printf("Computed pi = %22.20f\n", pi)
}

func handleConnection(conn net.Conn, id int, s *Sum, numWorkers, numSteps int) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	app := &MasterProtocol{
		mySum:      s,
		myId:       id,
		mySteps:    numSteps,
		numWorkers: numWorkers,
	}
	outmsg := app.PrepareRequest()
	writer.WriteString(outmsg + "\n")
	writer.Flush()
	inmsg, _ := reader.ReadString('\n')
	app.ProcessReply(strings.TrimSpace(inmsg))
}

type MasterProtocol struct {
	mySum      *Sum
	myId       int
	mySteps    int
	numWorkers int
}

func (app *MasterProtocol) PrepareRequest() string {
	return fmt.Sprintf("%d %d", app.mySteps, app.myId)
}

func (app *MasterProtocol) ProcessReply(input string) {
	repl, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println("Error parsing reply:", err)
		return
	}
	app.mySum.AddTo(repl)
}
