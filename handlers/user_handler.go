package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ChomuCake/uni-golang-labs/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

type userService interface {
	RegisterUser(user models.User) error
	LoginUser(user models.User) (models.User, error)
}

type tokenManagerUser interface {
	GenerateToken(user models.User) (string, error)
}

type UserHandler struct {
	uService userService
	tokenMng tokenManagerUser
}

func NewUserHandler(uService userService, tokenMng tokenManagerUser) *UserHandler {
	return &UserHandler{
		uService: uService,
		tokenMng: tokenMng,
	}
}

func (h *UserHandler) RegisterRoutesUser(router *httprouter.Router) {
	router.POST("/register", h.RegisterUser)
	router.POST("/login", h.LoginUser)
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.uService.RegisterUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existingUser, err := h.uService.LoginUser(user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, err := h.tokenMng.GenerateToken(existingUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Встановлення токена в заголовок відповіді
	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
}
