package init

import (
	"log"
	"moneyplanner/api/users"
	"moneyplanner/database"
	"moneyplanner/models"
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

// InitDoneResponse represents the response for initialization status check
type InitDoneResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	InitDone  bool   `json:"init_done"`
	IsNewDB   bool   `json:"is_new_db"`
	Error     string `json:"error,omitempty"`
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

	// Create admin user with default wallet, wallet group, and categories
	adminUserReq := &users.UserCreationRequest{
		Username:         adminUsername,
		Name:             adminName,
		Email:            adminEmail,
		Password:         adminPassword,
		UserType:         models.UserTypeHuman,
		WalletName:       walletName,
		WalletGroupName:  walletGroupName,
		CreateCategories: true,
	}

	setupResult, err := users.SetupUserWithDefaults(adminUserReq)
	if err != nil {
		return &InitResponse{
			Success: false,
			Error:   "Failed to setup admin user: " + err.Error(),
		}
	}

	initedData.AdminUser = setupResult.User
	initedData.DefaultWallet = setupResult.Wallet
	initedData.DefaultWalletGroup = setupResult.WalletGroup
	initedData.RootCategories = setupResult.RootCategories

	return &InitResponse{
		Success: true,
		Message: "Database initialized successfully",
		Data:    initedData,
	}
}

// CheckInitStatus checks if database initialization is complete
func CheckInitStatus(dbPath string) *InitDoneResponse {
	
	isNew := database.IsNewDatabase(dbPath)

	return &InitDoneResponse{
		Success:  true,
		Message:  "Initialization status retrieved",
		InitDone: !isNew,  // true if database is NOT new (i.e., initialized)
		IsNewDB:  isNew,   // true if database is new
	}
}
