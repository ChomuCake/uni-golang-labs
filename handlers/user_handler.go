package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	db "github.com/ChomuCake/uni-golang-labs/database"
	"github.com/ChomuCake/uni-golang-labs/models"
	"github.com/ChomuCake/uni-golang-labs/util"
	_ "github.com/go-sql-driver/mysql"
)

// DI

type userHandler struct {
	userDB db.UserDB // Використовуємо загальний інтерфейс роботи з даними UserDB(для юзерів)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	handler := &userHandler{
		userDB: &db.MySQLUserDB{
			DB: db.GetDB(),
		},
	}

	handler.regHandle(w, r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	handler := &userHandler{
		userDB: &db.MySQLUserDB{
			DB: db.GetDB(),
		},
	}

	handler.loginHandle(w, r)
}

func (h *userHandler) regHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.userDB.GetUserByUsername(user.Username)
	if err == nil {
		w.Header().Set("X-Error-Message", "User with this name is already registered")
		w.WriteHeader(http.StatusConflict) // Код 409 - Conflict, якщо користувач вже існує
		return
	}

	err = h.userDB.AddUser(user)
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

func (h *userHandler) loginHandle(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existingUser, err := h.userDB.GetUserByUsernameAndPassword(user.Username, user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenString, err := util.GenerateToken(existingUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Встановлення токена в заголовок відповіді
	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
}
