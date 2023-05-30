package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ChomuCake/uni-golang-labs/drepo"
	"github.com/ChomuCake/uni-golang-labs/models"
	"github.com/ChomuCake/uni-golang-labs/services"
	"github.com/ChomuCake/uni-golang-labs/util"
	"github.com/julienschmidt/httprouter"
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
	db.testDBName = "benchmark_test_db"
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

func BenchmarkGetUserExpenses(b *testing.B) {
	// Сворення тестової бд
	db := &TestDatabase{}

	// Підготовка тестової бази даних
	err := db.InitDB()
	if err != nil {
		b.Fatalf("Failed to initialize test database: %v", err)
	}
	defer db.CloseDB()

	// Створення репо витрат
	expenseDB := drepo.NewExpenseDBMySQL(db)

	// Створення репо юзерів
	userDB := drepo.NewUserDBMySQL(db)

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

	s := services.NewExpenseService(expenseDB, userDB)

	h := NewExpenseHandler(s, jwtToken)

	router := httprouter.New()
	h.RegisterRoutes(router)

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
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			b.Errorf("Expected status 200 OK, but got %d", rr.Code)
		}
	}

	duration := time.Since(startTime)
	b.Logf("Benchmark duration: %s\n", duration)
}
