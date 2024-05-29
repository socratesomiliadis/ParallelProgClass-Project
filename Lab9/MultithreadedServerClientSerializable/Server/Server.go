package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

const (
	PORT = "1234"
)

type Request struct {
	Action  string `json:"action"`
	Offset  int    `json:"offset,omitempty"`
	Message string `json:"message,omitempty"`
}

type Response struct {
	Message string `json:"message"`
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

	for {
		inmsg, _ := in.ReadString('\n')
		inmsg = strings.TrimSpace(inmsg)

		var req Request
		json.Unmarshal([]byte(inmsg), &req)

		if req.Action == "CLOSE" {
			break
		}

		outmsg := processRequest(req)
		resData, _ := json.Marshal(outmsg)
		out.WriteString(string(resData) + "\n")
		out.Flush()
	}
	conn.Close()
	fmt.Println("Data socket closed")
}

func processRequest(req Request) Response {
	fmt.Println("Received message from client:", req)

	var res Response
	switch strings.ToUpper(req.Action) {
	case "LOWERCASE":
		res.Message = strings.ToLower(req.Message)
	case "UPPERCASE":
		res.Message = strings.ToUpper(req.Message)
	case "ENCODE":
		res.Message = encode(req.Message, req.Offset)
	case "DECODE":
		res.Message = decode(req.Message, req.Offset)
	}

	fmt.Println("Send message to client:", res.Message)
	return res
}

func encode(message string, offset int) string {
	var result strings.Builder
	for _, char := range message {
		if char != ' ' {
			originalAlphabetPosition := char - 'a'
			newAlphabetPosition := (originalAlphabetPosition + rune(offset)) % 26
			newCharacter := 'a' + newAlphabetPosition
			result.WriteRune(newCharacter)
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}

func decode(message string, offset int) string {
	var result strings.Builder
	for _, char := range message {
		if char != ' ' {
			originalAlphabetPosition := char - 'a'
			newAlphabetPosition := (originalAlphabetPosition - rune(offset)) % 26
			if newAlphabetPosition < 0 {
				newAlphabetPosition += 26
			}
			newCharacter := 'a' + newAlphabetPosition
			result.WriteRune(newCharacter)
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}
