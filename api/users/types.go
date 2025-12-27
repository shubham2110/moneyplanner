package users

import "moneyplanner/models"

// UserCreationRequest represents the request to create a new user with default wallet and categories
type UserCreationRequest struct {
	Username          string		       	`json:"username"`
	Name              string				`json:"name"`
	Email             string				`json:"email"`
	Password          string				`json:"password"`
	UserType          models.UserType		`json:"user_type"`
	WalletName        string				`json:"wallet_name"`
	WalletGroupName   string				`json:"wallet_group_name"`
	CreateCategories  bool 					`json:"create_categories"`
}

// UserCreationResult contains the created user and related entities
type UserCreationResult struct {
	User               *models.User        `json:"user"`
	Wallet             *models.Wallet      `json:"wallet"`
	WalletGroup        *models.WalletGroup `json:"wallet_group"`
	RootCategories     []models.Category   `json:"root_categories,omitempty"`
	DummyTransactions  []models.Transaction `json:"dummy_transactions,omitempty"`
}

// UserUpdateRequest represents a request to update user details
type UserUpdateRequest struct {
	Name        *string            `json:"name,omitempty"`
	Email       *string            `json:"email,omitempty"`
	Password    *string            `json:"password,omitempty"`
	Type        *models.UserType   `json:"type,omitempty"`
	DefaultWalletID *uint          `json:"default_wallet_id,omitempty"`
}

// UserResponse represents a user response
type UserResponse struct {
	UserID           uint             `json:"user_id"`
	Username         string           `json:"username"`
	Name             string           `json:"name"`
	Email            string           `json:"email"`
	Type             models.UserType  `json:"type"`
	DefaultWalletID  *uint            `json:"default_wallet_id,omitempty"`
}

