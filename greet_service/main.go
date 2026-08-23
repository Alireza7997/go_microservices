package main

import (
	"fmt"
	"log/slog"
	"net"

	"microservice/greet_service/greet_pb"
	"microservice/greet_service/global"
	_ "microservice/greet_service/load"
	"microservice/greet_service/service"

	"google.golang.org/grpc"
)

func main() {
	addr := fmt.Sprintf("%s:%s", global.CFG.CurrentMicroservice.IP, global.CFG.CurrentMicroservice.Port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("failed to listen", "addr", addr, "err", err)
		return
	}

	s := grpc.NewServer()
	greet_pb.RegisterGreetServiceServer(s, service.New())

	slog.Info("greet service listening", "addr", addr)
	if err := s.Serve(lis); err != nil {
		slog.Error("grpc server failed", "err", err)
	}
}
