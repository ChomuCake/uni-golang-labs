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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.GetUserByUsername(user.Username)
	if err == nil {
		w.Header().Set("X-Error-Message", "User with this name is already registered")
		w.WriteHeader(http.StatusConflict) // Код 409 - Conflict, якщо користувач вже існує
		return
	}

	err = db.AddUser(user)
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existingUser, err := db.GetUserByUsernameAndPassword(user.Username, user.Password)
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
