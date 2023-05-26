package main

import (
	"log"
	"net/http"
	"text/template"

	db "github.com/ChomuCake/uni-golang-labs/database"
	"github.com/ChomuCake/uni-golang-labs/handlers"
	_ "github.com/go-sql-driver/mysql"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("frontend/*.html"))
	err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer db.CloseDB()

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/expenses", handlers.ExpensesHandler)
	http.HandleFunc("/expenses/", handlers.ExpensesHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
