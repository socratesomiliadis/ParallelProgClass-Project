package main

import "fmt"

func main() {
	size := 1000

	a := make([]float64, size)
	b := make([]float64, size)
	c := make([]float64, size)

	for i := 0; i < size; i++ {
		a[i] = 0.0
		b[i] = 1.0
		c[i] = 0.5
	}

	for i := 0; i < size; i++ {
		a[i] = b[i] + c[i]
	}

	// For debugging
	for i := 0; i < size; i++ {
		fmt.Println(a[i])
	}
}
