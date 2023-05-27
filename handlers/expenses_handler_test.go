package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ChomuCake/uni-golang-labs/models"
	"github.com/golang-jwt/jwt"
)

// MockExpenseDB є замінником реалізації ExpenseDB
type MockExpenseDB struct{}

func (db *MockExpenseDB) AddExpense(expense models.Expense) error {
	return nil
}

func (db *MockExpenseDB) GetUserExpenses(userID int) ([]models.Expense, error) {
	return []models.Expense{
		{ID: 1, Amount: 10.0, Date: time.Now(), UserID: 1},
		{ID: 2, Amount: 20.0, Date: time.Now(), UserID: 1},
	}, nil
}

func (db *MockExpenseDB) UpdateUserExpenses(expense models.Expense) error {
	return nil
}

func (db *MockExpenseDB) DeleteExpense(expenseID string) error {
	return nil
}

// MockUserDB є замінником реалізації UserDB
type MockUserDB struct{}

func (db *MockUserDB) GetUserByID(userID int) (models.User, error) {
	return models.User{ID: 1, Username: "John Doe"}, nil
}

func (db *MockUserDB) AddUser(user models.User) error {
	return nil
}

func (db *MockUserDB) GetUserByUsername(username string) (models.User, error) {
	return models.User{ID: 1, Username: "John Doe"}, nil
}

func (db *MockUserDB) GetUserByUsernameAndPassword(username, password string) (models.User, error) {
	return models.User{ID: 1, Username: "John Doe"}, nil
}

// MockTokenManager є замінником реалізації TokenManager
type MockTokenManager struct{}

func (tm *MockTokenManager) ExtractUserIDFromToken(token interface{}) (int, error) {
	claims, ok := token.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	return int(userID), nil
}

func (tm *MockTokenManager) GenerateToken(user models.User) (string, error) {
	return "token", nil
}

func (tm *MockTokenManager) VerifyToken(tokenString string) (interface{}, error) {
	return "token", nil
}

func (tm *MockTokenManager) ExtractToken(r *http.Request) string {
	return "token"
}

func (tm *MockTokenManager) ExtractUserIDFromRequest(r *http.Request) (int, error) {
	return 1, nil
}

func SetUpHandlerDep() *ExpenseHandler {
	h := &ExpenseHandler{
		ExpenseDB: &MockExpenseDB{},
		UserDB:    &MockUserDB{},
		TokenMng:  &MockTokenManager{},
	}
	return h
}

func TestExpensesHandler_PostExpense(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"amount": 10}`)
	req, err := http.NewRequest("POST", "/expenses", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusCreated)
	}
}
