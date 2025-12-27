package models

type WalletGroup struct {
	WalletGroupID   uint     `gorm:"primaryKey" json:"wallet_group_id"`
	WalletGroupName string   `json:"wallet_group_name"`

	// Relationships
	Wallets []Wallet `gorm:"many2many:wallet_wallet_groups;" json:"wallets,omitempty"`
}

func (WalletGroup) TableName() string {
	return "wallet_groups"
}
