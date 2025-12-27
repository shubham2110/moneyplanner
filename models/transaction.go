package models

import "time"

type Transaction struct {
	TransactionID    uint      `gorm:"primaryKey" json:"transaction_id"`
	CategoryID       uint      `json:"category_id"`
	Amount           float64   `json:"amount"`
	Note             *string   `json:"note"` // Nullable
	PersonID         *uint     `json:"person_id"` // Nullable foreign key
	WalletID         uint      `json:"wallet_id"`
	TransactionTime  time.Time `json:"transaction_time"`
	EntryTime        time.Time `json:"entry_time"`
	LastModifiedTime time.Time `json:"last_modified_time"`
	UserID           uint      `json:"user_id"`

	// Relationships
	Category Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Person   *Person   `gorm:"foreignKey:PersonID" json:"person,omitempty"`
	Wallet   Wallet    `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}
