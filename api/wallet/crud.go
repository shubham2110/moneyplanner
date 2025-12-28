package wallet

import (
	"fmt"
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
	"time"
)

// GetWalletByID retrieves a wallet by its ID
func GetWalletByID(walletID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := database.DB.First(&wallet, walletID).Error; err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}
	return &wallet, nil
}

// ListAllWallets retrieves all wallets from the database
func ListAllWallets() ([]models.Wallet, error) {
	var wallets []models.Wallet
	if err := database.DB.Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to list wallets: %w", err)
	}
	return wallets, nil
}

// CreateWallet creates a new wallet
func CreateWallet(req *WalletCreationRequest) (*models.Wallet, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	w := &models.Wallet{
		Name:             req.Name,
		Icon:             req.Icon,
		IsEnabled:        true,
		Balance:          0,
		LastModifiedTime: time.Now(),
	}

	// Optional overrides
	if req.IsEnabled != nil {
		w.IsEnabled = *req.IsEnabled
	}
	if req.Balance != nil {
		w.Balance = *req.Balance
	}

	if err := database.DB.Create(w).Error; err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	log.Printf("✓ Wallet '%s' created (ID: %d)", w.Name, w.WalletID)
	return w, nil
}

// UpdateWallet updates wallet details
func UpdateWallet(walletID uint, req *WalletUpdateRequest) (*models.Wallet, error) {
	wallet, err := GetWalletByID(walletID)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = *req.Name
		wallet.Name = *req.Name
	}

	if req.Icon != nil {
		updates["icon"] = *req.Icon
		wallet.Icon = *req.Icon
	}

	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
		wallet.IsEnabled = *req.IsEnabled
	}

	if req.Balance != nil {
		updates["balance"] = *req.Balance
		wallet.Balance = *req.Balance
	}

	// Only touch LastModifiedTime if we actually update something
	if len(updates) == 0 {
		return wallet, nil // No updates provided
	}
	updates["last_modified_time"] = time.Now()
	wallet.LastModifiedTime = updates["last_modified_time"].(time.Time)

	if err := database.DB.Model(&models.Wallet{}).
		Where("wallet_id = ?", walletID).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update wallet: %w", err)
	}

	log.Printf("✓ Wallet '%s' (ID: %d) updated", wallet.Name, walletID)
	return wallet, nil
}

// DeleteWallet deletes a wallet by ID
func DeleteWallet(walletID uint) error {
	// Check if wallet exists
	wallet, err := GetWalletByID(walletID)
	if err != nil {
		return err
	}

	// Delete the wallet
	if err := database.DB.Delete(&models.Wallet{}, walletID).Error; err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	log.Printf("✓ Wallet '%s' (ID: %d) deleted", wallet.Name, walletID)
	return nil
}
