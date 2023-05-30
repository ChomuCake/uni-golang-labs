package services

import (
	"errors"
	"sort"
	"time"

	"github.com/ChomuCake/uni-golang-labs/models"
)

type ByDate []models.Expense

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

type ExpenseDB interface {
	GetUserExpenses(userID int) ([]models.Expense, error)
	AddExpense(expense models.Expense) error
	DeleteExpense(expenseID string) error
	UpdateUserExpenses(expense models.Expense) error
}

type UserDB interface {
	GetUserByID(userID int) (models.User, error)
}

type ExpenseService struct {
	expenseDB ExpenseDB
	userDB    UserDB
}

func NewExpenseService(expenseDB ExpenseDB, userDB UserDB) *ExpenseService {
	return &ExpenseService{expenseDB, userDB}
}

func (s *ExpenseService) CreateExpense(userID int, expense models.Expense) error {
	// Перевірка, чи користувач існує
	_, err := s.userDB.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Парсинг рядкового значення дати
	parsedDate, err := time.Parse("2006-01-02", expense.RawDate)
	if err != nil {
		return errors.New("invalid date format")
	}

	// Оновлення поля Date
	expense.Date = parsedDate

	// Створення витрати
	err = s.expenseDB.AddExpense(expense)
	if err != nil {
		return errors.New("failed to create expense")
	}

	return nil
}

func (s *ExpenseService) GetExpenses(userID int, sortExpensesBy string) ([]models.Expense, error) {
	// Перевірка, чи користувач існує
	_, err := s.userDB.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	userExpenses, err := s.expenseDB.GetUserExpenses(userID)
	if err != nil {
		return nil, errors.New("failed to get user expenses")
	}

	switch sortExpensesBy {
	case "day":
		today := time.Now().Truncate(24 * time.Hour) // Отримуємо поточну дату без часу
		var todayExpenses []models.Expense

		// Фільтруємо витрати за сьогоднішній день
		for _, expense := range userExpenses {
			if expense.Date.Year() == today.Year() &&
				expense.Date.Month() == today.Month() &&
				expense.Date.Day() == today.Day() {
				todayExpenses = append(todayExpenses, expense)
			}
		}
		userExpenses = todayExpenses

	case "month":
		month := time.Now().Month() // Поточний місяць
		var monthExpenses []models.Expense

		// Фільтруємо витрати за поточний місяць
		for _, expense := range userExpenses {
			if expense.Date.Month() == month {
				monthExpenses = append(monthExpenses, expense)
			}
		}
		userExpenses = monthExpenses

	case "all":
		sort.SliceStable(userExpenses, func(i, j int) bool {
			return userExpenses[i].Date.Before(userExpenses[j].Date)
		})
	default:
		if sortExpensesBy != "" {
			return nil, errors.New("not correct sort parameter SortBy")
		}

		sort.SliceStable(userExpenses, func(i, j int) bool {
			return userExpenses[i].Date.Before(userExpenses[j].Date)
		})
	}

	return userExpenses, nil
}

func (s *ExpenseService) UpdateExpense(userID int, updatedExpense models.Expense) error {
	// Перевірка, чи користувач існує
	_, err := s.userDB.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Парсинг рядкового значення дати
	updatedExpense.Date, err = time.Parse("2006-01-02", updatedExpense.RawDate)
	if err != nil {
		return errors.New("failed to parse date expense")
	}

	// Оновлення витрати
	err = s.expenseDB.UpdateUserExpenses(updatedExpense)
	if err != nil {
		return errors.New("failed to update expense")
	}

	return nil
}

func (s *ExpenseService) DeleteExpense(userID int, expenseID string) error {
	// Перевірка, чи користувач існує
	_, err := s.userDB.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = s.expenseDB.DeleteExpense(expenseID)
	if err != nil {
		return errors.New("failed to delete expense")
	}

	return nil
}
