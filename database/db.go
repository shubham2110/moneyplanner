package database

import (
	"fmt"
	"log"
	"moneyplanner/models"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitOptions holds options for database initialization
type InitOptions struct {
	ForceMigrate bool
	DBPath       string
}

// InitDB initializes the database connection (does NOT run migrations)
func InitDB(dbPath string) error {
	var err error

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return nil
}

// IsNewDatabase checks if database file exists and has no tables
func IsNewDatabase(dbPath string) bool {
	// First check if file exists
	_, err := os.Stat(dbPath)
	if err != nil {
		// File doesn't exist, it's new
		return true
	}

	if DB != nil {
		return !DB.Migrator().HasTable("wallet_groups")
	}
	
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return true
	}

	return !DB.Migrator().HasTable("wallet_groups")
	// File exists, check if any tables exist by looking for a key table
	// If wallet_groups table doesn't exist, we consider it a new database
	
}

// IsMigrationNeeded checks if there are pending migrations and returns list of missing items
func IsMigrationNeeded() (needed bool, missingItems []string) {
	modelsToCheck := []interface{}{
		&models.User{},
		&models.Person{},
		&models.Wallet{},
		&models.WalletGroup{},
		&models.Category{},
		&models.Transaction{},
	}

	for _, model := range modelsToCheck {
		if !DB.Migrator().HasTable(model) {
			stmt := &gorm.Statement{DB: DB}
			stmt.Parse(model)
			missingItems = append(missingItems, "table: "+stmt.Schema.Table)
			needed = true
			continue
		}

		// Get all columns that the model expects
		stmt := &gorm.Statement{DB: DB}
		stmt.Parse(model)

		// Get actual database columns
		columnTypes, _ := DB.Migrator().ColumnTypes(model)
		actualColumnsMap := make(map[string]bool)
		for _, col := range columnTypes {
			actualColumnsMap[col.Name()] = true
		}

		for _, field := range stmt.Schema.Fields {
			// Skip fields without database column names
			if field.DBName == "" || field.DBName == "-" {
				continue
			}

			// Check if column exists in database
			if !actualColumnsMap[field.DBName] {
				missingItems = append(missingItems, "column: "+stmt.Schema.Table+"."+field.DBName)
				needed = true
				log.Printf("Migration needed: Missing column %s.%s (field: %s)", stmt.Schema.Table, field.DBName, field.Name)
			}
		}
	}

	return needed, missingItems
}

// MigrateDB runs all migrations
func MigrateDB() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Person{},
		&models.Wallet{},
		&models.WalletGroup{},
		&models.Category{},
		&models.Transaction{},
	)
}

// CloseDB closes the database connection
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
