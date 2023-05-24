package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Expense struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	Category string    `json:"category"`
	Amount   int       `json:"amount"`
	UserID   int       `json:"user_id"`
}

type ByDate []Expense

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

var (
	tmpl      *template.Template
	db        *sql.DB
	secretKey = []byte("fd9f5dc52a0b5728c5182c593e0fae7d821e6c7a0fe64b78e67450a0a6860d63")
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = getUserByUsername(user.Username)
	if err == nil {
		w.Header().Set("X-Error-Message", "User with this name is already registered")
		w.WriteHeader(http.StatusConflict) // Код 409 - Conflict, якщо користувач вже існує
		return
	}

	err = addUser(user)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func generateToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Токен дійсний протягом 24 годин
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existingUser, err := getUserByUsernameAndPassword(user.Username, user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenString, err := generateToken(existingUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Встановлення токена в заголовок відповіді
	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
}

func expensesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var expense Expense
		err := json.NewDecoder(r.Body).Decode(&expense)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Отримання токена з заголовка авторизації
		tokenString := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))

		// Перевірка токена
		token, err := verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, ok := claims["id"].(float64)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Перевірка, чи користувач існує
		existingUser, err := getUserByID(int(userID))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expense.Date = time.Now()
		expense.UserID = existingUser.ID

		err = addExpense(expense)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else if r.Method == http.MethodGet {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		// Перевірка токена
		token, err := verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, ok := claims["id"].(float64)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Перевірка, чи користувач існує
		existingUser, err := getUserByID(int(userID))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userExpenses, err := getUserExpenses(existingUser.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sortExpensesBy := r.URL.Query().Get("sort")

		switch sortExpensesBy {
		case "day":
			sort.Sort(ByDate(userExpenses))
		case "month":
			sort.SliceStable(userExpenses, func(i, j int) bool {
				return userExpenses[i].Date.Month() < userExpenses[j].Date.Month()
			})
		case "all":
			sort.SliceStable(userExpenses, func(i, j int) bool {
				return userExpenses[i].Date.Before(userExpenses[j].Date)
			})
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(userExpenses)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodDelete {
		// Отримання токена з заголовка авторизації
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		// Перевірка токена
		_, err := verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Розбиття URL шляху для отримання ID витрати
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 3 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expenseID := pathParts[2]

		err = deleteExpense(expenseID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func initDB() error {
	var err error
	dsn := "root:12345@tcp(localhost:3306)/test?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return nil
}

func closeDB() {
	db.Close()
}

func addUser(user User) error {
	stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func getUserByUsernameAndPassword(username, password string) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, username FROM users WHERE username = ? AND password = ?", username, password).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func getUserByUsername(username string) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, username FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func getUserByID(userID int) (User, error) {
	// Виконання запиту до бази даних для отримання користувача за його ідентифікатором
	query := "SELECT id, username FROM users WHERE id = ?"
	row := db.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("User not found")
		}
		return User{}, err
	}

	return user, nil
}

func getUserExpenses(userID int) ([]Expense, error) {
	// Виконання запиту до бази даних для отримання витрат користувача за його ідентифікатором
	query := "SELECT id, amount, category, date FROM expenses WHERE user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var expense Expense
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

func addExpense(expense Expense) error {
	// Виконання запиту до бази даних для збереження витрати
	query := "INSERT INTO expenses (amount, category, date, user_id) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, expense.Amount, expense.Category, expense.Date, expense.UserID)
	if err != nil {
		return err
	}

	return nil
}

func deleteExpense(expenseID string) error {
	// Виконання запиту до бази даних для видалення витрати за її ідентифікатором
	query := "DELETE FROM expenses WHERE id = ?"
	_, err := db.Exec(query, expenseID)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	tmpl = template.Must(template.ParseFiles("pagesHTML/index.html"))
	err := initDB()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer closeDB()

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/expenses", expensesHandler)
	http.HandleFunc("/expenses/", expensesHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
