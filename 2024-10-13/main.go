package main

import (
	"fmt"
	whatever "steve-fidika/2024-10-13/package_1"
	"time"
)

func main() {
	start := time.Now()

	result := whatever.RecursiveFibonacci(35)
	
	duration := time.Since(start)
	fmt.Printf("35th Fibonacci Sequence item: %d\n", result)
	fmt.Printf("Time taken: %v\n", duration)
}