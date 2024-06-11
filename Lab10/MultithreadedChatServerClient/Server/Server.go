package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

const (
	PORT = "1234"
	EXIT = "CLOSE"
)

var clientWriters = struct {
	sync.RWMutex
	writers []*bufio.Writer
}{}

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
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	clientWriters.Lock()
	clientWriters.writers = append(clientWriters.writers, writer)
	clientWriters.Unlock()

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}

		fmt.Println("Received message from client:", message)
		broadcastMessage(message, writer)
	}
}

func broadcastMessage(message string, sender *bufio.Writer) {
	clientWriters.RLock()
	defer clientWriters.RUnlock()

	for _, writer := range clientWriters.writers {
		if writer != sender {
			writer.WriteString(message)
			writer.Flush()
		}
	}
}
