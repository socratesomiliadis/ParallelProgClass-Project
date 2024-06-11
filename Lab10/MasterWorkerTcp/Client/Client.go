package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	HOST       = "localhost"
	PORT       = "1234"
	numWorkers = 4
)

func main() {
	conn, err := net.Dial("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	inmsg, _ := reader.ReadString('\n')
	app := &WorkerProtocol{numWorkers: numWorkers}
	outmsg := app.Compute(strings.TrimSpace(inmsg))
	writer.WriteString(outmsg + "\n")
	writer.Flush()
}

type WorkerProtocol struct {
	numWorkers int
}

func (app *WorkerProtocol) Compute(input string) string {
	parts := strings.Fields(input)
	numSteps, _ := strconv.Atoi(parts[0])
	id, _ := strconv.Atoi(parts[1])

	fmt.Printf("Worker %d calculates %d\n", id, numSteps)
	myStart := id * (numSteps / app.numWorkers)
	myStop := myStart + (numSteps / app.numWorkers)
	if id == (app.numWorkers - 1) {
		myStop = numSteps
	}

	myStep := 1.0 / float64(numSteps)
	fmt.Printf("Worker %d sums from %d to %d\n", id, myStart, myStop)

	var sum float64
	for i := myStart; i < myStop; i++ {
		x := (float64(i) + 0.5) * myStep
		sum += 4.0 / (1.0 + x*x)
	}

	theOutput := fmt.Sprintf("%f", sum)
	fmt.Printf("Worker %d result %s\n", id, theOutput)
	return theOutput
}
