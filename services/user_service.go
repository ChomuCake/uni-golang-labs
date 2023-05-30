package services

import (
	"database/sql"
	"errors"

	"github.com/ChomuCake/uni-golang-labs/models"
)

type detailUserDB interface {
	AddUser(user models.User) error
	GetUserByUsernameAndPassword(username, password string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserByID(userID int) (models.User, error)
}

type UserService struct {
	userDB detailUserDB
}

func NewUserService(userDB detailUserDB) *UserService {
	return &UserService{userDB}
}

func (s *UserService) RegisterUser(user models.User) error {

	_, err := s.userDB.GetUserByUsername(user.Username)
	if err == nil {
		return errors.New("user with such name is already exists")
	}

	err = s.userDB.AddUser(user)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("errNoRows")
		}
		return errors.New("registration failed")
	}

	return nil
}

func (s *UserService) LoginUser(user models.User) (models.User, error) {

	existingUser, err := s.userDB.GetUserByUsernameAndPassword(user.Username, user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("errNoRows")
		}
		return models.User{}, errors.New("user with that name isn't exists")
	}

	return existingUser, nil
}
