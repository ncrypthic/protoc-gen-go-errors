package main

import (
	"context"
	"errors"
	"net"

	test "github.com/ncrypthic/protoc-gen-go-errors/example/grpc/test"
	"google.golang.org/grpc"
)

type server struct {
	test.UnimplementedUserServiceServer
}

func (s *server) Register(context.Context, *test.RegisterRequest) (*test.RegisterResponse, error) {
	return nil, test.ERR_err_invalid_email
}

func main() {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	test.RegisterUserServiceServer(s, &server{})
	go s.Serve(l)

	c, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := test.NewUserServiceClient(c)
	_, err = client.Register(context.TODO(), &test.RegisterRequest{})
	if !errors.Is(err, test.ERR_err_invalid_email) {
		panic("unexpected error")
	}
	if errors.Is(err, test.ERR_err_email_exists) {
		panic("unexpected error")
	}
}
