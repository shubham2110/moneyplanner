package transactions

import (
	"fmt"
	"log"
	"moneyplanner/api/persons"
	"moneyplanner/database"
	"moneyplanner/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// adjustWalletBalance adjusts the wallet balance based on category root ID
func adjustWalletBalance(walletID uint, rootID uint, amount float64) error {
	var adjustment float64
	if rootID == 1 {
		adjustment = amount
	} else if rootID == 2 {
		adjustment = -amount
	} else {
		return nil // No adjustment for other root categories
	}

	if err := database.DB.Model(&models.Wallet{}).Where("wallet_id = ?", walletID).Update("balance", gorm.Expr("balance + ?", adjustment)).Error; err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}
	return nil
}

// GetTransactionByID retrieves a transaction by its ID
func GetTransactionByID(transactionID uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := database.DB.Preload("Category.Wallet").Preload("Person").Preload("Wallet").Preload("User").First(&transaction, transactionID).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return &transaction, nil
}

// ListAllTransactions retrieves all transactions with optional filters
func ListAllTransactions(filter *TransactionFilter) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := database.DB.Preload("Category.Wallet").Preload("Person").Preload("Wallet").Preload("User")

	// Apply filters
	if filter != nil {
		if filter.StartTransactionTime != nil && filter.EndTransactionTime != nil {
			query = query.Where("transaction_time BETWEEN ? AND ?", *filter.StartTransactionTime, *filter.EndTransactionTime)
		} else if filter.StartTransactionTime != nil {
			query = query.Where("transaction_time >= ?", *filter.StartTransactionTime)
		} else if filter.EndTransactionTime != nil {
			query = query.Where("transaction_time <= ?", *filter.EndTransactionTime)
		}

		if filter.StartEntryTime != nil && filter.EndEntryTime != nil {
			query = query.Where("entry_time BETWEEN ? AND ?", *filter.StartEntryTime, *filter.EndEntryTime)
		} else if filter.StartEntryTime != nil {
			query = query.Where("entry_time >= ?", *filter.StartEntryTime)
		} else if filter.EndEntryTime != nil {
			query = query.Where("entry_time <= ?", *filter.EndEntryTime)
		}

		if filter.StartLastModifiedTime != nil && filter.EndLastModifiedTime != nil {
			query = query.Where("last_modified_time BETWEEN ? AND ?", *filter.StartLastModifiedTime, *filter.EndLastModifiedTime)
		} else if filter.StartLastModifiedTime != nil {
			query = query.Where("last_modified_time >= ?", *filter.StartLastModifiedTime)
		} else if filter.EndLastModifiedTime != nil {
			query = query.Where("last_modified_time <= ?", *filter.EndLastModifiedTime)
		}

		if filter.UserID != nil {
			query = query.Where("user_id = ?", *filter.UserID)
		}

		if len(filter.CategoryIDs) > 0 {
			query = query.Where("category_id IN ?", filter.CategoryIDs)
		}

		if filter.WalletID != nil {
			query = query.Where("wallet_id = ?", *filter.WalletID)
		}

		if filter.PersonID != nil {
			query = query.Where("person_id = ?", *filter.PersonID)
		}

		if filter.FuzzyNote != nil {
			query = query.Where("note LIKE ?", "%"+*filter.FuzzyNote+"%")
		}

		if filter.AmountOp != nil && filter.AmountValue != nil {
			switch *filter.AmountOp {
			case "eq":
				query = query.Where("amount = ?", *filter.AmountValue)
			case "lt":
				query = query.Where("amount < ?", *filter.AmountValue)
			case "le":
				query = query.Where("amount <= ?", *filter.AmountValue)
			case "gt":
				query = query.Where("amount > ?", *filter.AmountValue)
			case "ge":
				query = query.Where("amount >= ?", *filter.AmountValue)
			}
		}
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	return transactions, nil
}

// ListTransactionsByUser retrieves transactions for a user
func ListTransactionsByUser(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := database.DB.Preload("Category.Wallet").Preload("Person").Preload("Wallet").Preload("User").Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to list transactions for user: %w", err)
	}
	return transactions, nil
}

