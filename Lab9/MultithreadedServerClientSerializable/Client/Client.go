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

var acceptedVerbs = map[string]bool{
	"LOWERCASE": true,
	"UPPERCASE": true,
	"ENCODE":    true,
	"DECODE":    true,
	"CLOSE":     true,
}

type Request struct {
	Action  string `json:"action"`
	Offset  int    `json:"offset,omitempty"`
	Message string `json:"message,omitempty"`
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

	fmt.Print("Select action (LOWERCASE, UPPERCASE, ENCODE, DECODE, CLOSE): ")
	action, _ := cp.user.ReadString('\n')
	action = strings.TrimSpace(strings.ToUpper(action))
	for !acceptedVerbs[action] {
		fmt.Print("Invalid action. Select action (LOWERCASE, UPPERCASE, ENCODE, DECODE, CLOSE): ")
		action, _ = cp.user.ReadString('\n')
		action = strings.TrimSpace(strings.ToUpper(action))
	}

	req.Action = action

	if action == "CLOSE" {
		return req, nil
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

		req.Offset, _ = strconv.Atoi(offset)
	}

	fmt.Print("Enter message: ")
	message, _ := cp.user.ReadString('\n')
	message = strings.TrimSpace(message)
	req.Message = message

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

	for {
		req, err := cp.prepareRequest()
		if err != nil {
			fmt.Println("Error preparing request:", err)
			return
		}

		if req.Action == "CLOSE" {
			break
		}

		reqData, _ := json.Marshal(req)
		out.WriteString(string(reqData) + "\n")
		out.Flush()

		inmsg, _ := in.ReadString('\n')
		cp.processReply([]byte(strings.TrimSpace(inmsg)))
	}

	req := Request{Action: "CLOSE"}
	reqData, _ := json.Marshal(req)
	out.WriteString(string(reqData) + "\n")
	out.Flush()

	fmt.Println("Data Socket closed")
}
