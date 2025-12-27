package models

type Category struct {
	CategoryID uint   `gorm:"primaryKey" json:"category_id"`
	Icon       string `json:"icon"`
	Name       string `json:"name"`
	ParentID   *uint  `json:"parent_id"` // Nullable foreign key to Category
	RootID     uint   `json:"root_id"`
	WalletID   uint   `json:"wallet_id"`
	IsGlobal   bool   `json:"is_global"`

	// Relationships
	Parent       *Category      `gorm:"foreignKey:ParentID;references:CategoryID" json:"parent,omitempty"`
	Transactions []Transaction  `gorm:"foreignKey:CategoryID" json:"transactions,omitempty"`
	Wallet       Wallet         `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
