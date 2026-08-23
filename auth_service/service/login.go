package service

import (
	"context"
	"database/sql"
	"microservice/auth_service/auth_pb"
	"microservice/auth_service/global"
	"microservice/general"
	"microservice/pkg/database"

	"github.com/kataras/iris/v12"
)

func (s *service) Login(ctx context.Context, req *auth_pb.LoginRequest) (*auth_pb.LoginResponse, error) {
	var (
		username = req.GetUserName()
		password = req.GetPassword()
	)

	response := &auth_pb.LoginResponse{}

	db := database.Connect(global.DB)
	defer db.Close()

	user, err := s.GetUserByUsername(db, username)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Err = &general.Error{
				Code:    iris.StatusUnauthorized,
				Message: "invalid username or password",
			}
			return response, nil
		}
		return nil, err
	}

	if !s.MatchPasswords(user.Password, password) {
		response.Err = &general.Error{
			Code:    iris.StatusUnauthorized,
			Message: "invalid username or password",
		}
		return response, nil
	}

	token, err := s.GenerateJWT(user.Id)
	if err != nil {
		return nil, err
	}

	response.UserID = user.Id
	response.Token = token

	return response, nil
}
