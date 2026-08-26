package service

import (
	"context"
	"github.com/Alireza7997/go_microservices/auth_service/auth_pb"
	"github.com/Alireza7997/go_microservices/auth_service/global"
	"github.com/Alireza7997/go_microservices/general"
	"github.com/Alireza7997/go_microservices/pkg/database"

	"github.com/kataras/iris/v12"
)

func (s *service) Register(ctx context.Context, req *auth_pb.RegisterRequest) (*auth_pb.RegisterResponse, error) {
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
			Message: "passwords should be the same",
		}
		return response, nil
	}

	// Create database connection and close it when we are finished using it
	db := database.Connect(global.DB)
	defer db.Close()

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
