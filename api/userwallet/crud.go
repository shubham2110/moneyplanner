package userwallet

import (
	"fmt"
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
)

// AttachWalletToUser creates a row in user_wallets (idempotent-ish: no error if already attached)
func AttachWalletToUser(userID, walletID uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var wallet models.Wallet
	if err := database.DB.First(&wallet, walletID).Error; err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}

	// Append association
	if err := database.DB.Model(&user).Association("Wallets").Append(&wallet); err != nil {
		return fmt.Errorf("failed to attach wallet to user: %w", err)
	}

	log.Printf("✓ Attached wallet %d to user %d", walletID, userID)
	return nil
}

// DetachWalletFromUser removes a row from user_wallets
func DetachWalletFromUser(userID, walletID uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var wallet models.Wallet
	if err := database.DB.First(&wallet, walletID).Error; err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}

	if err := database.DB.Model(&user).Association("Wallets").Delete(&wallet); err != nil {
		return fmt.Errorf("failed to detach wallet from user: %w", err)
	}

	log.Printf("✓ Detached wallet %d from user %d", walletID, userID)
	return nil
}

// ReplaceUserWallets sets the user's wallets to exactly this list (attach multiple in one go)
func ReplaceUserWallets(userID uint, walletIDs []uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if len(walletIDs) == 0 {
		// Clear all wallets
		if err := database.DB.Model(&user).Association("Wallets").Clear(); err != nil {
			return fmt.Errorf("failed to clear wallets: %w", err)
		}
		log.Printf("✓ Cleared all wallets for user %d", userID)
		return nil
	}

	var wallets []models.Wallet
	if err := database.DB.Where("wallet_id IN ?", walletIDs).Find(&wallets).Error; err != nil {
		return fmt.Errorf("failed to load wallets: %w", err)
	}
	if len(wallets) != len(walletIDs) {
		return fmt.Errorf("one or more wallet IDs not found")
	}

	if err := database.DB.Model(&user).Association("Wallets").Replace(wallets); err != nil {
		return fmt.Errorf("failed to replace user wallets: %w", err)
	}

	log.Printf("✓ Replaced wallets for user %d", userID)
	return nil
}

// ListUserWallets fetches user's wallets
func ListUserWallets(userID uint) ([]models.Wallet, error) {
	var user models.User
	if err := database.DB.Preload("Wallets").First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user.Wallets, nil
}
