package auth

import (
	"testing"
	"time"

	"riley/internal/config"
	"riley/internal/models"
	"riley/internal/sql"

	"github.com/golang-jwt/jwt/v5"
)

func TestCheckToken(t *testing.T) {
	wrongToken := "0123456789"
	secret := "secret"
	expiryDate := time.Now().Add(time.Hour * 24)

	if err := CheckToken(wrongToken, secret); err == nil {
		t.Error("Testing check token with wrong token: Wanted err, got nil")
	}

	db := sql.Connect(config.LoadTestConfig())

	user, err := models.UserCreate("exampleTestCheckToken@example.com", "password123%A%", db)
	if err != nil {
		t.Error("Testing check token: Wanted nil, got", err)
	}

	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Error("Delete user during check token test failed: Wanted nil, got", err)
		}
	}()

	correctToken, err := GenerateToken(expiryDate, user.ID, secret)
	if err != nil {
		t.Error("Testing generate token: Wanted nil, got", err)
	}

	if err := CheckToken(correctToken, secret); err != nil {
		t.Error("Testing check token with correct token: Wanted nil, got", err)
	}

	if err := CheckToken(correctToken, "wrong secret"); err == nil {
		t.Error("Testing check token with wrong secret: Wanted error, got nil")
	}

	if err := CheckToken("", secret); err == nil {
		t.Error("Testing check token with empty token: Wanted error, got nil")
	}

	if err := CheckToken(correctToken, ""); err == nil {
		t.Error("Testing check token with empty secret: Wanted error, got nil")
	}

	if err := CheckToken("", ""); err == nil {
		t.Error("Testing check token with empty token and secret: Wanted error, got nil")
	}
}

func TestGenerateToken(t *testing.T) {
	secret := "secret"
	expiryDate := time.Now().Add(time.Hour * 24)

	db := sql.Connect(config.LoadTestConfig())

	user, err := models.UserCreate("exampleTestGenerateToken@example.com", "password123%A%", db)
	if err != nil {
		t.Error("Testing generate token: Wanted nil, got", err)
	}

	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Error("Delete user during generate token test failed: Wanted nil, got", err)
		}
	}()

	token, err := GenerateToken(expiryDate, user.ID, secret)
	if err != nil {
		t.Error("Testing generate token: Wanted nil, got", err)
	}

	if token == "" {
		t.Error("Testing generate token: Wanted token, got empty string")
	}

	tokenUnwrapped, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Error("Testing generate token: Wanted nil, got", err)
	}

	if uint64(tokenUnwrapped.Claims.(jwt.MapClaims)["sub"].(float64)) != user.ID {
		t.Error("Testing generate token: Wanted", user.ID, "got", tokenUnwrapped.Claims.(jwt.MapClaims)["sub"])
	}

	if tokenUnwrapped.Claims.(jwt.MapClaims)["iss"] != "riley" {
		t.Error("Testing generate token: Wanted riley got", tokenUnwrapped.Claims.(jwt.MapClaims)["iss"])
	}

	tokenExpiryDate := tokenUnwrapped.Claims.(jwt.MapClaims)["exp"].(float64)
	expiryDateToCompare := float64(expiryDate.Unix())

	if tokenExpiryDate != expiryDateToCompare {
		t.Error("Testing generate token: Wanted", expiryDate.Unix(), "got", tokenUnwrapped.Claims.(jwt.MapClaims)["exp"])
	}
}
