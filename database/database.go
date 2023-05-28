package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() error {
	var err error
	dsn := "root:12345@tcp(localhost:3306)/test?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return nil
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	db.Close()
}
