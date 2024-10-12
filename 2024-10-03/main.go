package main

import (
	"fmt"
	"os"
	"test/other"
)

func main() {
	other.Other()

	var name string
	if (len(os.Args) > 1) {
		name = os.Args[1]
	} else {
		name = "Guest"
	}

	fmt.Printf("Hello there, %s\n", name)
}