package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Expense struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Amount   int    `json:"amount"`
}

var (
	tmpl              *template.Template
	currentUserID     int
	authUserID        int
	currentExpensesID int
	users             []User
	expenses          []Expense
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

	// Перевірка на існування користувача з таким же ім'ям
	for _, u := range users {
		if u.Username == user.Username {
			// Користувач з таким ім'ям вже існує
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	currentUserID++
	user.ID = currentUserID
	users = append(users, user)

	w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Perform user authentication
	for _, u := range users {
		if u.Username == user.Username && u.Password == user.Password {
			authUserID = u.ID
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func expensesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var expense Expense
		err := json.NewDecoder(r.Body).Decode(&expense)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Якщо користувач не авторизований
		if authUserID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		currentExpensesID++
		expense.ID = currentExpensesID
		expenses = append(expenses, expense)
		w.WriteHeader(http.StatusCreated)
	} else if r.Method == http.MethodGet {
		if authUserID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(expenses)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func init() {
	tmpl = template.Must(template.ParseFiles("pagesHTML/index.html"))
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/expenses", expensesHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
