package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TokenDefaultExpiryDate() time.Time {
	return time.Now().UTC().Add(time.Hour * 24)
}

func CheckToken(tokenString string, secret string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

func GetUserIDFromToken(tokenString string, secret string) (uint64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrInvalidKeyType
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, jwt.ErrInvalidKeyType
	}

	return uint64(userID), nil
}

func GenerateToken(expiryDate time.Time, userID uint64, secret string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": expiryDate.Unix(),
			"iat": time.Now().UTC().Unix(),
			"nbt": time.Now().UTC().Unix(),
			"sub": userID,
			"iss": "riley",
		},
	)

	return token.SignedString([]byte(secret))
}