// ListTransactionsByWallet retrieves transactions for a wallet with optional additional filters
func ListTransactionsByWallet(walletID uint, filter *TransactionFilter) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := database.DB.Preload("Category.Wallet").Preload("Person").Preload("Wallet").Preload("User").Where("wallet_id = ?", walletID)

	// Apply additional filters
	if filter != nil {
		if filter.StartTransactionTime != nil && filter.EndTransactionTime != nil {
			query = query.Where("transaction_time BETWEEN ? AND ?", *filter.StartTransactionTime, *filter.EndTransactionTime)
		} else if filter.StartTransactionTime != nil {
			query = query.Where("transaction_time >= ?", *filter.StartTransactionTime)
		} else if filter.EndTransactionTime != nil {
			query = query.Where("transaction_time <= ?", *filter.EndTransactionTime)
		}

		if filter.StartEntryTime != nil && filter.EndEntryTime != nil {
			query = query.Where("entry_time BETWEEN ? AND ?", *filter.StartEntryTime, *filter.EndEntryTime)
		} else if filter.StartEntryTime != nil {
			query = query.Where("entry_time >= ?", *filter.StartEntryTime)
		} else if filter.EndEntryTime != nil {
			query = query.Where("entry_time <= ?", *filter.EndEntryTime)
		}

		if filter.StartLastModifiedTime != nil && filter.EndLastModifiedTime != nil {
			query = query.Where("last_modified_time BETWEEN ? AND ?", *filter.StartLastModifiedTime, *filter.EndLastModifiedTime)
		} else if filter.StartLastModifiedTime != nil {
			query = query.Where("last_modified_time >= ?", *filter.StartLastModifiedTime)
		} else if filter.EndLastModifiedTime != nil {
			query = query.Where("last_modified_time <= ?", *filter.EndLastModifiedTime)
		}

		if filter.UserID != nil {
			query = query.Where("user_id = ?", *filter.UserID)
		}

		if len(filter.CategoryIDs) > 0 {
			query = query.Where("category_id IN ?", filter.CategoryIDs)
		}

		if filter.PersonID != nil {
			query = query.Where("person_id = ?", *filter.PersonID)
		}

		if filter.FuzzyNote != nil {
			query = query.Where("note LIKE ?", "%"+*filter.FuzzyNote+"%")
		}

		if filter.AmountOp != nil && filter.AmountValue != nil {
			switch *filter.AmountOp {
			case "eq":
				query = query.Where("amount = ?", *filter.AmountValue)
			case "lt":
				query = query.Where("amount < ?", *filter.AmountValue)
			case "le":
				query = query.Where("amount <= ?", *filter.AmountValue)
			case "gt":
				query = query.Where("amount > ?", *filter.AmountValue)
			case "ge":
				query = query.Where("amount >= ?", *filter.AmountValue)
			}
		}
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to list transactions for wallet: %w", err)
	}
	return transactions, nil
}

// CreateTransaction creates a new transaction
func CreateTransaction(req *TransactionCreationRequest) (*models.Transaction, error) {
	if req.WalletID == 0 {
		return nil, fmt.Errorf("wallet_id is required")
	}
	if req.CategoryID == 0 {
		return nil, fmt.Errorf("category_id is required")
	}
	if req.Amount == 0 {
		return nil, fmt.Errorf("amount is required")
	}
	if req.UserID == 0 {
		return nil, fmt.Errorf("user_id is required")
	}

	// Handle person
	var personID *uint
	if req.PersonName != nil && strings.TrimSpace(*req.PersonName) != "" {
		// Try to find existing person
		p, err := persons.GetPersonByName(*req.PersonName)
		if err != nil {
			// Create new person
			newPersonReq := &persons.PersonCreationRequest{
				PersonName: *req.PersonName,
				Alias:      "", // No alias provided
			}
			newPerson, err := persons.CreatePerson(newPersonReq)
			if err != nil {
				return nil, fmt.Errorf("failed to create person: %w", err)
			}
			personID = &newPerson.PersonID
		} else {
			personID = &p.PersonID
		}
	} else if req.PersonID != nil {
		personID = req.PersonID
	}

	now := time.Now()
	transactionTime := now
	if req.TransactionTime != nil {
		transactionTime = *req.TransactionTime
	}

	t := &models.Transaction{
		WalletID:         req.WalletID,
		CategoryID:       req.CategoryID,
		Amount:           req.Amount,
		PersonID:         personID,
		Note:             req.Note,
		TransactionTime:  transactionTime,
		EntryTime:        now,
		LastModifiedTime: now,
		UserID:           req.UserID,
	}

	if err := database.DB.Create(t).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Load relationships
	if err := database.DB.Preload("Category.Wallet").Preload("Person").Preload("Wallet").Preload("User").First(t, t.TransactionID).Error; err != nil {
		log.Printf("Warning: Failed to preload transaction relationships: %v", err)
	}

	// Adjust wallet balance
	if err := adjustWalletBalance(t.WalletID, t.Category.RootID, t.Amount); err != nil {
		log.Printf("Warning: Failed to adjust wallet balance: %v", err)
	}

	log.Printf("✓ Transaction created (ID: %d)", t.TransactionID)
	return t, nil
}

