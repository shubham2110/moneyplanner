# MoneyPlanner - Backend API

A Go-based financial management application for tracking wallets, categories, transactions, and budgets.

---

## 🚀 Quick Start

### 1. Start the Server

```bash
cd c:\Users\Shivani\Documents\Shubham\moneyplanner
go run main.go
```

**Output:**
```
MoneyPlanner server running on :8080
Initialize database by calling:
  POST http://localhost:8080/api/init
  Body: {"force_migrate": true}
```

### 2. Initialize Database

**CURL:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}'
```

**PowerShell:**
```powershell
$body = @{force_migrate = $true} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/init" `
  -Method Post -Headers @{"Content-Type"="application/json"} -Body $body
```

### 3. Verify Setup

```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}'
```

---

## 📁 Project Structure

```
moneyplanner/
├── main.go                      # Application entry point
├── go.mod / go.sum              # Dependency management
│
├── api/                         # API handlers
│   ├── router.go               # Route registration
│   └── init/
│       └── handler.go          # Database initialization logic
│
├── database/                    # Database layer
│   └── db.go                   # GORM setup, migrations, schema checks
│
├── models/                      # Data models
│   ├── user.go                 # User entity
│   ├── wallet.go               # Wallet entity
│   ├── walletgroup.go          # Wallet group entity
│   ├── category.go             # Category hierarchy
│   ├── transaction.go          # Transaction records
│   └── person.go               # Person entity
│
├── docs/                        # Swagger documentation
│   └── doc.go                  # API specification
│
├── moneyplanner.db             # SQLite database (auto-created)
│
└── Documentation/
    ├── README.md               # This file
    ├── API_DOCUMENTATION.md    # Complete API documentation
    ├── QUICK_REFERENCE.md      # Quick lookup guide
    ├── EXAMPLES.md             # Code examples in multiple languages
    └── TESTING.md              # Testing scripts and procedures
```

---

## 📚 Documentation

### For Quick Lookup
👉 **[QUICK_REFERENCE.md](./QUICK_REFERENCE.md)** - Most common commands at a glance

### For Complete API Reference
👉 **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)** - Full endpoint documentation with models

### For Code Examples
👉 **[EXAMPLES.md](./EXAMPLES.md)** - Examples in curl, PowerShell, JavaScript, Python

### For Testing
👉 **[TESTING.md](./TESTING.md)** - Test scripts and workflows

---

## 🔌 Current API Endpoints

### Initialize Database
```
POST /api/init
```

Initialize the MoneyPlanner database with schema and default data.

**Request:**
```json
{
  "force_migrate": true,
  "default_wallet_name": "My Wallet",
  "admin_username": "admin",
  "admin_password": "password123"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Database initialized successfully",
  "data": {
    "admin_user": { ... },
    "default_wallet": { ... },
    "default_wallet_group": { ... }
  }
}
```

See [API_DOCUMENTATION.md](./API_DOCUMENTATION.md#1-initialize-database) for complete details.

---

## 🗂️ Data Models

### User
System users with credentials and wallet associations.
- `user_id` (int)
- `username` (string)
- `email` (string)
- `name` (string)
- `type` (string: human/bot)
- `default_wallet_id` (int)

### Wallet
Financial accounts or wallets.
- `wallet_id` (int)
- `wallet_group_id` (int)
- `name` (string)
- `balance` (decimal)
- `is_enabled` (boolean)

### Category
Transaction categories with hierarchy.
- `category_id` (int)
- `name` (string)
- `type` (string: income/expense)
- `parent_id` (int, nullable)
- `root_id` (int, nullable)

### Transaction
Financial transaction records.
- `transaction_id` (int)
- `user_id` (int)
- `wallet_id` (int)
- `category_id` (int)
- `amount` (decimal)
- `description` (string)
- `transaction_date` (datetime)

### WalletGroup
Groups for organizing wallets.
- `wallet_group_id` (int)
- `wallet_group_name` (string)

### Person
People involved in transactions.
- `person_id` (int)
- `name` (string)
- `email` (string, nullable)
- `phone` (string, nullable)

---

## 🛠️ Technology Stack

- **Language:** Go 1.24
- **Database:** SQLite (glebarez/sqlite)
- **ORM:** GORM v1.25.5
- **HTTP Framework:** Go standard `net/http`
- **Documentation:** Swagger/OpenAPI (planned)

---

## 📥 Installation & Setup

### Prerequisites
- Go 1.24 or higher
- GORM and dependencies (auto-downloaded via go mod)

### Step 1: Clone/Navigate to Project
```bash
cd c:\Users\Shivani\Documents\Shubham\moneyplanner
```

### Step 2: Download Dependencies
```bash
go mod download
go mod tidy
```

### Step 3: Build (Optional)
```bash
go build -o moneyplanner.exe main.go
```

### Step 4: Run Server
```bash
go run main.go
```

### Step 5: Initialize Database
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}'
```

---

## 🧪 Testing

### Quick Test
```bash
# Run server in one terminal
go run main.go

# In another terminal, run tests
.\test-api.ps1
```

