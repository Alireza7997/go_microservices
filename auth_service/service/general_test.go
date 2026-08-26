package service

import (
	"testing"
	"time"

	"github.com/Alireza7997/go_microservices/auth_service/global"
	"github.com/Alireza7997/go_microservices/auth_service/models"

	"github.com/golang-jwt/jwt/v5"
)

func TestHashPasswordAndMatchPasswords(t *testing.T) {
	s := service{}

	hash, err := s.HashPassword("s3cret-password")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "s3cret-password" {
		t.Fatal("password was not hashed")
	}

	if !s.MatchPasswords(hash, "s3cret-password") {
		t.Error("expected matching password to validate")
	}
	if s.MatchPasswords(hash, "wrong-password") {
		t.Error("expected wrong password to fail validation")
	}
}

func TestGenerateJWT(t *testing.T) {
	global.SecretKeyBytes = []byte("test-secret-key")

	tokenStr, err := service{}.GenerateJWT(1234)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}

	claims := &models.Claims{}
	parsed, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("expected valid token")
	}
	if claims.UserId != 1234 {
		t.Fatalf("expected user id 1234, got %d", claims.UserId)
	}

	expiry := time.Until(claims.ExpiresAt.Time)
	if expiry <= 0 || expiry > 25*time.Hour {
		t.Fatalf("token expiry not set correctly: %v", expiry)
	}
}

func TestGenerateJWTCannotBeTampered(t *testing.T) {
	global.SecretKeyBytes = []byte("test-secret-key")

	tokenStr, _ := service{}.GenerateJWT(1)

	_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("some-other-secret"), nil
	})
	if err == nil {
		t.Fatal("expected token signed with a different key to be rejected")
	}
}
