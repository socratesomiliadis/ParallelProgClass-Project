package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	const (
		HOST = "localhost"
		PORT = "1234"
		EXIT = "CLOSE"
	)

	conn, err := net.Dial("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connection to", HOST, "established")

	var wg sync.WaitGroup

	wg.Add(1)
	go sendMessages(conn, &wg)

	wg.Add(1)
	go receiveMessages(conn, &wg)

	wg.Wait()

	fmt.Println("Data Socket closed")
}

func sendMessages(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	writer := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Send message, CLOSE for exit: ")
		scanner.Scan()
		text := scanner.Text()
		if strings.ToUpper(text) == "CLOSE" {
			break
		}
		writer.WriteString(text + "\n")
		writer.Flush()
	}
	conn.Close()
}

func receiveMessages(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed")
			break
		}
		fmt.Println("\nReceived message:", message)
		fmt.Print("Send message, CLOSE for exit: ")
	}
}
