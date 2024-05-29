package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var acceptedVerbs = map[string]bool{
	"LOWERCASE": true,
	"UPPERCASE": true,
	"ENCODE":    true,
	"DECODE":    true,
	"CLOSE":     true,
}

type ClientProtocol struct {
	user *bufio.Reader
}

func NewClientProtocol() *ClientProtocol {
	return &ClientProtocol{user: bufio.NewReader(os.Stdin)}
}

func (cp *ClientProtocol) prepareRequest() (string, error) {
	fmt.Print("Select action (LOWERCASE, UPPERCASE, ENCODE, DECODE, CLOSE): ")
	action, _ := cp.user.ReadString('\n')
	action = strings.TrimSpace(strings.ToUpper(action))
	for !acceptedVerbs[action] {
		fmt.Print("Invalid action. Select action (LOWERCASE, UPPERCASE, ENCODE, DECODE, CLOSE): ")
		action, _ = cp.user.ReadString('\n')
		action = strings.TrimSpace(strings.ToUpper(action))
	}

	if action == "CLOSE" {
		return action, nil
	}

	if action == "ENCODE" || action == "DECODE" {
		fmt.Print("Enter offset: ")
		offset, _ := cp.user.ReadString('\n')
		offset = strings.TrimSpace(offset)
		for !isNumeric(offset) {
			fmt.Print("Invalid offset. Enter offset: ")
			offset, _ = cp.user.ReadString('\n')
			offset = strings.TrimSpace(offset)
		}

		action += " " + offset
	}

	fmt.Print("Enter message: ")
	message, _ := cp.user.ReadString('\n')
	message = strings.TrimSpace(message)

	return action + " " + message, nil
}

func (cp *ClientProtocol) processReply(theInput string) {
	fmt.Println("Message received from server:", theInput)
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

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

	in := bufio.NewReader(conn)
	out := bufio.NewWriter(conn)

	cp := NewClientProtocol()

	outmsg, err := cp.prepareRequest()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		return
	}
	for outmsg != EXIT {
		out.WriteString(outmsg + "\n")
		out.Flush()
		inmsg, _ := in.ReadString('\n')
		cp.processReply(strings.TrimSpace(inmsg))
		outmsg, err = cp.prepareRequest()
		if err != nil {
			fmt.Println("Error preparing request:", err)
			return
		}
	}
	out.WriteString(outmsg + "\n")
	out.Flush()

	fmt.Println("Data Socket closed")
}
