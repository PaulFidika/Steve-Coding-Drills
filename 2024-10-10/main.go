package main

import (
	"fmt"
)

func main() {
	// if (len(os.Args) > 1) {
	// 	SaySomething(os.Args[1])
	// } else {
	// 	SaySomething("nothing provided")
	// }

	PrintTypes()
}

func SaySomething(message string) string {
	value := fmt.Sprintf("Here it is %s", message)
	fmt.Print(value)

	return value
}

func PrintTypes() {
	var a map[string]int
	fmt.Printf("Here: %T %v \n", a, a)
	a = map[string]int{"abc": 123, "xyz": 456}
	fmt.Printf("Here: %T %v \n", a, a)
}