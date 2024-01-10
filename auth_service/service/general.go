package service

import "microservice/auth/auth_pb"

type service struct {
	auth_pb.UnimplementedAuthServiceServer
}

var serviceInstance = service{}

func New() auth_pb.AuthServiceServer {
	return &serviceInstance
}

// func (*service) Register()
