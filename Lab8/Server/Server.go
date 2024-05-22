package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	PORT = "1234"
	EXIT = "CLOSE"
)

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

		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	in := bufio.NewReader(conn)
	out := bufio.NewWriter(conn)

	inmsg, _ := in.ReadString('\n')
	inmsg = strings.TrimSpace(inmsg)
	outmsg := processRequest(inmsg)
	for inmsg != EXIT {
		out.WriteString(outmsg + "\n")
		out.Flush()
		inmsg, _ = in.ReadString('\n')
		inmsg = strings.TrimSpace(inmsg)
		outmsg = processRequest(inmsg)
	}
	conn.Close()
	fmt.Println("Data socket closed")
}

func processRequest(theInput string) string {
	fmt.Println("Received message from client:", theInput)

	action := getFirstWord(theInput)
	theOutput := ""

	switch strings.ToUpper(action) {
	case "LOWERCASE":
		theOutput = strings.ToLower(theInput[len(action)+1:])
	case "UPPERCASE":
		theOutput = strings.ToUpper(theInput[len(action)+1:])
	case "ENCODE":
		offset := findFirstInteger(theInput)
		message := theInput[len(action)+2:]
		theOutput = encode(message, offset)
	case "DECODE":
		offset := findFirstInteger(theInput)
		message := theInput[len(action)+2:]
		theOutput = decode(message, offset)
	}

	fmt.Println("Send message to client:", theOutput)
	return theOutput
}

func getFirstWord(text string) string {
	words := strings.Fields(text)
	if len(words) > 0 {
		return words[0]
	}
	return text
}

func findFirstInteger(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			j := i
			for j < len(s) && s[j] >= '0' && s[j] <= '9' {
				j++
			}
			num, _ := strconv.Atoi(s[i:j])
			return num
		}
	}
	return 0
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
