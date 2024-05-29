package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	Steps int `json:"offset,omitempty"`
}

type Response struct {
	Message string `json:"message"`
}

type ClientProtocol struct {
	user *bufio.Reader
}

func NewClientProtocol() *ClientProtocol {
	return &ClientProtocol{user: bufio.NewReader(os.Stdin)}
}

func (cp *ClientProtocol) prepareRequest() (Request, error) {
	var req Request

	fmt.Print("Enter number of steps: ")
	steps, _ := cp.user.ReadString('\n')
	steps = strings.TrimSpace(steps)
	for !isNumeric(steps) {
		fmt.Print("Invalid number. Enter number of steps: ")
		steps, _ = cp.user.ReadString('\n')
		steps = strings.TrimSpace(steps)
	}

	req.Steps, _ = strconv.Atoi(steps)

	return req, nil
}

func (cp *ClientProtocol) processReply(reply []byte) {
	var res Response
	json.Unmarshal(reply, &res)
	fmt.Println("Message received from server:", res.Message)
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func main() {
	const (
		HOST = "localhost"
		PORT = "1234"
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

	req, err := cp.prepareRequest()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		return
	}

	reqData, _ := json.Marshal(req)
	out.WriteString(string(reqData) + "\n")
	out.Flush()

	inmsg, _ := in.ReadString('\n')
	cp.processReply([]byte(strings.TrimSpace(inmsg)))

	fmt.Println("Data Socket closed")
}
