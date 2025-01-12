package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/jimmyl02/bazel-playground-connectrpc/proto/testproto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
	path, handler := testproto.NewTestHandler(testserver)

	mux := http.NewServeMux()
	mux.Handle(path, handler)
	http.ListenAndServe(
		":8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
