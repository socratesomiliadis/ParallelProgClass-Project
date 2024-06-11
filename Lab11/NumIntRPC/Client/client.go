// client.go
package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error connecting to RPC server:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter number of steps: ")
		nSteps, _ := reader.ReadString('\n')
		nSteps = strings.TrimSpace(nSteps)

		if isNumeric(nSteps) {
			numSteps, _ := strconv.Atoi(nSteps)
			var result float64

			err = client.Call("NumInt.NumInt", numSteps, &result)
			if err != nil {
				fmt.Println("RPC error:", err)
				return
			}

			fmt.Printf("Computed pi = %22.20f\n", result)
			break
		} else {
			fmt.Println("Invalid number. Please enter a valid number of steps.")
		}
	}
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
