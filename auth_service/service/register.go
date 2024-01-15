package service

import (
	"context"
	"microservice/auth/auth_pb"
	"microservice/auth/global"
	"microservice/general"
	"microservice/pkg/database"

	"github.com/kataras/iris/v12"
)

func (s *service) Register(ctx context.Context, req *auth_pb.RegisterRequest) (*auth_pb.RegisterResponse, error) {
	// Create database connnection and closing it when we are finished using it
	db := database.Connect(global.DB)
	defer db.Close()

	// Getting request values
	var (
		username        = req.GetUserName()
		password        = req.GetPassword()
		passwordConfirm = req.GetConfirmPassword()
		email           = req.GetEmail()
	)

	// Create the response
	response := &auth_pb.RegisterResponse{}

	// Check if passwords match
	if passwordConfirm != password {
		response.Err = &general.Error{
			Code:    iris.StatusBadRequest,
			ErrMsg:  "",
			Message: "passwords should be the same",
		}
		return response, nil
	}

	// Hash the password
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Check if the username is already taken
	exists, err := s.UserExists(db, username)
	if err != nil {
		return nil, err
	}

	if exists {
		response.Err = &general.Error{
			Code:    iris.StatusBadRequest,
			ErrMsg:  "",
			Message: "username already taken",
		}
		return response, nil
	}

	// Create the user
	userID, err := s.CreateUser(ctx, db, username, hashedPassword, email)
	if err != nil {
		return nil, err
	}

	response.UserID = userID

	return response, err
}
