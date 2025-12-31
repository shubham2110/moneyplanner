package transactions

import (
	"time"
)

type TransactionCreationRequest struct {
	WalletID        uint       `json:"wallet_id"`
	CategoryID      uint       `json:"category_id"`
	Amount          float64    `json:"amount"`
	PersonID        *uint      `json:"person_id,omitempty"`
	PersonName      *string    `json:"person_name,omitempty"` // To create person if not exists
	Note            *string    `json:"note,omitempty"`
	TransactionTime *time.Time `json:"transaction_time,omitempty"`
	UserID          uint       `json:"user_id"`
}

type TransactionUpdateRequest struct {
	WalletID        *uint      `json:"wallet_id,omitempty"`
	CategoryID      *uint      `json:"category_id,omitempty"`
	Amount          *float64   `json:"amount,omitempty"`
	PersonID        *uint      `json:"person_id,omitempty"`
	PersonName      *string    `json:"person_name,omitempty"`
	Note            *string    `json:"note,omitempty"`
	TransactionTime *time.Time `json:"transaction_time,omitempty"`
}
