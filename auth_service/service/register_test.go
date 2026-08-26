package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Alireza7997/go_microservices/auth_service/auth_pb"
	"github.com/Alireza7997/go_microservices/auth_service/global"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/protobuf/proto"
)

func setupTestDB(t *testing.T) (*service, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	global.SecretKeyBytes = []byte("test-secret-key")
	global.DB = func() (*sql.DB, error) { return db, nil }
	return &service{}, mock
}

func TestRegisterPasswordsDoNotMatch(t *testing.T) {
	s, _ := setupTestDB(t)

	res, err := s.Register(context.Background(), &auth_pb.RegisterRequest{
		UserName:        "alice",
		Password:        "pw1",
		ConfirmPassword: "pw2",
		Email:           "alice@example.com",
	})
	if err != nil {
		t.Fatalf("Register returned transport error: %v", err)
	}
	if res.Err == nil || res.Err.Message != "passwords should be the same" {
		t.Fatalf("expected password mismatch error, got %+v", res.Err)
	}
}

func TestRegisterUsernameTaken(t *testing.T) {
	s, mock := setupTestDB(t)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users WHERE username = \$1`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	res, err := s.Register(context.Background(), &auth_pb.RegisterRequest{
		UserName:        "alice",
		Password:        "pw",
		ConfirmPassword: "pw",
		Email:           "alice@example.com",
	})
	if err != nil {
		t.Fatalf("Register returned transport error: %v", err)
	}
	if res.Err == nil || res.Err.Message != "username already taken" {
		t.Fatalf("expected username taken error, got %+v", res.Err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

// Integration-style test of the full register flow:
// duplicate check -> insert -> RETURNING id.
func TestRegisterFullFlow(t *testing.T) {
	s, mock := setupTestDB(t)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users WHERE username = \$1`).
		WithArgs("bob").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery(`INSERT INTO users \(username, password, email\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs("bob", sqlmock.AnyArg(), "bob@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(7))

	res, err := s.Register(context.Background(), &auth_pb.RegisterRequest{
		UserName:        "bob",
		Password:        "pw",
		ConfirmPassword: "pw",
		Email:           "bob@example.com",
	})
	if err != nil {
		t.Fatalf("Register returned transport error: %v", err)
	}
	if res.Err != nil {
		t.Fatalf("unexpected error in response: %+v", res.Err)
	}
	if !proto.Equal(res, &auth_pb.RegisterResponse{UserID: 7}) {
		t.Fatalf("expected UserID=7, got %+v", res)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users WHERE username = \$1`).
		WithArgs("alice").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	exists, err := (service{}).UserExists(db, "alice")
	if err != nil {
		t.Fatalf("UserExists failed: %v", err)
	}
	if !exists {
		t.Fatal("expected exists=true")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}
