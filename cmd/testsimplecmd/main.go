package main

import (
	"fmt"

	"github.com/moznion/go-optional"
)

func main() {
	fmt.Println("hello world!")

	some := optional.Some(true)
	fmt.Println(some.Unwrap())
}
