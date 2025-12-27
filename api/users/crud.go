package users

import (
	"fmt"
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
)

// GetUserByID retrieves a user by their ID
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

// ListAllUsers retrieves all users from the database
func ListAllUsers() ([]models.User, error) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// UpdateUser updates user details
func UpdateUser(userID uint, req *UserUpdateRequest) (*models.User, error) {
	user, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = *req.Name
		user.Name = *req.Name
	}

	if req.Email != nil {
		updates["email"] = *req.Email
		user.Email = *req.Email
	}

	if req.Password != nil {
		updates["password"] = *req.Password
		user.Password = *req.Password
	}

	if req.Type != nil {
		updates["type"] = *req.Type
		user.Type = *req.Type
	}

	if req.DefaultWalletID != nil {
		updates["default_wallet_id"] = *req.DefaultWalletID
		user.DefaultWalletID = req.DefaultWalletID
	}

	if len(updates) == 0 {
		return user, nil // No updates provided
	}

	if err := database.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	log.Printf("✓ User '%s' updated", user.Username)
	return user, nil
}

// DeleteUser deletes a user by ID
func DeleteUser(userID uint) error {
	// Check if user exists
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	// Delete the user
	if err := database.DB.Delete(&models.User{}, userID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	log.Printf("✓ User '%s' (ID: %d) deleted", user.Username, userID)
	return nil
}

// capitalizeUsername converts username to a capitalized name
// e.g., "john_doe" -> "John Doe"
func capitalizeUsername(username string) string {
	if username == "" {
		return ""
	}

	// Replace underscores and hyphens with spaces
	result := ""
	capitalizeNext := true

	for _, char := range username {
		if char == '_' || char == '-' {
			result += " "
			capitalizeNext = true
		} else if capitalizeNext {
			result += string(char - 32) // Convert to uppercase
			capitalizeNext = false
		} else {
			result += string(char)
		}
	}

	return result
}
