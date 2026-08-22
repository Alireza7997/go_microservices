package service

import (
	"context"
	"database/sql"
	"fmt"
	"microservice/auth_service/auth_pb"
	"microservice/auth_service/global"
	"microservice/auth_service/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const tokenLifeSpan = 24 * time.Hour

type service struct {
	auth_pb.UnimplementedAuthServiceServer
}

var serviceInstance = service{}

func New() auth_pb.AuthServiceServer {
	return &serviceInstance
}

func (s service) CreateUser(ctx context.Context, db *sql.DB, username, password, email string) (int64, error) {
	var userID int64
	err := db.QueryRowContext(
		ctx,
		fmt.Sprintf("INSERT INTO %s (username, password, email) VALUES ($1, $2, $3) RETURNING id", models.UserTable),
		username, password, email,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (service) UserExists(db *sql.DB, username string) (bool, error) {
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = $1", models.UserTable), username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (service) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

func (service) MatchPasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (service) GenerateJWT(id int64) (string, error) {
	claims := &models.Claims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifeSpan)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(global.SecretKeyBytes)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (service) GetUser(db *sql.DB, id int64) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(
		fmt.Sprintf("SELECT id, username, email FROM %s WHERE id = $1", models.UserTable),
		id,
	).Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

