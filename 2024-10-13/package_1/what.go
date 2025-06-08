package whatever

import "fmt"

func Whatever() {
	fmt.Println("Yeah whatever")
}

func RecursiveFibonacci(n int) int {
	if n <= 1 {
		return n
	} else {
		return RecursiveFibonacci(n - 1) + RecursiveFibonacci(n - 2)
	}
}