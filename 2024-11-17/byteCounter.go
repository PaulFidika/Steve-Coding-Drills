package main

import (
	"fmt"
	"io"
	"os"
)

type ByteCounter int

func (counter *ByteCounter) Write(p []byte) (int, error) {
	l := len(p)
	fmt.Println("I tried")
	*counter += ByteCounter(l)
	return len(p), nil
}

func byteCounter() {
	var c ByteCounter

	f1, _ := os.Open("a.txt")
	f2 := &c

	_, _ = io.Copy(f2, f1)
	f1.Seek(0, 0)
	_, _ = io.Copy(f2, f1)
	f1.Seek(0, 0)
	_, _ = io.Copy(f2, f1)

	fmt.Println(*f2)

	// fmt.Println("copied", n, "bytes")
	// fmt.Println(c)
}