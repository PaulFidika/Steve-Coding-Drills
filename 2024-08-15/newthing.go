package main

import (
	monkey "blahblah/utils"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, World!")
	result := monkey.AddStuff(2, 3)
	fmt.Println(result)

	var n int = 1_000_000
	var testSlice = make([]int, n)
	var testSlice2 = make([]int, 0, n)
	var testSlice3 = make([]int, 100_000)

	fmt.Println("Test slice: ", timeLoop1(testSlice, n))
	fmt.Println("Test slice 2: ", timeLoop2(testSlice2, n))
	fmt.Println("Test slice 3: ", timeLoop2(testSlice3, n))
}

func timeLoop1(slice []int, n int) time.Duration {
	var t0 = time.Now()
	for i := 0; i < n; i++ {
		slice[i] = 1
	}
	return time.Since(t0)
}

func timeLoop2(slice []int, n int) time.Duration {
	var t0 = time.Now()
	for len(slice) < n {
		slice = append(slice, 1)
	}
	return time.Since(t0)
}

