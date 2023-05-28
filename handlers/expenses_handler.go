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
	_ "github.com/go-sql-driver/mysql"
)

type ByDate []models.Expense

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

// DI

type ExpenseHandler struct {
	ExpenseDB db.ExpenseDB      // Використовуємо загальний інтерфейс роботи з даними ExpenseDB(для витрат)
	UserDB    db.UserDB         // Використовуємо загальний інтерфейс роботи з даними UserDB(для юзерів)
	TokenMng  util.TokenManager // Використовуємо загальний інтерфейс роботи з токенами
}

// Функція ExpensesHandler, яка обробляє запити. У цій функції ми створюємо екземпляр expenseHandler
// та передаємо йому залежність - екземпляр db.MySQLExpenseDB(конкретна реалізація)
// та екземпляр db.MySQLUserDB(конкретна реалізація)
func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	handler := &ExpenseHandler{
		ExpenseDB: &db.MySQLExpenseDB{
			DB: db.GetDB(),
		},
		UserDB: &db.MySQLUserDB{
			DB: db.GetDB(),
		},
		TokenMng: &util.JWTTokenManager{},
	}

	handler.Handle(w, r)
}

func (h *ExpenseHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var expense models.Expense
		err := json.NewDecoder(r.Body).Decode(&expense)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Отримання айді користувача з заголовка авторизації
		userID, err := h.TokenMng.ExtractUserIDFromRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Перевірка, чи користувач існує
		existingUser, err := h.UserDB.GetUserByID(int(userID))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expense.Date = time.Now()
		expense.UserID = existingUser.ID

		err = h.ExpenseDB.AddExpense(expense)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else if r.Method == http.MethodGet {
		// Отримання айді користувача з заголовка авторизації
		userID, err := h.TokenMng.ExtractUserIDFromRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Перевірка, чи користувач існує
		existingUser, err := h.UserDB.GetUserByID(int(userID))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userExpenses, err := h.ExpenseDB.GetUserExpenses(existingUser.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sortExpensesBy := r.URL.Query().Get("sort")

		switch sortExpensesBy {
		case "day":
			today := time.Now().Truncate(24 * time.Hour) // Отримуємо поточну дату без часу
			var todayExpenses []models.Expense

			// Фільтруємо витрати за сьогоднішній день
			for _, expense := range userExpenses {
				if expense.Date.Year() == today.Year() &&
					expense.Date.Month() == today.Month() &&
					expense.Date.Day() == today.Day() {
					todayExpenses = append(todayExpenses, expense)
				}
			}
			userExpenses = todayExpenses

		case "month":
			month := time.Now().Month() // Поточний місяць
			var monthExpenses []models.Expense

			// Фільтруємо витрати за поточний місяць
			for _, expense := range userExpenses {
				if expense.Date.Month() == month {
					monthExpenses = append(monthExpenses, expense)
				}
			}
			userExpenses = monthExpenses

		case "all":
			sort.SliceStable(userExpenses, func(i, j int) bool {
				return userExpenses[i].Date.Before(userExpenses[j].Date)
			})
		default:
			if sortExpensesBy != "" {
				w.WriteHeader(http.StatusMisdirectedRequest)
				return
			}

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
	} else if r.Method == http.MethodPut {
		// Отримання айді користувача з заголовка авторизації
		userID, err := h.TokenMng.ExtractUserIDFromRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Перевірка, чи користувач існує
		_, err = h.UserDB.GetUserByID(int(userID))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var updatedExpense models.Expense
		err = json.NewDecoder(r.Body).Decode(&updatedExpense)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Парсинг рядкового значення дати
		parsedDate, err := time.Parse("2006-01-02", updatedExpense.RawDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Оновлення поля Date
		updatedExpense.Date = parsedDate

		// Оновлення витрати
		err = h.ExpenseDB.UpdateUserExpenses(updatedExpense)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodDelete {
		// Отримання айді користувача з заголовка авторизації
		userID, err := h.TokenMng.ExtractUserIDFromRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Перевірка, чи користувач існує
		_, err = h.UserDB.GetUserByID(int(userID))
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

		err = h.ExpenseDB.DeleteExpense(expenseID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
