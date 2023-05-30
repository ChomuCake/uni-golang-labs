package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type RealDatabase struct {
	// реалізація основної бази даних
	RealDBName string
	db_real    *sql.DB
}

func (db *RealDatabase) InitDB() error {
	db.RealDBName = "test"
	var err error
	dsn := "root:12345@tcp(localhost:3306)/" + db.RealDBName + "?parseTime=true"
	db.db_real, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return nil
}

func (db *RealDatabase) GetDB() *sql.DB {
	return db.db_real
}

func (db *RealDatabase) CloseDB() {
	db.db_real.Close()
}
