package main

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"time"
)

func main() {
	// Circuit input size (number of bits)
	size := 28
	// Number of possible inputs (bit combinations)
	iterations := int(math.Pow(2, float64(size)))

	numGoroutines := runtime.GOMAXPROCS(0)

	// Start timing
	start := time.Now()

	done := make(chan bool, 1)

	chunkSize := iterations / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		from := i * chunkSize
		to := i*chunkSize + chunkSize
		if i == (numGoroutines - 1) {
			to = iterations
		}

		go checkCircuitWorker(from, to, size, done)
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Stop timing
	elapsedTime := time.Since(start)

	fmt.Println("All done")
	fmt.Println("time in ms =", elapsedTime.Milliseconds())
}

func checkCircuitWorker(from int, to int, size int, done chan bool) {
	for i := from; i < to; i++ {
		checkCircuit(i, size)
	}
	done <- true
}

func checkCircuit(z, size int) {
	// Convert z to binary format
	v := make([]bool, size)
	for i := size - 1; i >= 0; i-- {
		v[i] = z&(1<<i) != 0
	}

	// Check the output of the circuit for input v
	value :=
		(v[0] || v[1]) &&
			(!v[1] || !v[3]) &&
			(v[2] || v[3]) &&
			(!v[3] || !v[4]) &&
			(v[4] || !v[5]) &&
			(v[5] || !v[6]) &&
			(v[5] || v[6]) &&
			(v[6] || !v[15]) &&
			(v[7] || !v[8]) &&
			(!v[7] || !v[13]) &&
			(v[8] || v[9]) &&
			(v[8] || !v[9]) &&
			(!v[9] || !v[10]) &&
			(v[9] || v[11]) &&
			(v[10] || v[11]) &&
			(v[12] || v[13]) &&
			(v[13] || !v[14]) &&
			(v[14] || v[15]) &&
			(v[14] || v[16]) &&
			(v[17] || v[1]) &&
			(v[18] || !v[0]) &&
			(v[19] || v[1]) &&
			(v[19] || !v[18]) &&
			(!v[19] || !v[9]) &&
			(v[0] || v[17]) &&
			(!v[1] || v[20]) &&
			(!v[21] || v[20]) &&
			(!v[22] || v[20]) &&
			(!v[21] || !v[20]) &&
			(v[22] || !v[20])

	// If output == true print v and z
	if value {
		printResult(v, size, z)
	}
}

// Printing utility
func printResult(v []bool, size, z int) {
	result := strconv.Itoa(z)

	for i := 0; i < size; i++ {
		if v[i] {
			result += " 1"
		} else {
			result += " 0"
		}
	}

	fmt.Println(result)
}
