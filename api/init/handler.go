package init

import (
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
	"time"
)

// InitRequest represents the initialization request payload
type InitRequest struct {
	ForceMigrate        bool   `json:"force_migrate"`
	DefaultWalletName   string `json:"default_wallet_name,omitempty"`
	DefaultWalletGroup  string `json:"default_wallet_group,omitempty"`
	AdminUsername       string `json:"admin_username,omitempty"`
	AdminPassword       string `json:"admin_password,omitempty"`
	AdminEmail          string `json:"admin_email,omitempty"`
	AdminName           string `json:"admin_name,omitempty"`
}

// InitResponse represents the API response
type InitResponse struct {
	Success         bool        `json:"success"`
	Message         string      `json:"message"`
	Data            *InitedData `json:"data,omitempty"`
	Error           string      `json:"error,omitempty"`
	MissingItems    []string    `json:"missing_items,omitempty"` // List of missing tables/columns
}

// InitedData contains the created default entities
type InitedData struct {
	DatabaseIsNew      bool                `json:"database_is_new"`
	MigrationRequired  bool                `json:"migration_required"`
	AdminUser          *models.User        `json:"admin_user,omitempty"`
	DefaultWallet      *models.Wallet      `json:"default_wallet,omitempty"`
	DefaultWalletGroup *models.WalletGroup `json:"default_wallet_group,omitempty"`
	RootCategories     []models.Category   `json:"root_categories,omitempty"`
}

