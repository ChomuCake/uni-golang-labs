package services

import (
	// only for sql.ErrNoRows
	"errors"
	"testing"
	"time"

	"github.com/ChomuCake/uni-golang-labs/models"
)

var expensesBD []models.Expense

func ResetMockDB() {
	expensesBD = []models.Expense{
		{ID: 1, Amount: 10, Date: time.Now(), Category: "test", UserID: 1},
	}
}

// MockExpenseDB є замінником реалізації ExpenseDB
type MockExpenseDB struct {
}

func (db *MockExpenseDB) AddExpense(expense models.Expense) error {
	expensesBD = append(expensesBD, expense)
	return nil
}

func (db *MockExpenseDB) GetUserExpenses(userID int) ([]models.Expense, error) {
	if userID == 1 {
		return expectedExpenses, nil
	}
	return []models.Expense{}, nil
}

func (db *MockExpenseDB) UpdateUserExpenses(expense models.Expense) error {
	if expense.UserID == 2 {
		return errors.New("server error")
	}
	expensesBD[0] = expense
	return nil
}

func removeElement(slice []models.Expense, index int) []models.Expense {
	return append(slice[:index], slice[index+1:]...)
}

func (db *MockExpenseDB) DeleteExpense(expenseID string) error {
	if expenseID == "0" {
		expensesBD = removeElement(expensesBD, 0)
		return nil
	}

	return errors.New("server error")
}

// MockUserDB є замінником реалізації UserDB
type MockUserDB struct{}

func (db *MockUserDB) GetUserByID(userID int) (models.User, error) {
	if userID == 1 {
		return models.User{ID: 1, Username: "John Doe"}, nil
	}
	return models.User{}, errors.New("server error")

}

var expectedExpenses = []models.Expense{
	{ID: 1, Amount: 10, Date: time.Now(), Category: "test", UserID: 1},
	{ID: 2, Amount: 20, Date: time.Now(), Category: "test", UserID: 1},
	{ID: 3, Amount: 20, Date: time.Now().AddDate(0, 0, -1), Category: "test", UserID: 1},
	{ID: 4, Amount: 20, Date: time.Now().AddDate(0, 0, -32), Category: "test", UserID: 1},
}

var testUser = models.User{
	ID:       1,
	Username: "Test",
	Password: "12345",
}

func TestExpensesHandler_CreateExpense(t *testing.T) {
	// Arrange
	s := NewExpenseService(&MockExpenseDB{}, &MockUserDB{})
	ResetMockDB()

	// Act
	err := s.CreateExpense(testUser.ID, expectedExpenses[0])

	// Assert
	if err != nil {
		t.Errorf("Received an error: отримано %v, очікувалося %v",
			err, nil)
	}

	if len(expensesBD) != 2 {
		t.Errorf("Received incorrect number of expenses: отримано %v, очікувалося %v",
			len(expensesBD), 2)
	}

}

func TestExpenseService_CreateExpense(t *testing.T) {
	// Arrange
	mockExpenseDB := &MockExpenseDB{}
	mockUserDB := &MockUserDB{}
	s := NewExpenseService(mockExpenseDB, mockUserDB)
	ResetMockDB()

	// Act
	err := s.CreateExpense(testUser.ID, expectedExpenses[0])

	// Assert
	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

	if len(expensesBD) != 2 {
		t.Errorf("Received incorrect number of expenses: received %v, expected %v", len(expensesBD), 2)
	}
}

func TestExpenseService_GetExpenses_SortByDay(t *testing.T) {
	// Arrange
	mockExpenseDB := &MockExpenseDB{}
	mockUserDB := &MockUserDB{}
	s := NewExpenseService(mockExpenseDB, mockUserDB)
	ResetMockDB()

	// Act
	expenses, err := s.GetExpenses(testUser.ID, "day")

	// Assert
	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

	expectedCount := 2
	if len(expenses) != expectedCount {
		t.Errorf("Received incorrect number of expenses: received %v, expected %v", len(expenses), expectedCount)
	}
}

func TestExpenseService_GetExpenses_SortByMonth(t *testing.T) {
	// Arrange
	mockExpenseDB := &MockExpenseDB{}
	mockUserDB := &MockUserDB{}
	s := NewExpenseService(mockExpenseDB, mockUserDB)
	ResetMockDB()

	// Act
	expenses, err := s.GetExpenses(testUser.ID, "month")

	// Assert
	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

	expectedCount := 3
	if len(expenses) != expectedCount {
		t.Errorf("Received incorrect number of expenses: received %v, expected %v", len(expenses), expectedCount)
	}
}

func TestExpenseService_GetExpenses_SortByAll(t *testing.T) {
	// Arrange
	mockExpenseDB := &MockExpenseDB{}
	mockUserDB := &MockUserDB{}
	s := NewExpenseService(mockExpenseDB, mockUserDB)
	ResetMockDB()

	// Act
	expenses, err := s.GetExpenses(testUser.ID, "all")

	// Assert
	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

	expectedCount := 4
	if len(expenses) != expectedCount {
		t.Errorf("Received incorrect number of expenses: received %v, expected %v", len(expenses), expectedCount)
	}
}

func TestExpenseService_GetExpenses_SortByInvalid(t *testing.T) {
	// Arrange
	mockExpenseDB := &MockExpenseDB{}
	mockUserDB := &MockUserDB{}
	s := NewExpenseService(mockExpenseDB, mockUserDB)
	ResetMockDB()

	// Act
	_, err := s.GetExpenses(testUser.ID, "invalid")

	// Assert
	expectedError := "not correct sort parameter SortBy"
	if err == nil || err.Error() != expectedError {
		t.Errorf("Received incorrect error: received %v, expected %v", err, expectedError)
	}
}

func TestExpenseService_UpdateExpense(t *testing.T) {
	// Arrange
	mockExpenseDB := &MockExpenseDB{}
	mockUserDB := &MockUserDB{}
	s := NewExpenseService(mockExpenseDB, mockUserDB)
	ResetMockDB()
	ExpenseRaw := expectedExpenses[1]
	ExpenseRaw.RawDate = time.Now().Format("2006-01-02")
	ExpectedExpense := expectedExpenses[1]
	ExpectedExpense.Date, _ = time.Parse("2006-01-02", ExpenseRaw.RawDate)
	// Act
	err := s.UpdateExpense(testUser.ID, ExpenseRaw)

	// Assert
	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}
}

func TestExpenseService_DeleteExpense(t *testing.T) {
	// Arrange
	mockExpenseDB := &MockExpenseDB{}
	mockUserDB := &MockUserDB{}
	s := NewExpenseService(mockExpenseDB, mockUserDB)
	ResetMockDB()

	// Act
	err := s.DeleteExpense(testUser.ID, "0")

	// Assert
	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

	if len(expensesBD) != 0 {
		t.Errorf("Failed to delete expense")
	}
}
