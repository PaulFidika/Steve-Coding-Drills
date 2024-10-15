package main

import "fmt"

func main() {
	myList := [6]int{1, 2, 3}
	fmt.Printf("My list: %v %v %v \n", myList, len(myList), cap(myList))

	mySlice := make([]int, 3, 6)

	mySlice = append(mySlice, 4, 5, 6)
	fmt.Printf("My slice %v", mySlice)

	mySlice = append(mySlice, 9)
	fmt.Printf("Capacity: %v", cap(mySlice))
}
