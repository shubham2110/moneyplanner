package models

import "time"

type Wallet struct {
	WalletID         uint      `gorm:"primaryKey" json:"wallet_id"`
	Name             string    `json:"name"`
	Icon             string    `json:"icon"`
	IsEnabled        bool      `json:"is_enabled"`
	Balance          float64   `json:"balance"`
	LastModifiedTime time.Time `json:"last_modified_time"`

	// Relationships
	Categories   []Category    `gorm:"foreignKey:WalletID;references:WalletID" json:"categories,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:WalletID" json:"transactions,omitempty"`
	Users        []User        `gorm:"many2many:user_wallets;" json:"users,omitempty"`
	WalletGroups []WalletGroup `gorm:"many2many:wallet_wallet_groups;" json:"wallet_groups,omitempty"`
}

func (Wallet) TableName() string {
	return "wallets"
}
