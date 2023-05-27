package util

import "github.com/ChomuCake/uni-golang-labs/models"

type TokenManager interface {
	VerifyToken(tokenString string) (interface{}, error)
	GenerateToken(user models.User) (string, error)
	ExtractUserIDFromToken(interface{}) (int, error)
}
