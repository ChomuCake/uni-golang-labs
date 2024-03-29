//------------------------Інтеграційний тест------------------------------

package drepo

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/ChomuCake/uni-golang-labs/models"
)

type TestDatabase struct {
	// реалізація тестової бaзи даних
	testDBName string
	db_test    *sql.DB
}

func (db *TestDatabase) GetDB() *sql.DB {
	return db.db_test
}

func (db *TestDatabase) InitDB() error {
	// Формування рядка підключення до тестової бази даних
	db.testDBName = "test_db"
	dsn := "root:12345@tcp(localhost:3306)/" + db.testDBName + "?parseTime=true"

	// Встановлення з'єднання з тестовою базою даних
	var err error
	db.db_test, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %v", err)
	}

	// Очищення тестової бази даних перед початком тестів
	if err := db.СlearTestDB(); err != nil {
		return fmt.Errorf("failed to clear test database: %v", err)
	}

	return nil
}

func (db *TestDatabase) CloseDB() {
	db.db_test.Close()
}

func (db *TestDatabase) СlearTestDB() error {
	// Видалення і створення бази даних
	_, err := db.db_test.Exec("DROP DATABASE IF EXISTS " + db.testDBName)
	if err != nil {
		return fmt.Errorf("failed to drop test database: %v", err)
	}

	_, err = db.db_test.Exec("CREATE DATABASE " + db.testDBName)
	if err != nil {
		return fmt.Errorf("failed to create test database: %v", err)
	}

	_, err = db.db_test.Exec("USE " + db.testDBName)
	if err != nil {
		return fmt.Errorf("failed to switch to test database: %v", err)
	}

	// Створення таблиці `users`
	_, err = db.db_test.Exec(`
		CREATE TABLE users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	// Створення таблиці `expenses`
	_, err = db.db_test.Exec(`
		CREATE TABLE expenses (
			id INT AUTO_INCREMENT PRIMARY KEY,
			date DATE NOT NULL,
			category VARCHAR(255) NOT NULL,
			amount INT NOT NULL,
			user_id INT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create expenses table: %v", err)
	}

	return nil
}

func TestGetUserExpensesIntegration(t *testing.T) {
	// Сворення тестової бд
	db := &TestDatabase{}

	// Підготовка тестової бази даних
	err := db.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	defer db.CloseDB()

	// Створення репо витрат
	ExpenseDB := NewExpenseDBMySQL(db)

	// Створення репо юзерів
	userDB := NewUserDBMySQL(db)

	newUser := models.User{
		Username: "TestName",
		Password: "12345",
	}

	// GetUserByID повинен повертати тільки ім'я та айді користувача
	expectedUser := models.User{
		Username: newUser.Username,
		ID:       1,
	}

	// Створення об'єкту моделі витрат
	newExpense := models.Expense{
		ID:       1,
		Date:     time.Now().Truncate(24 * time.Hour).UTC(),
		Category: "TestExpenses",
		Amount:   100,
		UserID:   expectedUser.ID,
	}

	// GetUserExpenses повинен усе крім юзерАЙді(бо нема сенсу)
	expectedExpenses := models.Expense{
		ID:       1,
		Date:     time.Now().Truncate(24 * time.Hour).UTC(),
		Category: "TestExpenses",
		Amount:   100,
	}

	ExpensesUpdate := models.Expense{
		ID:       newExpense.ID,
		Date:     newExpense.Date,
		Category: "Updated " + newExpense.Category,
		Amount:   999 + newExpense.Amount,
	}

	// Тестування створення і отримання користувача
	// Результат після створення користувача він має отримуватись з бд
	t.Run("create and get User", func(t *testing.T) {
		err = userDB.AddUser(newUser)

		if err != nil {
			t.Errorf("failed to add user with error: %v", err)
		}

		fmt.Println(newUser)
		user, err := userDB.GetUserByID(expectedUser.ID)
		if err != nil {
			t.Errorf("failed to get user with error: %v", err)
		}

		if !reflect.DeepEqual(expectedUser, user) {
			t.Errorf("expenses data is corrupted; actual: %v, expected: %v", user, expectedUser)
		}
	})

	// Тестування створення і отримання витрат користувача
	// Результат користувач повинен отримувати нову витрату після створення її у бд
	t.Run("create and get UserExpneses", func(t *testing.T) {
		err = ExpenseDB.AddExpense(newExpense)

		if err != nil {
			t.Errorf("failed to add expense with error: %v", err)
		}

		fmt.Println(newExpense)
		expense, err := ExpenseDB.GetUserExpenses(expectedUser.ID)
		if err != nil {
			t.Errorf("failed to get user expneses with error: %v", err)
		}

		if !reflect.DeepEqual(expectedExpenses, expense[0]) {
			t.Errorf("expenses data is corrupted; actual: %v, expected: %v", expense[0], expectedExpenses)
		}
	})

	// Тестування оновлення і отримання витрат користувача
	// Результат користувач повинен отримувати оновлені витрати після оновлення їх у бд
	t.Run("update and get UserExpnese", func(t *testing.T) {
		err = ExpenseDB.UpdateUserExpenses(ExpensesUpdate)

		if err != nil {
			t.Errorf("failed update expense with error: %v", err)
		}

		fmt.Println(ExpensesUpdate)
		expense, err := ExpenseDB.GetUserExpenses(expectedUser.ID)
		if err != nil {
			t.Errorf("failed to get user expneses with error: %v", err)
		}

		if !reflect.DeepEqual(ExpensesUpdate, expense[0]) {
			t.Errorf("expenses data is corrupted; actual: %v, expected: %v", expense[0], ExpensesUpdate)
		}
	})

	// Тестування видалення і отримання витрат користувача
	// Результат користувач повинен отримувати 0 витрат після видалення їх з бд
	t.Run("delete and get UserExpnese", func(t *testing.T) {
		err = ExpenseDB.DeleteExpense(strconv.Itoa(ExpensesUpdate.ID))

		if err != nil {
			t.Errorf("failed to delete expense with error: %v", err)
		}

		expense, err := ExpenseDB.GetUserExpenses(expectedUser.ID)
		if err != nil {
			t.Errorf("failed to get user expneses with error: %v", err)
		}

		if !reflect.DeepEqual(0, len(expense)) {
			t.Errorf("expenses data is corrupted; actual: %v, expected: %v", expense, 0)
		}
	})

	// Тестування отримання користувача за ім'ям, та за ім'ям і паролем
	// Результат користувач повинен бути однаковим при кожному отримані з бд
	t.Run("get user by username and get user by username and password", func(t *testing.T) {
		userGet1, err := userDB.GetUserByUsername(newUser.Username)
		if err != nil {
			t.Errorf("failed to get user with error: %v", err)
		}

		userGet2, err := userDB.GetUserByUsernameAndPassword(newUser.Username, newUser.Password)
		if err != nil {
			t.Errorf("failed to get user with error: %v", err)
		}

		if !reflect.DeepEqual(userGet1, userGet2) {
			t.Errorf("expenses data is corrupted; actual: %v, expected: %v", userGet1, userGet2)
		}
	})

	// Закінчення тестування
	log.Println("Integration test completed.")
}
