package main

import (
	"log"
	"microservice/auth/auth_pb"
	"microservice/auth/service"

	"net"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", "127.0.0.4:6666")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	auth_pb.RegisterAuthServiceServer(s, service.New())

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
