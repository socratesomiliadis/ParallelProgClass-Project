package main

import (
	"fmt"
	"sync"
)

func routineA(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Hello world from thread %d of routine A\n", id)
	fmt.Printf("Bye world from thread %d of routine A\n", id)
}

func routineB(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Hello world from thread %d of routine B\n", id)
	fmt.Printf("Bye world from thread %d of routine B\n", id)
}

func main() {
	fmt.Println("Hello world from main thread")
	var wg sync.WaitGroup
	
	for i := 0; i < 3; i++ {
		
		wg.Add(1)
		go routineA(i, &wg)

		wg.Add(1)
		go routineB(i, &wg)
	}
	
	for i := 0; i < 6; i++ {
        wg.Wait()
    }
    
    fmt.Println("Bye world from main")    
}
