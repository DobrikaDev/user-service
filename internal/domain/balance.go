package domain

import "time"

type Balance struct {
	ID      string `json:"id" db:"id"`
	UserID  string `json:"user_id" db:"user_id"`
	Balance int    `json:"balance" db:"balance"`
}

type BalanceOperation struct {
	ID          string               `json:"id" db:"id"`
	BalanceID   string               `json:"balance_id" db:"balance_id"`
	Amount      int                  `json:"amount" db:"amount"`
	Type        BalanceOperationType `json:"type" db:"type"`
	Description string               `json:"description" db:"description"`
	CreatedAt   time.Time            `json:"created_at" db:"created_at"`
}
