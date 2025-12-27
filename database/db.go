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

	// File exists, check if any tables exist by looking for a key table
	// If wallet_groups table doesn't exist, we consider it a new database
	return !DB.Migrator().HasTable("wallet_groups")
}

// IsMigrationNeeded checks if there are pending migrations
func IsMigrationNeeded() bool {
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
			return true
		}

		// Check if all columns exist
		stmt := &gorm.Statement{DB: DB}
		stmt.Parse(model)

		for _, field := range stmt.Schema.Fields {
			if field.FieldType.Kind() == 4 { // Skip pointers
				continue
			}

			if !DB.Migrator().HasColumn(model, field.DBName) {
				log.Printf("Missing column: %s.%s", stmt.Schema.Table, field.DBName)
				return true
			}
		}
	}

	return false
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
