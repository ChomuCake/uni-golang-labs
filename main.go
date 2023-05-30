package main

import (
	"log"
	"net/http"

	db "github.com/ChomuCake/uni-golang-labs/database"
	"github.com/ChomuCake/uni-golang-labs/handlers"
	"github.com/ChomuCake/uni-golang-labs/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func init() {
	err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer db.CloseDB()

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	router := httprouter.New()
	router.POST("/expenses", handlers.NewExpenseHandler(&db.MySQLExpenseDB{DB: db.GetDB()}, &db.MySQLUserDB{DB: db.GetDB()}, &util.JWTTokenManager{}).CreateExpense)
	router.GET("/expenses", handlers.NewExpenseHandler(&db.MySQLExpenseDB{DB: db.GetDB()}, &db.MySQLUserDB{DB: db.GetDB()}, &util.JWTTokenManager{}).GetExpenses)
	router.PUT("/expenses", handlers.NewExpenseHandler(&db.MySQLExpenseDB{DB: db.GetDB()}, &db.MySQLUserDB{DB: db.GetDB()}, &util.JWTTokenManager{}).UpdateExpense)
	router.DELETE("/expenses/:id", handlers.NewExpenseHandler(&db.MySQLExpenseDB{DB: db.GetDB()}, &db.MySQLUserDB{DB: db.GetDB()}, &util.JWTTokenManager{}).DeleteExpense)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