// UpdateTransaction updates transaction details
func UpdateTransaction(transactionID uint, req *TransactionUpdateRequest) (*models.Transaction, error) {
	transaction, err := GetTransactionByID(transactionID)
	if err != nil {
		return nil, err
	}

	// Save old values for balance adjustment
	oldRootID := transaction.Category.RootID
	oldAmount := transaction.Amount
	oldWalletID := transaction.WalletID

	updates := map[string]interface{}{}

	if req.WalletID != nil {
		updates["wallet_id"] = *req.WalletID
		transaction.WalletID = *req.WalletID
	}

	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
		transaction.CategoryID = *req.CategoryID
	}

	if req.Amount != nil {
		updates["amount"] = *req.Amount
		transaction.Amount = *req.Amount
	}

	if req.Note != nil {
		updates["note"] = *req.Note
		transaction.Note = req.Note
	}

	if req.TransactionTime != nil {
		updates["transaction_time"] = *req.TransactionTime
		transaction.TransactionTime = *req.TransactionTime
	}

	// Handle person update
	if req.PersonName != nil && strings.TrimSpace(*req.PersonName) != "" {
		p, err := persons.GetPersonByName(*req.PersonName)
		if err != nil {
			newPersonReq := &persons.PersonCreationRequest{
				PersonName: *req.PersonName,
				Alias:      "",
			}
			newPerson, err := persons.CreatePerson(newPersonReq)
			if err != nil {
				return nil, fmt.Errorf("failed to create person: %w", err)
			}
			updates["person_id"] = newPerson.PersonID
			transaction.PersonID = &newPerson.PersonID
		} else {
			updates["person_id"] = p.PersonID
			transaction.PersonID = &p.PersonID
		}
	} else if req.PersonID != nil {
		updates["person_id"] = *req.PersonID
		transaction.PersonID = req.PersonID
	}

	updates["last_modified_time"] = time.Now()

	if len(updates) == 0 {
		return transaction, nil // No updates provided
	}

	if err := database.DB.Model(&models.Transaction{}).
		Where("transaction_id = ?", transactionID).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Reverse old effect on balance
	if oldRootID == 1 {
		if err := adjustWalletBalance(oldWalletID, oldRootID, -oldAmount); err != nil {
			log.Printf("Warning: Failed to reverse wallet balance: %v", err)
		}
	} else if oldRootID == 2 {
		if err := adjustWalletBalance(oldWalletID, oldRootID, oldAmount); err != nil {
			log.Printf("Warning: Failed to reverse wallet balance: %v", err)
		}
	}

	// Reload relationships
	if err := database.DB.Preload("Category.Wallet").Preload("Person").Preload("Wallet").Preload("User").First(transaction, transactionID).Error; err != nil {
		log.Printf("Warning: Failed to preload transaction relationships: %v", err)
	}

	// Apply new effect on balance
	if transaction.Category.RootID == 1 {
		if err := adjustWalletBalance(transaction.WalletID, transaction.Category.RootID, transaction.Amount); err != nil {
			log.Printf("Warning: Failed to adjust new wallet balance: %v", err)
		}
	} else if transaction.Category.RootID == 2 {
		if err := adjustWalletBalance(transaction.WalletID, transaction.Category.RootID, -transaction.Amount); err != nil {
			log.Printf("Warning: Failed to adjust new wallet balance: %v", err)
		}
	}

	log.Printf("✓ Transaction (ID: %d) updated", transactionID)
	return transaction, nil
}

// DeleteTransaction deletes a transaction by ID
func DeleteTransaction(transactionID uint) error {
	// Check if transaction exists
	transaction, err := GetTransactionByID(transactionID)
	if err != nil {
		return err
	}

	// Reverse the effect on balance
	if transaction.Category.RootID == 1 {
		if err := adjustWalletBalance(transaction.WalletID, transaction.Category.RootID, -transaction.Amount); err != nil {
			log.Printf("Warning: Failed to reverse wallet balance on delete: %v", err)
		}
	} else if transaction.Category.RootID == 2 {
		if err := adjustWalletBalance(transaction.WalletID, transaction.Category.RootID, transaction.Amount); err != nil {
			log.Printf("Warning: Failed to reverse wallet balance on delete: %v", err)
		}
	}

	// Delete the transaction
	if err := database.DB.Delete(&models.Transaction{}, transactionID).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	log.Printf("✓ Transaction (ID: %d) deleted", transactionID)
	return nil
}
