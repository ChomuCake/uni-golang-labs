package services

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/ChomuCake/uni-golang-labs/models"
)

type MockUserDBDetail struct {
	mockAddUser                      func(user models.User) error
	mockGetUserByUsernameAndPassword func(username, password string) (models.User, error)
	mockGetUserByUsername            func(username string) (models.User, error)
	mockGetUserByID                  func(userID int) (models.User, error)
}

func (m *MockUserDBDetail) AddUser(user models.User) error {
	if m.mockAddUser != nil {
		return m.mockAddUser(user)
	}
	return nil
}

func (m *MockUserDBDetail) GetUserByUsernameAndPassword(username, password string) (models.User, error) {
	if m.mockGetUserByUsernameAndPassword != nil {
		return m.mockGetUserByUsernameAndPassword(username, password)
	}
	return models.User{}, nil
}

func (m *MockUserDBDetail) GetUserByUsername(username string) (models.User, error) {
	if m.mockGetUserByUsername != nil {
		return m.mockGetUserByUsername(username)
	}
	return models.User{}, nil
}

func (m *MockUserDBDetail) GetUserByID(userID int) (models.User, error) {
	if m.mockGetUserByID != nil {
		return m.mockGetUserByID(userID)
	}
	return models.User{}, nil
}

func TestUserService_RegisterUser_Success(t *testing.T) {
	// Arrange
	MockUserDBDetail := &MockUserDBDetail{
		mockGetUserByUsername: func(username string) (models.User, error) {
			return models.User{}, sql.ErrNoRows
		},
		mockAddUser: func(user models.User) error {
			return nil
		},
	}
	s := NewUserService(MockUserDBDetail)

	// Act
	err := s.RegisterUser(testUser)

	// Assert
	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}
}

func TestUserService_RegisterUser_UserExists(t *testing.T) {
	// Arrange
	MockUserDBDetail := &MockUserDBDetail{
		mockGetUserByUsername: func(username string) (models.User, error) {
			return models.User{}, nil
		},
	}
	s := NewUserService(MockUserDBDetail)

	// Act
	err := s.RegisterUser(testUser)

	// Assert
	expectedError := "user with such name is already exists"
	if err == nil || err.Error() != expectedError {
		t.Errorf("Received incorrect error: received %v, expected %v", err, expectedError)
	}
}

func TestUserService_RegisterUser_AddUserError(t *testing.T) {
	// Arrange
	MockUserDBDetail := &MockUserDBDetail{
		mockGetUserByUsername: func(username string) (models.User, error) {
			return models.User{}, sql.ErrNoRows
		},
		mockAddUser: func(user models.User) error {
			return errors.New("registration failed")
		},
	}
	s := NewUserService(MockUserDBDetail)

	// Act
	err := s.RegisterUser(testUser)

	// Assert
	expectedError := "registration failed"
	if err == nil || err.Error() != expectedError {
		t.Errorf("Received incorrect error: received %v, expected %v", err, expectedError)
	}
}

func TestUserService_LoginUser_Success(t *testing.T) {
	// Arrange
	MockUserDBDetail := &MockUserDBDetail{
		mockGetUserByUsernameAndPassword: func(username, password string) (models.User, error) {
			return testUser, nil
		},
	}
	s := NewUserService(MockUserDBDetail)

	// Act
	user, err := s.LoginUser(testUser)

	// Assert
	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

	if user != testUser {
		t.Errorf("Received incorrect user")
	}
}

func TestUserService_LoginUser_UserNotFound(t *testing.T) {
	// Arrange
	MockUserDBDetail := &MockUserDBDetail{
		mockGetUserByUsernameAndPassword: func(username, password string) (models.User, error) {
			return models.User{}, sql.ErrNoRows
		},
	}
	s := NewUserService(MockUserDBDetail)

	// Act
	_, err := s.LoginUser(testUser)

	// Assert
	expectedError := "errNoRows"
	if err == nil || err.Error() != expectedError {
		t.Errorf("Received incorrect error: received %v, expected %v", err, expectedError)
	}
}
