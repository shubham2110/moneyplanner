package database

import (
	"bufio"
	"fmt"
	"log"
	"moneyplanner/models"
	"os"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection and handles migrations intelligently
func InitDB(dbPath string) error {
	var err error

	// Check if database file exists
	dbExists := fileExists(dbPath)

	// Connect to SQLite database (creates file if doesn't exist)
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if !dbExists {
		// New database - run migrations automatically
		log.Println("Creating new database...")
		err = MigrateDB()
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
		log.Println("✓ Database created with all tables")
		return nil
	}

	// Existing database - check if migrations are needed
	if isMigrationNeeded() {
		log.Println("⚠️  Pending migrations detected")
		if promptForMigration() {
			log.Println("Running migrations...")
			err = MigrateDB()
			if err != nil {
				return fmt.Errorf("migration failed: %w", err)
			}
			log.Println("✓ Migrations completed successfully")
		} else {
			log.Println("Migrations skipped by user")
		}
	} else {
		log.Println("✓ Database schema is up to date")
	}

	return nil
}

// fileExists checks if a file exists
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// isMigrationNeeded checks if there are any pending migrations by comparing ORM models with database schema
func isMigrationNeeded() bool {
	models := []interface{}{
		&models.User{},
		&models.Person{},
		&models.Wallet{},
		&models.WalletGroup{},
		&models.Category{},
		&models.Transaction{},
	}

	for _, model := range models {
		// Check if table exists
		if !DB.Migrator().HasTable(model) {
			return true
		}

		// Check if all columns defined in the model exist in the database
		if !DB.Migrator().HasColumn(model, "id") && !allColumnsExist(model) {
			return true
		}
	}

	return false
}

// allColumnsExist checks if all columns from the model exist in the database
func allColumnsExist(model interface{}) bool {
	// Get all fields from the model
	stmt := &gorm.Statement{DB: DB}
	stmt.Parse(model)

	for _, field := range stmt.Schema.Fields {
		// Skip primary key and associations for simplicity
		if field.PrimaryKey || field.FieldType.Kind() == 4 { // 4 is pointer kind
			continue
		}

		if !DB.Migrator().HasColumn(model, field.DBName) {
			log.Printf("Missing column: %s.%s", stmt.Schema.Table, field.DBName)
			return false
		}
	}

	return true
}

// promptForMigration asks the user if they want to run migrations
func promptForMigration() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to run pending migrations? (yes/no): ")
	
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading input, migrations skipped")
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "yes" || response == "y"
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
