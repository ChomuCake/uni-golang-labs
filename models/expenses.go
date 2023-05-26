package models

import (
	"time"
)

type Expense struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	RawDate  string    `json:"rawdate"`
	Category string    `json:"category"`
	Amount   int       `json:"amount"`
	UserID   int       `json:"user_id"`
}
