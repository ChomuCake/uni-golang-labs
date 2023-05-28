package handlers

import (
	"bytes"
	"database/sql" // only for sql.ErrNoRows
	"encoding/json"
	"errors"
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
	if expense.RawDate == "err" {
		return errors.New("server error")
	}
	return nil
}

func (db *MockExpenseDB) GetUserExpenses(userID int) ([]models.Expense, error) {
	if userID == 3 {
		return nil, errors.New("server error")
	}
	return []models.Expense{
		{ID: 1, Amount: 10, Date: fixedTime, Category: "test", UserID: 1},
		{ID: 2, Amount: 20, Date: fixedTime, Category: "test", UserID: 1},                    //day
		{ID: 3, Amount: 20, Date: fixedTime.AddDate(0, 0, -1), Category: "test", UserID: 1},  //month
		{ID: 4, Amount: 20, Date: fixedTime.AddDate(0, 0, -32), Category: "test", UserID: 1}, // all
	}, nil
}

func (db *MockExpenseDB) UpdateUserExpenses(expense models.Expense) error {
	if expense.Amount == -1 {
		return errors.New("server error")
	}
	return nil
}

func (db *MockExpenseDB) DeleteExpense(expenseID string) error {
	if expenseID == "99" {
		return errors.New("not found")
	}
	return nil
}

// MockUserDB є замінником реалізації UserDB
type MockUserDB struct{}

func (db *MockUserDB) GetUserByID(userID int) (models.User, error) {
	if userID == 1 {
		return models.User{ID: 1, Username: "John Doe"}, nil
	}
	if userID == 3 {
		return models.User{ID: 3, Username: "Joe Doe"}, nil
	}
	return models.User{}, errors.New("server error")

}

func (db *MockUserDB) AddUser(user models.User) error {
	if user.Username == "ErrNoRows" {
		return sql.ErrNoRows
	}

	if user.Username == "ServerError" {
		return errors.New("server error")
	}
	return nil
}

func (db *MockUserDB) GetUserByUsername(username string) (models.User, error) {
	if username == "ErrNoRows" {
		return models.User{}, errors.New("server error")
	}

	if username == "ServerError" {
		return models.User{}, errors.New("server error")
	}

	if username == "Reg" {
		return models.User{}, errors.New("server error")
	}

	return models.User{ID: 1, Username: "John Doe"}, nil
}

func (db *MockUserDB) GetUserByUsernameAndPassword(username, password string) (models.User, error) {
	if username == "ErrNoRows" {
		return models.User{}, sql.ErrNoRows
	}

	if username == "ServerError" {
		return models.User{}, errors.New("server error")

	}

	if username == "Incorrect" {
		return models.User{ID: -1, Username: "Incorrect"}, nil
	}

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
	if user.Username == "Incorrect" {
		return "", jwt.ErrInvalidKey
	}

	return "token", nil
}

func (tm *MockTokenManager) VerifyToken(tokenString string) (interface{}, error) {
	return "token", nil
}

func (tm *MockTokenManager) ExtractToken(r *http.Request) string {
	return "token"
}

func (tm *MockTokenManager) ExtractUserIDFromRequest(r *http.Request) (int, error) {
	token := r.Header.Get("Token")

	if token == "Incorrect" {
		return 0, jwt.ErrInvalidKey
	}

	if token == "TokenWithoutUserInDB" {
		return 2, nil
	}

	if token == "TokenWithID3InDB" {
		return 3, nil
	}

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

var fixedTime time.Time

func SetTimeNow() {
	fixedTime = time.Now()
}

// ---------------- POST TESTS --------------------
func TestExpensesHandler_PostExpense(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"amount": 10}`)
	req, err := http.NewRequest("POST", "/expenses", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Correct")

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

func TestExpensesHandler_PostExpense_IncorrectTokenInRequest(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"amount": 10}`)
	req, err := http.NewRequest("POST", "/expenses", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Incorrect")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusUnauthorized)
	}
}

func TestExpensesHandler_PostExpense_IncorrectBodyRequest(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"amount": "invalid", "date": "2023-05-27", "user_id": 1}`)
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
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusBadRequest)
	}
}

func TestExpensesHandler_PostExpense_IncorrectUserIdInRequest(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"amount": 10}`)
	req, err := http.NewRequest("POST", "/expenses", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "TokenWithoutUserInDB")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusUnauthorized)
	}
}

