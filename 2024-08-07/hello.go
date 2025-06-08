package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")

	dumbMap := map[string]uint8{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	for key, value := range dumbMap {
		fmt.Println(key, value)
	}

	fmt.Println(len(dumbMap))

	myArray := [5]uint8{5, 5, 5}

	fmt.Println(len(myArray), cap(myArray))

	for _, item := range myArray {
		fmt.Println(item)
	}

	anotherSlice := make([]uint8, 3, 5)

	fmt.Println(len(anotherSlice), cap(anotherSlice))
	fmt.Println(&anotherSlice[2])

	anotherSlice = append(anotherSlice, 5, 6)

	fmt.Println(len(anotherSlice), cap(anotherSlice))
	fmt.Println(&anotherSlice[2])

	myArray2 := [2]uint8{1, 2}
	fmt.Printf("%p\n", &myArray2)

	secondary(myArray2)
	fmt.Println(myArray2)

	// Slice
	slice := []uint8{1, 2, 3}

	// Arrays
	array1 := [3]uint8{1, 2, 3}
	array2 := [...]uint8{1, 2, 3}
} 

func secondary(myArray [2]uint8) {
	myArray[0] = 10
	fmt.Println(myArray)
}
