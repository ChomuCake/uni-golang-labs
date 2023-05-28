package database

import "github.com/ChomuCake/uni-golang-labs/models"

// ExpenseDB визначає інтерфейс для роботи з даними витрат
type ExpenseDB interface {
	GetUserExpenses(userID int) ([]models.Expense, error)
	AddExpense(expense models.Expense) error
	DeleteExpense(expenseID string) error
	UpdateUserExpenses(expense models.Expense) error
}
