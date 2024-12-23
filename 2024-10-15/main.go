package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	scan := bufio.NewScanner(os.Stdin)
	words := make(map[string]int)

	scan.Split(bufio.ScanWords)

	for scan.Scan() {
		words[scan.Text()]++
	}

	fmt.Println(len(words), "unique words")

	type kv struct {
		key string
		value int
	}

	var ss []kv

	for key, value := range words {
		ss = append(ss, kv{key, value})
	}

	sort.Slice(ss, func(i int, j int) bool {
		return ss[i].value > ss[j].value
	})

	for _, slice := range ss[:3] {
		fmt.Println(slice.key, ": ", slice.value)
	}

	num, _ := returnSomething()
}

func returnSomething() (int, error) {
	return 5, nil
}

var monkey = func(name string) {
	fmt.Printf("His name is %v\n", name)
}