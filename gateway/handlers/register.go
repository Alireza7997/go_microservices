package handlers

import (
	"net/http"

	"microservice/auth_service/auth_pb"
	"microservice/gateway/calls"
	"microservice/gateway/dto"
	g "microservice/gateway/global"
	"microservice/gateway/handlers/utils"
	"microservice/pkg/errors"
)

func register(w http.ResponseWriter, r *http.Request) {
	req := &dto.RegisterRequest{}

	utils.ParseBody(r.Body, req)
	if req.Username == "" || req.Password == "" || req.PasswordConfirm == "" || req.Email == "" {
		panic(errors.New(http.StatusBadRequest, "missing field", "username, password, password_confirm and email are required"))
	}

	calls.WithAuthService(func(client auth_pb.AuthServiceClient) {
		resp, err := client.Register(r.Context(), &auth_pb.RegisterRequest{
			UserName:        req.Username,
			Password:        req.Password,
			ConfirmPassword: req.PasswordConfirm,
			Email:           req.Email,
		})
		if err != nil {
			panic(errors.New(http.StatusServiceUnavailable, err.Error(), "auth service unavailable"))
		}
		if resp.GetErr() != nil {
			panic(errors.New(resp.GetErr().GetCode(), resp.GetErr().GetMessage(), resp.GetErr().GetMessage()))
		}
	})

	w.WriteHeader(http.StatusCreated)
}

var Register = g.Handler{Handler: register}
