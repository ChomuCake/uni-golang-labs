package database

import (
	"github.com/ChomuCake/uni-golang-labs/models"
	_ "github.com/go-sql-driver/mysql"
)

// --------------------------- Логіка роботи з даними для витрат ---------------------------

func GetUserExpenses(userID int) ([]models.Expense, error) {
	// Виконання запиту до бази даних для отримання витрат користувача за його ідентифікатором
	query := "SELECT id, amount, category, date FROM expenses WHERE user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ID, &expense.Amount, &expense.Category, &expense.Date)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func AddExpense(expense models.Expense) error {
	// Виконання запиту до бази даних для збереження витрати
	query := "INSERT INTO expenses (amount, category, date, user_id) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, expense.Amount, expense.Category, expense.Date, expense.UserID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteExpense(expenseID string) error {
	// Виконання запиту до бази даних для видалення витрати за її ідентифікатором
	query := "DELETE FROM expenses WHERE id = ?"
	_, err := db.Exec(query, expenseID)
	if err != nil {
		return err
	}

	return nil
}
