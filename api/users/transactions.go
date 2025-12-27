package users

import (
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
	"time"
)

// CreateDummyTransactions creates dummy transactions (amount 0) for all categories
// This is used to satisfy foreign key constraints
func CreateDummyTransactions(walletID uint, userID uint, categoryIDs []uint) ([]models.Transaction, error) {
	var transactions []models.Transaction

	for _, categoryID := range categoryIDs {
		dummyTxn := &models.Transaction{
			CategoryID:       categoryID,
			Amount:           0,
			WalletID:         walletID,
			UserID:           userID,
			TransactionTime:  time.Now(),
			EntryTime:        time.Now(),
			LastModifiedTime: time.Now(),
		}
		if err := database.DB.Create(dummyTxn).Error; err != nil {
			log.Printf("Warning: Failed to create dummy transaction for category %d: %v", categoryID, err)
		} else {
			log.Printf("✓ Dummy transaction created for category %d", categoryID)
			transactions = append(transactions, *dummyTxn)
		}
	}

	return transactions, nil
}

// GetAllCategoryIDsForWallet retrieves all category IDs (root and subcategories) for a wallet
func GetAllCategoryIDsForWallet(walletID uint, rootCategoryIDs []uint) ([]uint, error) {
	var categoryIDs []uint

	// Add root categories
	categoryIDs = append(categoryIDs, rootCategoryIDs...)

	// Get all subcategories
	var subcategories []models.Category
	if err := database.DB.Where("parent_id IS NOT NULL AND wallet_id = ?", walletID).Find(&subcategories).Error; err != nil {
		return nil, err
	}

	for _, subcat := range subcategories {
		categoryIDs = append(categoryIDs, subcat.CategoryID)
	}

	return categoryIDs, nil
}
