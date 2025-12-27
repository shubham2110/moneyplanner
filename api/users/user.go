package users

import (
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
)

// CreateUser creates a new user in the database
func CreateUser(username, name, email, password string, userType models.UserType, defaultWalletID *uint) (*models.User, error) {
	user := &models.User{
		Username:        username,
		Name:            name,
		Email:           email,
		Password:        password,
		Type:            userType,
		DefaultWalletID: defaultWalletID,
	}
	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}
	log.Printf("✓ User '%s' created", username)
	return user, nil
}
