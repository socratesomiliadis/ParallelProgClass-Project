package main

import (
	"fmt"
)

func routineA(id int, done chan bool) {

	fmt.Printf("Hello world from thread %d of routine A\n", id)
	fmt.Printf("Bye world from thread %d of routine A\n", id)

	done <- true
}

func routineB(id int, done chan bool) {
	fmt.Printf("Hello world from thread %d of routine B\n", id)
	fmt.Printf("Bye world from thread %d of routine B\n", id)

	done <- true
}

func main() {
	fmt.Println("Hello world from main")
	done := make(chan bool, 1)
	
	for i := 0; i < 3; i++ {
		go routineA(i, done)
		go routineB(i, done)
	}	
    for i := 0; i < 6; i++ {
        <-done
    } 
     
    fmt.Println("Bye world from main")
    close(done) 
}
