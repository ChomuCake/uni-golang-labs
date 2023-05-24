package database

import (
	"database/sql"
	"fmt"

	"github.com/ChomuCake/uni-golang-labs/models"
	_ "github.com/go-sql-driver/mysql"
)

// --------------------------- Логіка роботи з даними для юзера ---------------------------

func AddUser(user models.User) error {
	stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByUsernameAndPassword(username, password string) (models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username FROM users WHERE username = ? AND password = ?", username, password).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserByID(userID int) (models.User, error) {
	// Виконання запиту до бази даних для отримання користувача за його ідентифікатором
	query := "SELECT id, username FROM users WHERE id = ?"
	row := db.QueryRow(query, userID)

	var user models.User
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("User not found")
		}
		return models.User{}, err
	}

	return user, nil
}
