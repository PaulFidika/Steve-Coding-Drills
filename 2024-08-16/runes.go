package main

import (
	"fmt"
	"strings"
)

type gasEngine struct {
	mpg uint8
	gallons uint8
	owner owner
}

type owner struct {
	name string
}

func main() {
	var myString = "维基百科"
	var myRune = []rune(myString)

	fmt.Println(len(myString))
	fmt.Println(myString[0])

	fmt.Println(len(myRune))
	fmt.Println(myRune[0])

	var strBuilder = strings.Builder{}
	strBuilder.WriteString(myString)
	fmt.Println(strBuilder.String())

	var myEngine = gasEngine{ 
		mpg: 24, 
		gallons: 30,
		owner: owner{"Bobby"},
	}
	fmt.Println(myEngine.mpg, myEngine.gallons, myEngine.owner.name)
}