package main

import "fmt"

func main() {
	size := 10

	a := make([][]float64, size)
	b := make([][]float64, size)
	c := make([][]float64, size)

	for i := 0; i < size; i++ {
		a[i] = make([]float64, size)
		b[i] = make([]float64, size)
		c[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			a[i][j] = 0.1
			b[i][j] = 0.3
			c[i][j] = 0.5
		}
	}

	// Matrix addition
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			a[i][j] = b[i][j] + c[i][j]
		}
	}

	// For debugging
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			fmt.Print(a[i][j], " ")
		}
		fmt.Println()
	}
}
