package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		string_slice := strings.Split(line, "world")
		str := strings.Join(string_slice, "Earth")

		fmt.Println(str)
	}

	var a []int

	a = append(a, 1)
	fmt.Println(a)
}
