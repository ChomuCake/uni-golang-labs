package database

import "github.com/ChomuCake/uni-golang-labs/models"

// UserDB визначає інтерфейс для роботи з даними юзерів
type UserDB interface {
	AddUser(user models.User) error
	GetUserByUsernameAndPassword(username, password string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserByID(userID int) (models.User, error)
}
