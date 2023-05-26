package database

import (
	"database/sql"

	"github.com/ChomuCake/uni-golang-labs/models"
	_ "github.com/go-sql-driver/mysql"
)

// --------------------------- Логіка роботи з даними для витрат (MySQL) ---------------------------
type MySQLExpenseDB struct {
	DB *sql.DB
}

func (db *MySQLExpenseDB) GetUserExpenses(userID int) ([]models.Expense, error) {
	// Виконання запиту до бази даних для отримання витрат користувача за його ідентифікатором
	query := "SELECT id, amount, category, date FROM expenses WHERE user_id = ?"
	rows, err := db.DB.Query(query, userID)
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

func (db *MySQLExpenseDB) AddExpense(expense models.Expense) error {
	// Виконання запиту до бази даних для збереження витрати
	query := "INSERT INTO expenses (amount, category, date, user_id) VALUES (?, ?, ?, ?)"
	_, err := db.DB.Exec(query, expense.Amount, expense.Category, expense.Date, expense.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (db *MySQLExpenseDB) DeleteExpense(expenseID string) error {
	// Виконання запиту до бази даних для видалення витрати за її ідентифікатором
	query := "DELETE FROM expenses WHERE id = ?"
	_, err := db.DB.Exec(query, expenseID)
	if err != nil {
		return err
	}

	return nil
}

func (db *MySQLExpenseDB) UpdateUserExpenses(expense models.Expense) error {
	// Виконання запиту до бази даних для оновлення витрати
	query := "UPDATE expenses SET amount = ?, category = ?, date = ? WHERE id = ?"
	_, err := db.DB.Exec(query, expense.Amount, expense.Category, expense.Date, expense.ID)
	if err != nil {
		return err
	}

	return nil
}
