package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	sourceStr := "Hello, OTUS!"
	fmt.Println(reverse.String(sourceStr))
}
