package main

import (
	"fmt"

	"github.com/jimmyl02/bazel-playground-connectrpc/proto/testproto"
)

type TestSever struct{}

func (s *TestSever) Serve() {
	fmt.Println("server is running")

	a := testproto.SayHiRequest{}
	fmt.Println("a", a.Name)
}

func main() {
	fmt.Println("asdf")
}
