package main

import (
	"fmt"
	"log/slog"
	"net"

	"microservice/chat_service/chat_pb"
	"microservice/chat_service/global"
	_ "microservice/chat_service/load"
	"microservice/chat_service/service"

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
	chat_pb.RegisterChatServiceServer(s, service.New())

	slog.Info("chat service listening", "addr", addr)
	if err := s.Serve(lis); err != nil {
		slog.Error("grpc server failed", "err", err)
	}
}
