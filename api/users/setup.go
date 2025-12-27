package users

import (
	"log"
	"errors"
	"fmt"
	"moneyplanner/models"
	"time"
	"moneyplanner/database"
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


// CreateWalletGroup creates a new wallet group
func CreateWalletGroup(name string) (*models.WalletGroup, error) {
	walletGroup := &models.WalletGroup{
		WalletGroupName: name,
	}
	if err := database.DB.Create(walletGroup).Error; err != nil {
		return nil, err
	}
	log.Printf("✓ Wallet group '%s' created", name)
	return walletGroup, nil
}

// CreateWallet creates a new wallet with given parameters
func CreateWallet(name, icon string, initialBalance float64) (*models.Wallet, error) {
	wallet := &models.Wallet{
		Name:             name,
		Icon:             icon,
		IsEnabled:        true,
		Balance:          initialBalance,
		LastModifiedTime: time.Now(),
	}
	if err := database.DB.Create(wallet).Error; err != nil {
		return nil, err
	}
	log.Printf("✓ Wallet '%s' created", name)
	return wallet, nil
}

// AssociateWalletWithGroup associates a wallet with a wallet group
func AssociateWalletWithGroup(wallet *models.Wallet, walletGroup *models.WalletGroup) error {
	database.DB.Model(wallet).Association("WalletGroups").Append(walletGroup)
	return nil
}

// AssociateUserWithWallet associates a user with a wallet
func AssociateUserWithWallet(user *models.User, wallet *models.Wallet) error {
	database.DB.Model(user).Association("Wallets").Append(wallet)
	return nil
}



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


// CreateRootCategories creates the default Income and Expense root categories
func CreateRootCategories(walletID uint) ([]models.Category, error) {
	var rootCategories []models.Category

	incomeCategory := &models.Category{
		Icon:     "💵",
		Name:     "Income",
		WalletID: walletID,
		IsGlobal: true,
		RootID:   0, // Will be set after creation
	}
	if err := database.DB.Create(incomeCategory).Error; err != nil {
		return nil, err
	}
	incomeCategory.RootID = incomeCategory.CategoryID
	database.DB.Save(incomeCategory)
	log.Println("✓ Income root category created")
	rootCategories = append(rootCategories, *incomeCategory)

	expenseCategory := &models.Category{
		Icon:     "💸",
		Name:     "Expense",
		WalletID: walletID,
		IsGlobal: true,
		RootID:   0,
	}
	if err := database.DB.Create(expenseCategory).Error; err != nil {
		return nil, err
	}
	expenseCategory.RootID = expenseCategory.CategoryID
	database.DB.Save(expenseCategory)
	log.Println("✓ Expense root category created")
	rootCategories = append(rootCategories, *expenseCategory)

	return rootCategories, nil
}

// CreateIncomeSubcategories creates default income subcategories
func CreateIncomeSubcategories(walletID uint, parentCategoryID uint) error {
	incomeSubcategories := []string{"Salary", "Refund", "Bonus", "Interest"}
	for _, catName := range incomeSubcategories {
		subcat := &models.Category{
			Name:     catName,
			Icon:     "📊",
			ParentID: &parentCategoryID,
			RootID:   parentCategoryID,
			WalletID: walletID,
			IsGlobal: false,
		}
		if err := database.DB.Create(subcat).Error; err != nil {
			log.Printf("Warning: Failed to create %s subcategory: %v", catName, err)
		} else {
			log.Printf("✓ %s subcategory created", catName)
		}
	}
	return nil
}

// CreateExpenseSubcategories creates default expense subcategories
func CreateExpenseSubcategories(walletID uint, parentCategoryID uint) error {
	expenseSubcategories := []string{"Groceries", "House Maintenance", "Investment", "Utilities", "Transport", "Entertainment"}
	for _, catName := range expenseSubcategories {
		subcat := &models.Category{
			Name:     catName,
			Icon:     "📉",
			ParentID: &parentCategoryID,
			RootID:   parentCategoryID,
			WalletID: walletID,
			IsGlobal: false,
		}
		if err := database.DB.Create(subcat).Error; err != nil {
			log.Printf("Warning: Failed to create %s subcategory: %v", catName, err)
		} else {
			log.Printf("✓ %s subcategory created", catName)
		}
	}
	return nil
}

// CreateAllSubcategories creates both income and expense subcategories
func CreateAllSubcategories(walletID uint, incomeCategoryID uint, expenseCategoryID uint) error {
	if err := CreateIncomeSubcategories(walletID, incomeCategoryID); err != nil {
		return err
	}
	if err := CreateExpenseSubcategories(walletID, expenseCategoryID); err != nil {
		return err
	}
	return nil
}
