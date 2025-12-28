package walletgroupwallet

import (
	"fmt"
	"log"

	"moneyplanner/database"
	"moneyplanner/models"
)

// AttachWalletToGroup adds wallet -> walletgroup
func AttachWalletToGroup(walletGroupID, walletID uint) error {
	var wg models.WalletGroup
	if err := database.DB.First(&wg, walletGroupID).Error; err != nil {
		return fmt.Errorf("wallet group not found: %w", err)
	}

	var w models.Wallet
	if err := database.DB.First(&w, walletID).Error; err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}

	if err := database.DB.Model(&wg).Association("Wallets").Append(&w); err != nil {
		return fmt.Errorf("failed to attach wallet to wallet group: %w", err)
	}

	log.Printf("✓ Attached wallet %d to wallet group %d", walletID, walletGroupID)
	return nil
}

// DetachWalletFromGroup removes wallet -> walletgroup
func DetachWalletFromGroup(walletGroupID, walletID uint) error {
	var wg models.WalletGroup
	if err := database.DB.First(&wg, walletGroupID).Error; err != nil {
		return fmt.Errorf("wallet group not found: %w", err)
	}

	var w models.Wallet
	if err := database.DB.First(&w, walletID).Error; err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}

	if err := database.DB.Model(&wg).Association("Wallets").Delete(&w); err != nil {
		return fmt.Errorf("failed to detach wallet from wallet group: %w", err)
	}

	log.Printf("✓ Detached wallet %d from wallet group %d", walletID, walletGroupID)
	return nil
}

// ReplaceWalletsInGroup sets group wallets to exactly walletIDs
func ReplaceWalletsInGroup(walletGroupID uint, walletIDs []uint) error {
	var wg models.WalletGroup
	if err := database.DB.First(&wg, walletGroupID).Error; err != nil {
		return fmt.Errorf("wallet group not found: %w", err)
	}

	if len(walletIDs) == 0 {
		if err := database.DB.Model(&wg).Association("Wallets").Clear(); err != nil {
			return fmt.Errorf("failed to clear wallets from wallet group: %w", err)
		}
		log.Printf("✓ Cleared wallets for wallet group %d", walletGroupID)
		return nil
	}

	var wallets []models.Wallet
	if err := database.DB.Where("wallet_id IN ?", walletIDs).Find(&wallets).Error; err != nil {
		return fmt.Errorf("failed to load wallets: %w", err)
	}
	if len(wallets) != len(walletIDs) {
		return fmt.Errorf("one or more wallet IDs not found")
	}

	if err := database.DB.Model(&wg).Association("Wallets").Replace(wallets); err != nil {
		return fmt.Errorf("failed to replace wallets in wallet group: %w", err)
	}

	log.Printf("✓ Replaced wallets for wallet group %d", walletGroupID)
	return nil
}

// ListWalletsInGroup lists wallets for a group
func ListWalletsInGroup(walletGroupID uint) ([]models.Wallet, error) {
	var wg models.WalletGroup
	if err := database.DB.Preload("Wallets").First(&wg, walletGroupID).Error; err != nil {
		return nil, fmt.Errorf("wallet group not found: %w", err)
	}
	return wg.Wallets, nil
}
