package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println(os.Args[1])
		switch os.Args[1] {
		case "coordinates":
			fmt.Println("Running coordinates.go")
			runCoordinates()
		case "whatever":
			fmt.Println("Running whatever.go")
			whatever()
		case "byteCounter":
			fmt.Println("Running the byte counter")
			byteCounter()
		default:
			fmt.Println("Unknown file")
		}

	} else {
		fmt.Println("Please enter which file you want to run")
	}
}