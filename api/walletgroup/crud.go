package walletgroup

import (
	"fmt"
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
)

// GetWalletGroupByID retrieves a wallet group by ID
func GetWalletGroupByID(walletGroupID uint) (*models.WalletGroup, error) {
	var wg models.WalletGroup
	if err := database.DB.First(&wg, walletGroupID).Error; err != nil {
		return nil, fmt.Errorf("wallet group not found: %w", err)
	}
	return &wg, nil
}

// ListAllWalletGroups lists all wallet groups
func ListAllWalletGroups() ([]models.WalletGroup, error) {
	var groups []models.WalletGroup
	if err := database.DB.Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("failed to list wallet groups: %w", err)
	}
	return groups, nil
}

// CreateWalletGroup creates a wallet group
func CreateWalletGroup(req *WalletGroupCreationRequest) (*models.WalletGroup, error) {
	if req.WalletGroupName == "" {
		return nil, fmt.Errorf("wallet_group_name is required")
	}

	wg := &models.WalletGroup{
		WalletGroupName: req.WalletGroupName,
	}

	if err := database.DB.Create(wg).Error; err != nil {
		return nil, fmt.Errorf("failed to create wallet group: %w", err)
	}

	log.Printf("✓ WalletGroup '%s' created (ID: %d)", wg.WalletGroupName, wg.WalletGroupID)
	return wg, nil
}

// UpdateWalletGroup updates wallet group details
func UpdateWalletGroup(walletGroupID uint, req *WalletGroupUpdateRequest) (*models.WalletGroup, error) {
	wg, err := GetWalletGroupByID(walletGroupID)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}

	if req.WalletGroupName != nil {
		updates["wallet_group_name"] = *req.WalletGroupName
		wg.WalletGroupName = *req.WalletGroupName
	}

	if len(updates) == 0 {
		return wg, nil
	}

	if err := database.DB.Model(&models.WalletGroup{}).
		Where("wallet_group_id = ?", walletGroupID).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update wallet group: %w", err)
	}

	log.Printf("✓ WalletGroup '%s' (ID: %d) updated", wg.WalletGroupName, walletGroupID)
	return wg, nil
}

// DeleteWalletGroup deletes a wallet group
func DeleteWalletGroup(walletGroupID uint) error {
	wg, err := GetWalletGroupByID(walletGroupID)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(&models.WalletGroup{}, walletGroupID).Error; err != nil {
		return fmt.Errorf("failed to delete wallet group: %w", err)
	}

	log.Printf("✓ WalletGroup '%s' (ID: %d) deleted", wg.WalletGroupName, walletGroupID)
	return nil
}
