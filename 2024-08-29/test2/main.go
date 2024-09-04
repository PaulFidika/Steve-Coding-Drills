package main

import (
	"fmt"

	"github.com/aws/smithy-go/ptr"
)

func main() {
	str := ptr.String("Hello World")
	fmt.Println(str)
}