See [TESTING.md](./TESTING.md) for:
- PowerShell test scripts
- Bash test scripts
- Windows batch scripts
- Docker testing
- Performance testing
- Monitoring scripts

---

## 🗄️ Database

### Location
`moneyplanner.db` - SQLite database in project root

### Auto-Initialization
Database is automatically created on first `/api/init` call with:
- **Tables:** users, wallets, wallet_groups, categories, transactions, persons
- **Relationships:** Foreign keys, many-to-many associations
- **Default Data:** Admin user, wallets, categories, sample transactions

### Schema
All tables are created via GORM migrations in `database/db.go`:
```go
database.MigrateDB() // Runs AutoMigrate on all models
```

---

## 📊 Default Data Created

When initializing a new database, the system creates:

### Admin User
- Username: `admin`
- Password: `password123`
- Email: `admin@example.com`
- Type: `human`

### Wallet Structure
- **Group:** "Default"
- **Wallet:** "My Wallet" (Balance: $0.00)

### Categories
**Income Categories:**
- Salary
- Refund
- Gift
- Interest
- Other Income

**Expense Categories:**
- Groceries
- Transportation
- Utilities
- Entertainment
- Medical
- Shopping
- Food & Dining
- Other Expense

### Transactions
- One dummy transaction per category (amount = 0.00)
- Prevents foreign key constraint violations

---

## 🔒 Security Notes

⚠️ **Current Status:** No authentication/authorization
- All endpoints are publicly accessible
- Default admin credentials are hardcoded
- Suitable for development/testing only

**Future Improvements:**
- JWT authentication
- Password hashing (bcrypt)
- Role-based access control
- API key management
- Rate limiting

---

## 📈 Future Endpoints (Planned)

- `GET /api/users` - List users
- `POST /api/users` - Create user
- `GET /api/wallets` - List wallets
- `GET /api/wallets/{id}` - Get wallet details
- `POST /api/wallets` - Create wallet
- `GET /api/transactions` - List transactions
- `POST /api/transactions` - Record transaction
- `GET /api/categories` - List categories
- `GET /api/reports/summary` - Financial summary
- `GET /api/reports/spending` - Spending analysis

---

## 🐛 Troubleshooting

### Issue: "Port 8080 already in use"
```bash
# Find and kill process using port 8080
netstat -ano | findstr :8080      # Windows
lsof -i :8080 | grep LISTEN        # macOS/Linux
```

### Issue: "Cannot connect to database"
- Check `moneyplanner.db` exists in project root
- Verify database file permissions
- Try deleting old database and reinitializing:
  ```bash
  rm moneyplanner.db
  go run main.go
  curl -X POST http://localhost:8080/api/init \
    -H "Content-Type: application/json" \
    -d '{"force_migrate": true}'
  ```

### Issue: "Migration needed" response
Run with `force_migrate: true`:
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}'
```

### Issue: "Invalid JSON" error
Validate JSON at https://jsonlint.com/ or use:
```bash
echo '{"force_migrate": true}' | python -m json.tool
```

---

## 📝 API Examples

### Initialize with Defaults
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}'
```

### Initialize with Custom Credentials
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{
    "force_migrate": true,
    "admin_username": "john",
    "admin_password": "JohnPass123",
    "admin_email": "john@example.com",
    "admin_name": "John Doe"
  }'
```

### Check Status
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}'
```

See [EXAMPLES.md](./EXAMPLES.md) for examples in PowerShell, JavaScript, Python, and more.

---

## 🔗 API Base URL

```
http://localhost:8080
```

**API Prefix:**
```
/api
```

---

## 📞 Support

For detailed information:
- 📖 Full API docs: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
- ⚡ Quick reference: [QUICK_REFERENCE.md](./QUICK_REFERENCE.md)
- 💻 Code examples: [EXAMPLES.md](./EXAMPLES.md)
- 🧪 Testing guide: [TESTING.md](./TESTING.md)

---

## 📄 License

Private project - MoneyPlanner Financial Management System

---

## 🎯 Project Status

✅ **Complete:**
- Database models and relationships
- API initialization endpoint
- Default data creation
- Database schema validation
- GORM ORM integration

🔄 **In Progress:**
- Swagger UI integration
- API documentation refinement

📋 **Planned:**
- User management endpoints
- Transaction CRUD operations
- Category management
- Financial reports and analytics
- Authentication & authorization
- Database backup/export
- Mobile app support

---

## 👨‍💻 Development

### Build from Source
```bash
go build -o moneyplanner.exe main.go
```

### Run Tests
```bash
# PowerShell
.\test-api-basic.ps1

# Bash
./test-api.sh

# Node.js
node test.js
```

### Database Inspection

View database with SQLite:
```bash
sqlite3 moneyplanner.db

# Common queries:
.tables
.schema users
SELECT COUNT(*) FROM transactions;
```

---

## 📅 Last Updated

December 27, 2025

---

## 🙏 Acknowledgments

Built with:
- [GORM](https://gorm.io/) - Go ORM
- [glebarez/sqlite](https://github.com/glebarez/sqlite) - SQLite driver
- Go standard library

---

**Ready to get started? See [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) for the fastest way to initialize your database!**
