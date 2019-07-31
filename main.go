package main

import (
	"github.com/ravenhurst/golang-playground/foo"
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	fmt.Printf("foo => %s", foo.GetFoo())
}
