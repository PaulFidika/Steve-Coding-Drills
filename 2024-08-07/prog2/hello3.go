package main

import "fmt"

func example() {
    defer fmt.Println("This prints last")
    defer fmt.Println("This prints second")
	fmt.Println("This prints first")
}