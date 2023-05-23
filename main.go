package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
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
}

type ByDate []Expense

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

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
		expense.Date = time.Now()
		expenses = append(expenses, expense)
		w.WriteHeader(http.StatusCreated)
	} else if r.Method == http.MethodGet {
		if authUserID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sortExpensesBy := r.URL.Query().Get("sort")

		switch sortExpensesBy {
		case "day":
			sort.Sort(ByDate(expenses))
		case "month":
			sort.SliceStable(expenses, func(i, j int) bool {
				return expenses[i].Date.Month() < expenses[j].Date.Month()
			})
		case "all":
			sort.SliceStable(expenses, func(i, j int) bool {
				return expenses[i].Date.Before(expenses[j].Date)
			})
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(expenses)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodDelete {
		// Parse expense ID from URL path
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 3 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expenseID := pathParts[2]

		// Convert expenseID to an integer
		expenseIDInt, err := strconv.Atoi(expenseID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Find and delete the expense from the expenses slice
		for i, expense := range expenses {
			if expense.ID == expenseIDInt {
				expenses = append(expenses[:i], expenses[i+1:]...)
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		return
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
	http.HandleFunc("/expenses/", expensesHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
