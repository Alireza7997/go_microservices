package service

import (
	"context"
	"database/sql"
	"microservice/auth/auth_pb"
	"microservice/auth/models"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	auth_pb.UnimplementedAuthServiceServer
}

var serviceInstance = service{}

func New() auth_pb.AuthServiceServer {
	return &serviceInstance
}

func (s service) CreateUser(ctx context.Context, db *sql.DB, username, password, email string) (int64, error) {
	result, err := db.Exec("INSERT INTO ? (username, password, email) VALUES (?, ?, ?)", models.UserTable, username, password, email)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (service) UserExists(db *sql.DB, username string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM ? WHERE username = ?", models.UserTable, username).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, err
	}

	return false, err
}

func (service) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}
