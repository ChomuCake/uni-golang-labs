package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	db "github.com/ChomuCake/uni-golang-labs/database"
	"github.com/ChomuCake/uni-golang-labs/models"
	"github.com/ChomuCake/uni-golang-labs/util"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

type ByDate []models.Expense

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var expense models.Expense
		err := json.NewDecoder(r.Body).Decode(&expense)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Отримання токена з заголовка авторизації
		tokenString := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))

		// Перевірка токена
		token, err := util.VerifyToken(tokenString)
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
		existingUser, err := db.GetUserByID(int(userID))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expense.Date = time.Now()
		expense.UserID = existingUser.ID

		err = db.AddExpense(expense)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else if r.Method == http.MethodGet {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		// Перевірка токена
		token, err := util.VerifyToken(tokenString)
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
		existingUser, err := db.GetUserByID(int(userID))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userExpenses, err := db.GetUserExpenses(existingUser.ID)
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
		_, err := util.VerifyToken(tokenString)
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

		err = db.DeleteExpense(expenseID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
