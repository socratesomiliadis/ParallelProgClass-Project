package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	numSteps := int64(10000)
	sum := 0.0

	/* start timing */
	startTime := time.Now()

	step := 1.0 / float64(numSteps)
	/* do computation */
	for i := int64(0); i < numSteps; i++ {
		x := (float64(i) + 0.5) * step
		sum += 4.0 / (1.0 + x*x)
	}
	pi := sum * step

	/* end timing and print result */
	endTime := time.Now()
	fmt.Printf("sequential program results with %d steps\n", numSteps)
	fmt.Printf("computed pi = %22.20f\n", pi)
	fmt.Printf("difference between estimated pi and math.Pi = %22.20f\n", math.Abs(pi-math.Pi))
	fmt.Printf("time to compute = %f seconds\n", endTime.Sub(startTime).Seconds())
}
