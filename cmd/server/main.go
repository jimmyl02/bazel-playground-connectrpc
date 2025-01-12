package main

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/jimmyl02/bazel-playground-connectrpc/proto/testproto"
)

type TestSever struct{}

func (s *TestSever) SayHi(ctx context.Context, req *connect.Request[testproto.SayHiRequest]) (*connect.Response[testproto.SayHiResponse], error) {
	fmt.Println("received message", req)
	res := connect.NewResponse(&testproto.SayHiResponse{
		Response: "message received!",
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func main() {
	fmt.Println("beginning run")
	testserver := &TestSever{}
	testproto.NewTestHandler(testserver)
}
