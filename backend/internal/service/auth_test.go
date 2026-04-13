package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const testSecret = "test-jwt-secret-for-unit-tests"

func TestGenerateAndValidateToken(t *testing.T) {
	svc := &AuthService{jwtSecret: testSecret}

	userID := uuid.New()
	email := "test@example.com"

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	gotID, gotEmail, err := svc.ValidateToken(signed)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}
	if gotID != userID {
		t.Errorf("user_id = %v, want %v", gotID, userID)
	}
	if gotEmail != email {
		t.Errorf("email = %q, want %q", gotEmail, email)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	svc := &AuthService{jwtSecret: testSecret}

	claims := jwt.MapClaims{
		"user_id": uuid.New().String(),
		"email":   "test@example.com",
		"exp":     time.Now().Add(-time.Hour).Unix(),
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(testSecret))

	_, _, err := svc.ValidateToken(signed)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	svc := &AuthService{jwtSecret: testSecret}

	claims := jwt.MapClaims{
		"user_id": uuid.New().String(),
		"email":   "test@example.com",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("wrong-secret"))

	_, _, err := svc.ValidateToken(signed)
	if err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}

func TestValidateToken_InvalidFormat(t *testing.T) {
	svc := &AuthService{jwtSecret: testSecret}

	_, _, err := svc.ValidateToken("not-a-jwt-token")
	if err == nil {
		t.Fatal("expected error for invalid token format, got nil")
	}
}

func TestValidateToken_MissingUserID(t *testing.T) {
	svc := &AuthService{jwtSecret: testSecret}

	claims := jwt.MapClaims{
		"email": "test@example.com",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(testSecret))

	_, _, err := svc.ValidateToken(signed)
	if err == nil {
		t.Fatal("expected error for missing user_id claim, got nil")
	}
}