// InitializeDatabase handles the full database initialization with migrations and default data
func InitializeDatabase(req *InitRequest, dbPath string) *InitResponse {
	log.Println("Starting database initialization...")

	// Connect to database
	if err := database.InitDB(dbPath); err != nil {
		return &InitResponse{
			Success: false,
			Message: "Failed to connect to database",
			Error:   err.Error(),
		}
	}

	isNewDB := database.IsNewDatabase(dbPath)
	initedData := &InitedData{
		DatabaseIsNew: isNewDB,
	}

	// Handle existing database
	if !isNewDB {
		migrationNeeded, missingItems := database.IsMigrationNeeded()
		initedData.MigrationRequired = migrationNeeded

		if migrationNeeded && !req.ForceMigrate {
			return &InitResponse{
				Success:      false,
				Message:      "Database schema does not match ORM models",
				Error:        "Pending migrations detected. Set force_migrate: true to proceed",
				Data:         initedData,
				MissingItems: missingItems,
			}
		}

		if migrationNeeded || req.ForceMigrate {
			log.Println("Running migrations...")
			if err := database.MigrateDB(); err != nil {
				return &InitResponse{
					Success: false,
					Message: "Migration failed",
					Error:   err.Error(),
				}
			}
			log.Println("✓ Migrations completed")
		}

		return &InitResponse{
			Success: true,
			Message: "Database schema is in sync",
			Data:    initedData,
		}
	}

	// New database - run migrations and create defaults
	log.Println("Creating new database...")
	if err := database.MigrateDB(); err != nil {
		return &InitResponse{
			Success: false,
			Message: "Failed to create database tables",
			Error:   err.Error(),
		}
	}
	log.Println("✓ Database tables created")

	// Set defaults
	walletGroupName := "Default"
	if req.DefaultWalletGroup != "" {
		walletGroupName = req.DefaultWalletGroup
	}

	walletName := "Default Wallet"
	if req.DefaultWalletName != "" {
		walletName = req.DefaultWalletName
	}

	adminUsername := "admin"
	if req.AdminUsername != "" {
		adminUsername = req.AdminUsername
	}

	adminPassword := "admin"
	if req.AdminPassword != "" {
		adminPassword = req.AdminPassword
	}

	adminEmail := "admin@moneyplanner.local"
	if req.AdminEmail != "" {
		adminEmail = req.AdminEmail
	}

	adminName := "Admin User"
	if req.AdminName != "" {
		adminName = req.AdminName
	}

	// Create default wallet group
	walletGroup := &models.WalletGroup{
		WalletGroupName: walletGroupName,
	}
	if err := database.DB.Create(walletGroup).Error; err != nil {
		return &InitResponse{
			Success: false,
			Error:   "Failed to create wallet group: " + err.Error(),
		}
	}
	log.Println("✓ Default wallet group created")
	initedData.DefaultWalletGroup = walletGroup

	// Create default wallet
	wallet := &models.Wallet{
		Name:             walletName,
		Icon:             "💰",
		IsEnabled:        true,
		Balance:          0,
		LastModifiedTime: time.Now(),
	}
	if err := database.DB.Create(wallet).Error; err != nil {
		return &InitResponse{
			Success: false,
			Error:   "Failed to create wallet: " + err.Error(),
		}
	}
	log.Println("✓ Default wallet created")
	initedData.DefaultWallet = wallet

	// Associate wallet with group
	database.DB.Model(wallet).Association("WalletGroups").Append(walletGroup)

	// Create default admin user
	adminUser := &models.User{
		Username:       adminUsername,
		Name:           adminName,
		Email:          adminEmail,
		Password:       adminPassword,
		Type:           models.UserTypeHuman,
		DefaultWalletID: &wallet.WalletID,
	}
	if err := database.DB.Create(adminUser).Error; err != nil {
		return &InitResponse{
			Success: false,
			Error:   "Failed to create admin user: " + err.Error(),
		}
	}
	log.Println("✓ Default admin user created")
	initedData.AdminUser = adminUser

	// Associate user with wallet
	database.DB.Model(adminUser).Association("Wallets").Append(wallet)

	// Create root categories
	incomeCategory := &models.Category{
		Icon:     "💵",
		Name:     "Income",
		WalletID: wallet.WalletID,
		IsGlobal: true,
		RootID:   0, // Will be set after creation
	}
	if err := database.DB.Create(incomeCategory).Error; err != nil {
		return &InitResponse{
			Success: false,
			Error:   "Failed to create Income category: " + err.Error(),
		}
	}
	incomeCategory.RootID = incomeCategory.CategoryID
	database.DB.Save(incomeCategory)
	log.Println("✓ Income root category created")
	initedData.RootCategories = append(initedData.RootCategories, *incomeCategory)

	expenseCategory := &models.Category{
		Icon:     "💸",
		Name:     "Expense",
		WalletID: wallet.WalletID,
		IsGlobal: true,
		RootID:   0,
	}
	if err := database.DB.Create(expenseCategory).Error; err != nil {
		return &InitResponse{
			Success: false,
			Error:   "Failed to create Expense category: " + err.Error(),
		}
	}
	expenseCategory.RootID = expenseCategory.CategoryID
	database.DB.Save(expenseCategory)
	log.Println("✓ Expense root category created")
	initedData.RootCategories = append(initedData.RootCategories, *expenseCategory)

	// Create income subcategories
	incomeSubcategories := []string{"Salary", "Refund", "Bonus", "Interest"}
	for _, catName := range incomeSubcategories {
		subcat := &models.Category{
			Name:     catName,
			Icon:     "📊",
			ParentID: &incomeCategory.CategoryID,
			RootID:   incomeCategory.CategoryID,
			WalletID: wallet.WalletID,
			IsGlobal: false,
		}
		if err := database.DB.Create(subcat).Error; err != nil {
			log.Printf("Warning: Failed to create %s subcategory: %v", catName, err)
		} else {
			log.Printf("✓ %s subcategory created", catName)
		}
	}

	// Create expense subcategories
	expenseSubcategories := []string{"Groceries", "House Maintenance", "Investment", "Utilities", "Transport", "Entertainment"}
	for _, catName := range expenseSubcategories {
		subcat := &models.Category{
			Name:     catName,
			Icon:     "📉",
			ParentID: &expenseCategory.CategoryID,
			RootID:   expenseCategory.CategoryID,
			WalletID: wallet.WalletID,
			IsGlobal: false,
		}
		if err := database.DB.Create(subcat).Error; err != nil {
			log.Printf("Warning: Failed to create %s subcategory: %v", catName, err)
		} else {
			log.Printf("✓ %s subcategory created", catName)
		}
	}

	// Create dummy transactions (amount 0) for all categories to satisfy foreign key constraints
	log.Println("Creating dummy transactions for categories...")
	allCategories := []uint{
		incomeCategory.CategoryID,
		expenseCategory.CategoryID,
	}

	// Get all subcategories
	var subcategories []models.Category
	database.DB.Where("parent_id IS NOT NULL").Find(&subcategories)
	for _, subcat := range subcategories {
		allCategories = append(allCategories, subcat.CategoryID)
	}

	// Create dummy transaction for each category
	for _, categoryID := range allCategories {
		dummyTxn := &models.Transaction{
			CategoryID:       categoryID,
			Amount:           0,
			WalletID:         wallet.WalletID,
			UserID:           adminUser.UserID,
			TransactionTime:  time.Now(),
			EntryTime:        time.Now(),
			LastModifiedTime: time.Now(),
		}
		if err := database.DB.Create(dummyTxn).Error; err != nil {
			log.Printf("Warning: Failed to create dummy transaction for category %d: %v", categoryID, err)
		} else {
			log.Printf("✓ Dummy transaction created for category %d", categoryID)
		}
	}

	return &InitResponse{
		Success: true,
		Message: "Database initialized successfully",
		Data:    initedData,
	}
}
