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

type TransactionFilter struct {
	StartTransactionTime  *time.Time `json:"start_transaction_time,omitempty"`
	EndTransactionTime    *time.Time `json:"end_transaction_time,omitempty"`
	StartEntryTime        *time.Time `json:"start_entry_time,omitempty"`
	EndEntryTime          *time.Time `json:"end_entry_time,omitempty"`
	StartLastModifiedTime *time.Time `json:"start_last_modified_time,omitempty"`
	EndLastModifiedTime   *time.Time `json:"end_last_modified_time,omitempty"`
	UserID                *uint      `json:"user_id,omitempty"`
	CategoryIDs           []uint     `json:"category_ids,omitempty"`
	WalletID              *uint      `json:"wallet_id,omitempty"`
	PersonID              *uint      `json:"person_id,omitempty"`
	FuzzyNote             *string    `json:"fuzzy_note,omitempty"`
	AmountOp              *string    `json:"amount_op,omitempty"` // eq, lt, le, gt, ge
	AmountValue           *float64   `json:"amount_value,omitempty"`
}
