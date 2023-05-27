package util

import (
	"time"

	"github.com/ChomuCake/uni-golang-labs/models"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

type JWTTokenManager struct{}

var secretKey = []byte("fd9f5dc52a0b5728c5182c593e0fae7d821e6c7a0fe64b78e67450a0a6860d63")

func (tm *JWTTokenManager) GenerateToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Токен дійсний протягом 24 годин
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (tm *JWTTokenManager) VerifyToken(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}

func (tm *JWTTokenManager) ExtractUserIDFromToken(token interface{}) (int, error) {
	parsedToken, ok := token.(*jwt.Token)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	return int(userID), nil
}
