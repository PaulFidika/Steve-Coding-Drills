package main

import (
	"fmt"
	"time"
)

func main() {
	countTo100Million()
}

func countTo100Million() {
	start := time.Now()
	var res int

	for i := 1; i <= 6_000_000_000; i++ {
		res += 1
	}

	duration := time.Since(start)
	fmt.Printf("Counting to 100,000,000 took %v\n", duration)
}