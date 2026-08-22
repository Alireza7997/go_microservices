package handlers

import (
	"microservice/auth_service/auth_pb"
	"microservice/gateway/calls"
	"microservice/gateway/dto"
	g "microservice/gateway/global"
	"microservice/gateway/handlers/utils"
	"net/http"
)

func register(w http.ResponseWriter, r *http.Request) {
	req := &dto.RegisterRequest{}
	ctx := r.Context()

	utils.ParseBody(r.Body, req)

	s := calls.NewAuthService()
	s.Call(func(service auth_pb.AuthServiceClient) {
		grpcRes, err := service.Register(ctx, &auth_pb.RegisterRequest{
			UserName:        req.Username,
			Password:        req.Password,
			ConfirmPassword: req.PasswordConfirm,
			Email:           req.Email,
		})
		if grpcRes != nil {
			s.Check(grpcRes.Err, nil)
		} else {
			s.Check(nil, err)
		}
	})

	w.WriteHeader(http.StatusCreated)
}

var Register = g.Handler{
	Handler: register,
}
