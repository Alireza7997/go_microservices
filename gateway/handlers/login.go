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

func login(w http.ResponseWriter, r *http.Request) {
	req := &dto.LoginRequest{}

	utils.ParseBody(r.Body, req)
	if req.Username == "" || req.Password == "" {
		panic(errors.New(http.StatusBadRequest, "missing field", "username and password are required"))
	}

	var res *auth_pb.LoginResponse
	calls.WithAuthService(func(client auth_pb.AuthServiceClient) {
		resp, err := client.Login(r.Context(), &auth_pb.LoginRequest{
			UserName: req.Username,
			Password: req.Password,
		})
		if err != nil {
			panic(errors.New(http.StatusServiceUnavailable, err.Error(), "auth service unavailable"))
		}
		if resp.GetErr() != nil {
			panic(errors.New(resp.GetErr().GetCode(), resp.GetErr().GetMessage(), resp.GetErr().GetMessage()))
		}
		res = resp
	})

	utils.WriteJSON(w, http.StatusOK, dto.LoginResponse{
		UserID:   res.GetUserID(),
		Token:    res.GetToken(),
		Username: req.Username,
	})
}

var Login = g.Handler{Handler: login}
