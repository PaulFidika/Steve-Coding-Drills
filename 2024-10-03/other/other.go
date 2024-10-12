package other

import (
	"fmt"
	"strings"
)

func Other() {
	fmt.Println("whatever", "hello")
}

func Greet(name string) string {
	return fmt.Sprintf("Hello there, %s\n!", name)
}

func Multi(names []string) string {
	if len(names) < 1 {
		names = []string{"Bob"}
	}

	return fmt.Sprintf("Hey guys %s", strings.Join(names, ", "))
}