func TestExpensesHandler_PostExpense_ServerError(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"rawdate": "err"}`)
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
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusInternalServerError)
	}
}

// -------------- END POST TESTS --------------

// -------------- GET TESTS --------------
func TestExpensesHandler_GetExpensesDay(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/expenses?sort=day", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "Correct")
	SetTimeNow()

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusOK)
	}

	var expenses []models.Expense
	err = json.Unmarshal(rr.Body.Bytes(), &expenses)
	if err != nil {
		t.Fatal(err)
	}

	if len(expenses) != 2 {
		t.Errorf("Отримано некоректну кількість витрат: отримано %d, очікувалося %d",
			len(expenses), 2)
	}
}

func TestExpensesHandler_GetExpensesMonth(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/expenses?sort=month", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "Correct")
	SetTimeNow()

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusOK)
	}

	var expenses []models.Expense
	err = json.Unmarshal(rr.Body.Bytes(), &expenses)
	if err != nil {
		t.Fatal(err)
	}

	if len(expenses) != 3 {
		t.Errorf("Отримано некоректну кількість витрат: отримано %d, очікувалося %d",
			len(expenses), 3)
	}
}

func TestExpensesHandler_GetExpensesAll(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/expenses?sort=all", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "Correct")
	SetTimeNow()

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusOK)
	}

	var expenses []models.Expense
	err = json.Unmarshal(rr.Body.Bytes(), &expenses)
	if err != nil {
		t.Fatal(err)
	}

	if len(expenses) != 4 {
		t.Errorf("Отримано некоректну кількість витрат: отримано %d, очікувалося %d",
			len(expenses), 4)
	}
}

func TestExpensesHandler_GetExpenses(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/expenses?sort=", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "Correct")
	SetTimeNow()

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusOK)
	}

	var expenses []models.Expense
	err = json.Unmarshal(rr.Body.Bytes(), &expenses)
	if err != nil {
		t.Fatal(err)
	}

	if len(expenses) != 4 {
		t.Errorf("Отримано некоректну кількість витрат: отримано %d, очікувалося %d",
			len(expenses), 4)
	}
}

func TestExpensesHandler_GetExpenses_IncorrectTokenInRequest(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/expenses?sort=day", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "Incorrect")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusUnauthorized)
	}
}

func TestExpensesHandler_GetExpenses_InvalidSortParameter(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/expenses?sort=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "Correct")
	SetTimeNow()

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusMisdirectedRequest {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusMisdirectedRequest)
	}
}

func TestExpensesHandler_GetExpenses_ServerError(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/expenses?sort=day", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "TokenWithID3InDB")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusInternalServerError)
	}
}

func TestExpensesHandler_GetExpenses_IncorrectUserIdInRequest(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/expenses?sort=day", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "TokenWithoutUserInDB")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusUnauthorized)
	}
}

// -------------- END GET TESTS --------------

// -------------- PUT TESTS --------------
func TestExpensesHandler_PutExpense(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"rawdate": "2023-05-27"}`)
	req, err := http.NewRequest("PUT", "/expenses", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Correct")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusOK)
	}
}

func TestExpensesHandler_PutExpense_IncorrectTokenInRequest(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"amount": 10}`)
	req, err := http.NewRequest("PUT", "/expenses", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Incorrect")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusUnauthorized)
	}
}

func TestExpensesHandler_PutExpense_IncorrectBodyRequest(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"amount": "invalid", "date": "2023-05-27", "user_id": 1}`)
	req, err := http.NewRequest("PUT", "/expenses", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusBadRequest)
	}
}

func TestExpensesHandler_PutExpense_IncorrectUserIdInRequest(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"amount": 10}`)
	req, err := http.NewRequest("PUT", "/expenses", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "TokenWithoutUserInDB")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusUnauthorized)
	}
}

func TestExpensesHandler_PutExpense_IncorrectDateFormat(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"rawdate": "2023-05-27-2","amount": 1}`)
	req, err := http.NewRequest("PUT", "/expenses", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusBadRequest)
	}
}

func TestExpensesHandler_PutExpense_ServerError(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"rawdate": "2023-05-27","amount": -1}`)
	req, err := http.NewRequest("PUT", "/expenses", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusInternalServerError)
	}
}

// -------------- END PUT TESTS --------------

// -------------- DELETE TESTS --------------
func TestExpensesHandler_DeleteExpense_IncorrectTokenInRequest(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"amount": 10}`)
	req, err := http.NewRequest("DELETE", "/expenses/1", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Incorrect")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusUnauthorized)
	}
}

func TestExpensesHandler_DeleteExpense_IncorrectBodyRequest(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("DELETE", "/expenses/1/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusBadRequest)
	}
}

func TestExpensesHandler_DeleteExpense_IncorrectUserIdInRequest(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("DELETE", "/expenses", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "TokenWithoutUserInDB")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusUnauthorized)
	}
}

func TestExpensesHandler_DeleteExpense(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("DELETE", "/expenses/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Correct")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusOK)
	}
}

func TestExpensesHandler_DeleteExpense_NotFound(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("DELETE", "/expenses/99", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Correct")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusNotFound)
	}
}

// -------------- END DELETE TESTS --------------

// -------------- NOTALLOWEDMETHOD TESTS --------------
func TestExpensesHandler_UnknownMethodExpense_NotAllowed(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("NOTALLOWED", "/expenses/99", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Correct")

	handler := SetUpHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.Handle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusMethodNotAllowed)
	}
}

// -------------- END NOTALLOWEDMETHOD TESTS --------------
