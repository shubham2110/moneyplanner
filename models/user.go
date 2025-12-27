package models

type UserType string

const (
	UserTypeHuman UserType = "human"
	UserTypeBot   UserType = "bot"
)

type User struct {
	UserID           uint      `gorm:"primaryKey" json:"user_id"`
	Username         string    `gorm:"uniqueIndex" json:"username"`
	Name             string    `json:"name"`
	Email            string    `gorm:"uniqueIndex" json:"email"`
	Password         string    `json:"password"`
	Type             UserType  `json:"type"`
	DefaultWalletID  *uint     `json:"default_wallet_id"` // Nullable
	
	// Relationships
	Wallets      []Wallet      `gorm:"many2many:user_wallets;" json:"wallets,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
	DefaultWallet *Wallet      `gorm:"foreignKey:DefaultWalletID" json:"default_wallet,omitempty"`
}

func (User) TableName() string {
	return "users"
}
