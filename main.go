package main

import (
	"log"
	"net/http"

	db "github.com/ChomuCake/uni-golang-labs/database"
	"github.com/ChomuCake/uni-golang-labs/drepo"
	"github.com/ChomuCake/uni-golang-labs/handlers"
	"github.com/ChomuCake/uni-golang-labs/services"
	"github.com/ChomuCake/uni-golang-labs/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func main() {
	DB := &db.RealDatabase{}
	defer DB.CloseDB()

	err := DB.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()

	expenseDB := drepo.NewExpenseDBMySQL(DB)
	userDB := drepo.NewUserDBMySQL(DB)

	tokenManager := util.JWTTokenManager{}
	expenseService := services.NewExpenseService(expenseDB, userDB)
	expenseHandler := handlers.NewExpenseHandler(expenseService, tokenManager)
	expenseHandler.RegisterRoutes(router)

	userService := services.NewUserService(userDB)
	userHandler := handlers.NewUserHandler(userService, tokenManager)
	userHandler.RegisterRoutesUser(router)

	fs := http.FileServer(http.Dir("./frontend"))
	router.NotFound = fs

	log.Fatal(http.ListenAndServe(":8080", router))
}
