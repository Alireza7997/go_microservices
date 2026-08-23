package service

import (
	"context"
	"fmt"
	"time"

	"microservice/greet_service/greet_pb"
)

type service struct {
	greet_pb.UnimplementedGreetServiceServer
}

func New() greet_pb.GreetServiceServer {
	return &service{}
}

func (service) Ping(ctx context.Context, req *greet_pb.PingRequest) (*greet_pb.PingResponse, error) {
	name := req.GetName()
	if name == "" {
		name = "stranger"
	}

	return &greet_pb.PingResponse{
		Greeting:  fmt.Sprintf("hello, %s!", name),
		ServerTime: time.Now().Unix(),
	}, nil
}
