package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Alireza7997/go_microservices/auth_service/auth_pb"
	"github.com/Alireza7997/go_microservices/auth_service/global"
	_ "github.com/Alireza7997/go_microservices/auth_service/load"
	"github.com/Alireza7997/go_microservices/auth_service/service"

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
	auth_pb.RegisterAuthServiceServer(s, service.New())

	slog.Info("auth service listening", "addr", addr)
	if err := s.Serve(lis); err != nil {
		slog.Error("grpc server failed", "err", err)
	}
}
