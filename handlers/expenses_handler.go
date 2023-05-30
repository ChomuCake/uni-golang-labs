package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/ChomuCake/uni-golang-labs/models"
	_ "github.com/go-sql-driver/mysql"
)

type ByDate []models.Expense

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

// інтерфейс expenseService, tokenManager описується в тому ж файлі що і використовується
type expenseService interface {
	CreateExpense(userID int, expense models.Expense) error
	GetExpenses(userID int, sortExpensesBy string) ([]models.Expense, error)
	UpdateExpense(userID int, updatedExpense models.Expense) error
	DeleteExpense(userID int, expenseID string) error
}

type tokenManager interface {
	ExtractUserIDFromRequest(r *http.Request) (int, error)
}

type ExpenseHandler struct {
	expService expenseService
	tokenMng   tokenManager
}

func NewExpenseHandler(expService expenseService, tokenMng tokenManager) *ExpenseHandler {
	return &ExpenseHandler{
		expService: expService,
		tokenMng:   tokenMng,
	}
}

func (h *ExpenseHandler) RegisterRoutes(router *httprouter.Router) {
	router.POST("/expenses", h.CreateExpense)
	router.GET("/expenses", h.GetExpenses)
	router.PUT("/expenses", h.UpdateExpense)
	router.DELETE("/expenses/:id", h.DeleteExpense)
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var expense models.Expense
	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отримання айді користувача з заголовка авторизації
	userID, err := h.tokenMng.ExtractUserIDFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Створення витрат
	err = h.expService.CreateExpense(userID, expense)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Отримання айді користувача з заголовка авторизації
	userID, err := h.tokenMng.ExtractUserIDFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sortExpensesBy := r.URL.Query().Get("sort")

	// Отримання витрат
	userExpenses, err := h.expService.GetExpenses(userID, sortExpensesBy)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userExpenses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var updatedExpense models.Expense
	err := json.NewDecoder(r.Body).Decode(&updatedExpense)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отримання айді користувача з заголовка авторизації
	userID, err := h.tokenMng.ExtractUserIDFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Оновлення витрати
	err = h.expService.UpdateExpense(userID, updatedExpense)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Отримання айді користувача з заголовка авторизації
	userID, err := h.tokenMng.ExtractUserIDFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Видалення витрати
	err = h.expService.DeleteExpense(userID, params.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
