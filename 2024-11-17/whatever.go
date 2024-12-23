package main

import (
	"fmt"
	"strconv"
	"strings"
)

type IntSlice []int

func (intSlice IntSlice) String() string {
	var strs []string

	for _, v := range intSlice {
		strs = append(strs, strconv.Itoa(v))
	}

	return "[" + strings.Join(strs, ", ") + "]"
}

type MyStruct struct {
	Names []string
	Ages []int
}

func (mine *MyStruct) Add(name string, age int) bool {
	mine.Ages = append(mine.Ages, age)
	mine.Names = append(mine.Names, name)

	return true
}

type NewView []int

func (view NewView) String() string {
	var strs []string

	for _, v := range view {
		strs = append(strs, strconv.Itoa(v))
	}

	return "[" + strings.Join(strs, ", ") + "]"
}

func (view *NewView) Shuffle() string {
	*view = append(*view, (*view)[len(*view) - 1])

	return "Okay"
}


type NewerView interface {
	String() string
	Shuffle() string
}

func whatever() {
	var whatever IntSlice = []int{1, 2, 3}
	var s fmt.Stringer = whatever
	var view NewView = []int{4, 5, 6}
	var view2 NewerView = &view

	view.Shuffle()
	view.Shuffle()

	fmt.Printf("Hello %v %v %v", whatever, s, view2)
}