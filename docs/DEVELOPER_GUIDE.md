# MoneyPlanner - Developer Guide

Guide for adding new API endpoints and extending the application.

---

## Table of Contents

1. [Adding a New Endpoint](#adding-a-new-endpoint)
2. [Database Model Creation](#database-model-creation)
3. [API Handler Implementation](#api-handler-implementation)
4. [Testing New Endpoints](#testing-new-endpoints)
5. [Code Structure Conventions](#code-structure-conventions)
6. [Common Patterns](#common-patterns)
7. [Debugging](#debugging)

---

## Adding a New Endpoint

### Example: Create User Listing Endpoint

This guide shows how to add `GET /api/users` endpoint.

#### Step 1: Create API Handler

**File:** `api/users/handler.go` (new file)

```go
package users

import (
	"encoding/json"
	"log"
	"net/http"
	"moneyplanner/database"
	"moneyplanner/models"
)

// GetUsersResponse represents the response for listing users
type GetUsersResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    []models.User   `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// GetUsers retrieves all users from database
func GetUsers() *GetUsersResponse {
	db := database.GetDB()
	if db == nil {
		return &GetUsersResponse{
			Success: false,
			Message: "Database connection failed",
			Error:   "Database not initialized",
		}
	}

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v\n", err)
		return &GetUsersResponse{
			Success: false,
			Message: "Failed to fetch users",
			Error:   err.Error(),
		}
	}

	return &GetUsersResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	}
}
```

#### Step 2: Add Route Handler

**File:** `api/router.go` (modify existing file)

```go
package api

import (
	"encoding/json"
	"log"
	"net/http"

	initAPI "moneyplanner/api/init"
	usersAPI "moneyplanner/api/users"  // Add this import
)

// RegisterRoutes registers all API routes
func RegisterRoutes(mux *http.ServeMux) {
	log.Println("Registering API routes...")

	// Existing routes
	mux.HandleFunc("/api/init", handleInit)
	
	// Add new route
	mux.HandleFunc("/api/users", handleGetUsers)

	log.Println("✓ API routes registered")
}

// Existing handler...
func handleInit(w http.ResponseWriter, r *http.Request) {
	// ... existing code ...
}

// New handler
func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call the service
	resp := usersAPI.GetUsers()

	// Set appropriate status code
	statusCode := http.StatusOK
	if !resp.Success {
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
```

#### Step 3: Update Main Database File

**File:** `database/db.go` (modify existing file)

Add a global database variable:

```go
package database

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	// ... other imports ...
)

var globalDB *gorm.DB  // Add this

func InitDB(dbPath string) error {
	// ... existing code ...
	globalDB = db  // Store the database connection
	return nil
}

// GetDB returns the global database instance
func GetDB() *gorm.DB {
	return globalDB
}
```

#### Step 4: Test the Endpoint

**CURL:**
```bash
curl -X GET http://localhost:8080/api/users \
  -H "Content-Type: application/json"
```

**PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/users" `
  -Method Get | ConvertTo-Json -Depth 10
```

---

## Database Model Creation

### Adding a New Model

**Example: Create a Budget Model**

**File:** `models/budget.go` (new file)

```go
package models

import (
	"time"
	"gorm.io/gorm"
)

// Budget represents a spending budget for a category
type Budget struct {
	BudgetID    uint       `gorm:"primaryKey" json:"budget_id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CategoryID  uint       `gorm:"index" json:"category_id"`
	Category    *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Amount      float64    `json:"amount"`
	Period      string     `json:"period"` // monthly, yearly, weekly
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name
func (Budget) TableName() string {
	return "budgets"
}
```

### Register Model in Database Initialization

**File:** `database/db.go` (modify existing file)

```go
import (
	// ... other imports ...
	"moneyplanner/models"
)

func MigrateDB() error {
	// Add Budget model
	return globalDB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.WalletGroup{},
		&models.Category{},
		&models.Transaction{},
		&models.Person{},
		&models.Budget{},  // Add this
	)
}
```

---

## API Handler Implementation

### Request/Response Pattern

**Standard Request Structure:**
```go
type CreateTransactionRequest struct {
	UserID      uint    `json:"user_id"`
	WalletID    uint    `json:"wallet_id"`
	CategoryID  uint    `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type CreateTransactionResponse struct {
	Success     bool               `json:"success"`
	Message     string             `json:"message"`
	Data        *models.Transaction `json:"data,omitempty"`
	Error       string             `json:"error,omitempty"`
}
```

### Handler with Validation

**Example: Create Transaction Handler**

```go
func CreateTransaction(req *CreateTransactionRequest, userID uint) *CreateTransactionResponse {
	// Validation
	if req.Amount <= 0 {
		return &CreateTransactionResponse{
			Success: false,
			Message: "Invalid amount",
			Error:   "Amount must be greater than 0",
		}
	}

	if req.CategoryID == 0 || req.WalletID == 0 {
		return &CreateTransactionResponse{
			Success: false,
			Message: "Invalid request",
			Error:   "CategoryID and WalletID are required",
		}
	}

	db := database.GetDB()
	if db == nil {
		return &CreateTransactionResponse{
			Success: false,
			Message: "Database error",
			Error:   "Database not initialized",
		}
	}

	// Verify wallet exists and user has access
	var wallet models.Wallet
	if err := db.First(&wallet, req.WalletID).Error; err != nil {
		return &CreateTransactionResponse{
			Success: false,
			Message: "Wallet not found",
			Error:   err.Error(),
		}
	}

	// Create transaction
	transaction := models.Transaction{
		UserID:      userID,
		WalletID:    req.WalletID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		TransactionDate: time.Now(),
	}

	if err := db.Create(&transaction).Error; err != nil {
		log.Printf("Error creating transaction: %v\n", err)
		return &CreateTransactionResponse{
			Success: false,
			Message: "Failed to create transaction",
			Error:   err.Error(),
		}
	}

	return &CreateTransactionResponse{
		Success: true,
		Message: "Transaction created successfully",
		Data:    &transaction,
	}
}
```

---

## Testing New Endpoints

### Unit Testing

**File:** `api/users/handler_test.go`

```go
package users

import (
	"testing"
	"moneyplanner/models"
	"moneyplanner/database"
)

func TestGetUsers(t *testing.T) {
	// Initialize test database
	database.InitDB(":memory:")

	// Create test user
	db := database.GetDB()
	testUser := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
		Type:     "human",
	}
	db.Create(&testUser)

	// Test GetUsers
	resp := GetUsers()

	if !resp.Success {
		t.Errorf("Expected success, got error: %s", resp.Error)
	}

	if len(resp.Data) == 0 {
		t.Error("Expected users to be returned")
	}
}
```

### Integration Testing

**File:** `test-users.sh`

```bash
#!/bin/bash

API="http://localhost:8080/api/users"

echo "Testing GET /api/users"
curl -s -X GET "$API" | jq '.data | length'

echo "Expected: at least 1 user (admin)"
```

---

## Code Structure Conventions

### File Organization

```
api/
├── router.go           # Main route registration
├── init/
│   └── handler.go     # /api/init endpoint
├── users/
│   ├── handler.go     # User handlers
│   └── handler_test.go
├── transactions/
│   ├── handler.go     # Transaction handlers
│   └── handler_test.go
└── reports/
    └── handler.go     # Report handlers
```

### Naming Conventions

- **Packages:** `lowercase` (e.g., `users`, `transactions`)
- **Functions:** `PascalCase` for exported, `camelCase` for private
- **Files:** `lowercase_with_underscores` (e.g., `handler.go`, `handler_test.go`)
- **Constants:** `UPPER_CASE`
- **Interfaces:** `Reader`, `Writer` (suffix with -er)

### Error Handling

```go
// ✓ Good
if err := operation(); err != nil {
	log.Printf("Error description: %v\n", err)
	return &Response{
		Success: false,
		Error: err.Error(),
	}
}

// ✗ Avoid
if err != nil {
	panic(err)  // Never panic in handlers
}
```

---

## Common Patterns

### Middleware Pattern

**File:** `api/middleware.go`

```go
package api

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs all requests
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		next(w, r)

		duration := time.Since(start)
		log.Printf("Completed in %v\n", duration)
	}
}

// Usage in router.go:
// mux.HandleFunc("/api/users", LoggingMiddleware(handleGetUsers))
```

### Query Parameter Handling

```go
func handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	categoryID := r.URL.Query().Get("category_id")

	// Provide defaults
	if limit == "" {
		limit = "10"
	}

	// Call service with filters
	resp := transactionAPI.GetTransactions(limit, offset, categoryID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
```

### Pagination Pattern

```go
type GetTransactionsResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    []models.Transaction   `json:"data,omitempty"`
	Pagination *PaginationInfo     `json:"pagination,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

type PaginationInfo struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"has_more"`
}

func GetTransactions(limit, offset int) *GetTransactionsResponse {
	db := database.GetDB()

	var transactions []models.Transaction
	var total int64

	// Get total count
	db.Model(&models.Transaction{}).Count(&total)

	// Get paginated data
	db.Limit(limit).Offset(offset).Find(&transactions)

	hasMore := (int64(offset) + int64(limit)) < total

	return &GetTransactionsResponse{
		Success: true,
		Data:    transactions,
		Pagination: &PaginationInfo{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}
}
```

---

## Debugging

### Enable Detailed Logging

**In main.go:**
```go
import "gorm.io/logger"

// In InitDB:
db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
	Logger: logger.Default.LogMode(logger.Info),
})
```

### Log Database Queries

```go
// Print generated SQL
db.Session(&gorm.Session{Logger: logger.Default}).
	Where("amount > ?", 100).
	Find(&transactions)
```

### Common Issues

**Issue: Foreign key constraint error**
```go
// ✗ Wrong - violates FK constraint
transaction := models.Transaction{
	CategoryID: 999, // Non-existent category
}
db.Create(&transaction) // Error!

// ✓ Correct - verify FK exists first
var category models.Category
if err := db.First(&category, 999).Error; err != nil {
	return fmt.Errorf("category not found")
}

transaction := models.Transaction{
	CategoryID: 999,
}
db.Create(&transaction)
```

**Issue: Gorm query returns nil slices**
```go
// ✗ Wrong - might be nil
var users []models.User
db.Find(&users)
if users != nil { // Users is never nil!
	// This is always true
}

// ✓ Correct - check with len
if len(users) > 0 {
	// Has results
}
```

---

## Checklist for New Endpoint

Before considering a new endpoint complete:

- [ ] Handler function created in appropriate package
- [ ] Route registered in `api/router.go`
- [ ] Request/response structs defined with JSON tags
- [ ] Input validation implemented
- [ ] Error handling with proper HTTP status codes
- [ ] Logging added for debugging
- [ ] Unit tests written
- [ ] Integration tests created
- [ ] Documentation updated
- [ ] Example commands provided
- [ ] Edge cases handled
- [ ] Database models updated if needed

---

## Performance Optimization

### Index Database Columns

```go
type Transaction struct {
	TransactionID uint   `gorm:"primaryKey" json:"transaction_id"`
	UserID        uint   `gorm:"index" json:"user_id"`           // Add index
	WalletID      uint   `gorm:"index" json:"wallet_id"`         // Add index
	CategoryID    uint   `gorm:"index" json:"category_id"`       // Add index
	TransactionDate time.Time `gorm:"index" json:"transaction_date"` // Add index
	// ...
}
```

### Batch Operations

```go
// ✓ Good - single operation
var users []models.User
db.Where("created_at > ?", lastWeek).Find(&users)

// ✗ Avoid - N+1 problem
var users []models.User
db.Find(&users)
for _, user := range users {
	var wallets []models.Wallet
	db.Where("user_id = ?", user.UserID).Find(&wallets) // Individual query!
}

// ✓ Better - eager load
var users []models.User
db.Preload("Wallets").Find(&users)
```

---

## Version Control Best Practices

### Commit Messages

```
api: add GET /api/users endpoint

- Create users/handler.go with GetUsers function
- Register route in api/router.go
- Add database.GetDB() helper function
- Include comprehensive unit tests

Closes #123
```

### Before Committing

```bash
go fmt ./...        # Format code
go vet ./...        # Check for issues
go test ./...       # Run tests
```

---

Last Updated: December 27, 2025

**Next Steps:**
1. Review existing endpoint in `api/init/handler.go` for patterns
2. Follow the step-by-step guide above for new endpoints
3. Test thoroughly before committing
4. Update documentation for visibility
