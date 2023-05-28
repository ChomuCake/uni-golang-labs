// nolint
package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ChomuCake/uni-golang-labs/database"
	"github.com/ChomuCake/uni-golang-labs/models"
	"github.com/ChomuCake/uni-golang-labs/util"
)

const testDBName = "benchmark_test_db"

func InitTestDB() (*sql.DB, error) {
	// Формування рядка підключення до тестової бази даних
	dsn := "root:12345@tcp(localhost:3306)/" + testDBName + "?parseTime=true"

	// Встановлення з'єднання з тестовою базою даних
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %v", err)
	}

	// Очищення тестової бази даних перед початком тестів
	if err := СlearTestDB(db); err != nil {
		return nil, fmt.Errorf("failed to clear test database: %v", err)
	}

	return db, nil
}

func СlearTestDB(db *sql.DB) error {
	// Видалення і створення бази даних
	_, err := db.Exec("DROP DATABASE IF EXISTS " + testDBName)
	if err != nil {
		return fmt.Errorf("failed to drop test database: %v", err)
	}

	_, err = db.Exec("CREATE DATABASE " + testDBName)
	if err != nil {
		return fmt.Errorf("failed to create test database: %v", err)
	}

	_, err = db.Exec("USE " + testDBName)
	if err != nil {
		return fmt.Errorf("failed to switch to test database: %v", err)
	}

	// Створення таблиці `users`
	_, err = db.Exec(`
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
	_, err = db.Exec(`
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

func BenchmarkGetUserExpenses(b *testing.B) {
	testDB, err := InitTestDB()
	if err != nil {
		b.Fatalf("Failed to initialize test database: %v", err)
	}
	defer testDB.Close()

	expenseDB := &database.MySQLExpenseDB{
		DB: testDB,
	}

	userDB := &database.MySQLUserDB{
		DB: testDB,
	}

	jwtToken := &util.JWTTokenManager{}

	benchmarkUser := models.User{
		ID:       1,
		Username: "benchmarkUser",
		Password: "12345",
	}

	benmarkExpense := models.Expense{
		Amount:   100,
		Category: "Test",
		Date:     time.Now().UTC(),
		UserID:   1,
	}

	err = userDB.AddUser(benchmarkUser)
	if err != nil {
		b.Errorf("failed to add user with error: %v", err)
	}

	err = expenseDB.AddExpense(benmarkExpense)
	if err != nil {
		b.Errorf("failed to add expense with error: %v", err)
	}

	tokenStr, err := jwtToken.GenerateToken(benchmarkUser)
	if err != nil {
		b.Errorf("failed to generate token with error: %v", err)
	}

	handler := &ExpenseHandler{
		ExpenseDB: expenseDB,
		UserDB:    userDB,
		TokenMng:  jwtToken,
	}

	req, err := http.NewRequest("GET", "/expenses?sort=all", nil)
	if err != nil {
		b.Fatalf("Failed to create request: %v", err)
	}

	// Моделюємо авторизованого користувача, додавши заголовок авторизації
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	b.ResetTimer()

	startTime := time.Now()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			b.Errorf("Expected status 200 OK, but got %d", rr.Code)
		}
	}

	duration := time.Since(startTime)
	b.Logf("Benchmark duration: %s\n", duration)
}
