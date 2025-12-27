package users

import (
	"log"
	"errors"
	"fmt"
	"moneyplanner/models"
)

// SetupUserWithDefaults creates a complete user setup including wallet, wallet group, and categories
// This function handles all the logic for creating a user with their default wallet and associated entities
func SetupUserWithDefaults(req *UserCreationRequest) (*UserCreationResult, error) {
	result := &UserCreationResult{}

	if req.Username == "" {
		return nil, errors.New("username is mandatory")
	}

	// Set defaults for optional fields
	name := req.Name
	if name == "" {
		defaultName := capitalizeUsername(req.Username)
		name = defaultName
	}

	email := req.Email
	if email == "" {
		defaultEmail := fmt.Sprintf("%s@moneyplanner.local", req.Username)
		email = defaultEmail
	}

	password := req.Password
	if password == "" {
		defaultPassword := "password123"
		password = defaultPassword
	}

	userType := req.UserType
	if userType == "" {
		defaultType := models.UserTypeHuman
		userType = defaultType
	}

	// Handle wallet assignment or creation
	walletName := req.WalletName
	if walletName == "" {
		defaultWalletName := fmt.Sprintf("%s's Wallet", name)
		walletName = defaultWalletName
	}

	walletGroupName := req.WalletGroupName
	if walletGroupName == "" {
		defaultWalletGroupName := fmt.Sprintf("%s's Group", name)
		walletGroupName = defaultWalletGroupName
	}

	// Step 1: Create wallet group
	walletGroup, err := CreateWalletGroup(walletGroupName)
	if err != nil {
		return nil, err
	}
	result.WalletGroup = walletGroup

	// Step 2: Create wallet
	wallet, err := CreateWallet(walletName, "💰", 0)
	if err != nil {
		return nil, err
	}
	result.Wallet = wallet

	// Step 3: Associate wallet with group
	if err := AssociateWalletWithGroup(wallet, walletGroup); err != nil {
		log.Printf("Warning: Failed to associate wallet with group: %v", err)
	}

	// Step 4: Create user
	user, err := CreateUser(req.Username, name, email, password, userType, &wallet.WalletID)
	if err != nil {
		return nil, err
	}
	result.User = user

	// Step 5: Associate user with wallet
	if err := AssociateUserWithWallet(user, wallet); err != nil {
		log.Printf("Warning: Failed to associate user with wallet: %v", err)
	}

	if !req.CreateCategories{
		fmt.Println("Pass the create categories")
	}

	// Step 6: Create categories if requested
	if req.CreateCategories {
		rootCategories, err := CreateRootCategories(wallet.WalletID)
		if err != nil {
			return nil, err
		}
		result.RootCategories = rootCategories

		// Get root category IDs
		var rootCategoryIDs []uint
		for _, cat := range rootCategories {
			rootCategoryIDs = append(rootCategoryIDs, cat.CategoryID)
		}

		// Create subcategories
		if len(rootCategories) >= 2 {
			incomeCategoryID := rootCategories[0].CategoryID
			expenseCategoryID := rootCategories[1].CategoryID
			if err := CreateAllSubcategories(wallet.WalletID, incomeCategoryID, expenseCategoryID); err != nil {
				return nil, err
			}
		}

		// Step 7: Create dummy transactions for all categories
		allCategoryIDs, err := GetAllCategoryIDsForWallet(wallet.WalletID, rootCategoryIDs)
		if err != nil {
			return nil, err
		}

		dummyTxns, err := CreateDummyTransactions(wallet.WalletID, user.UserID, allCategoryIDs)
		if err != nil {
			return nil, err
		}
		result.DummyTransactions = dummyTxns
	}

	log.Printf("✓ User setup completed for '%s'", req.Username)
	return result, nil
}
