# MoneyPlanner API Documentation

**Base URL:** `http://localhost:8080`

**API Base Path:** `/api`

**Current Version:** 1.0

---

## Table of Contents

1. [Authentication](#authentication)
2. [Data Types & Models](#data-types--models)
3. [API Endpoints](#api-endpoints)
4. [Sample Commands](#sample-commands)
5. [Error Handling](#error-handling)
6. [Response Codes](#response-codes)

---

## Authentication

Currently, the API does not require authentication. All endpoints are publicly accessible.

> **Note:** Authentication will be added in future versions.

---

## Data Types & Models

### User
Represents a user in the MoneyPlanner system.

```json
{
  "user_id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "name": "Administrator",
  "type": "human",
  "default_wallet_id": 1
}
```

**Fields:**
- `user_id` (integer): Unique identifier
- `username` (string): Unique username
- `email` (string): User email address
- `name` (string): Full name
- `type` (string): User type - `human` or `bot`
- `default_wallet_id` (integer, nullable): Default wallet for transactions

---

### Wallet
Represents a financial wallet or account.

```json
{
  "wallet_id": 1,
  "wallet_group_id": 1,
  "name": "My Checking Account",
  "balance": 5000.50,
  "icon": "💳",
  "is_enabled": true
}
```

**Fields:**
- `wallet_id` (integer): Unique identifier
- `wallet_group_id` (integer): Parent wallet group
- `name` (string): Wallet name
- `balance` (decimal): Current balance
- `icon` (string): Emoji or icon representation
- `is_enabled` (boolean): Whether wallet is active

---

### WalletGroup
Groups multiple wallets together for organization.

```json
{
  "wallet_group_id": 1,
  "wallet_group_name": "Default"
}
```

**Fields:**
- `wallet_group_id` (integer): Unique identifier
- `wallet_group_name` (string): Group name

---

### Category
Transaction category hierarchy for organizing expenses and income.

```json
{
  "category_id": 1,
  "wallet_id": 1,
  "name": "Income",
  "type": "income",
  "parent_id": null,
  "root_id": 1,
  "is_enabled": true
}
```

**Fields:**
- `category_id` (integer): Unique identifier
- `wallet_id` (integer): Associated wallet
- `name` (string): Category name
- `type` (string): `income` or `expense`
- `parent_id` (integer, nullable): Parent category (for subcategories)
- `root_id` (integer, nullable): Root category ID
- `is_enabled` (boolean): Whether category is active

---

### Transaction
Records a financial transaction.

```json
{
  "transaction_id": 1,
  "user_id": 1,
  "wallet_id": 1,
  "category_id": 1,
  "person_id": null,
  "amount": 500.00,
  "description": "Monthly salary",
  "transaction_date": "2025-12-27T00:00:00Z",
  "created_at": "2025-12-27T10:30:00Z"
}
```

**Fields:**
- `transaction_id` (integer): Unique identifier
- `user_id` (integer): User who created transaction
- `wallet_id` (integer): Source/destination wallet
- `category_id` (integer): Transaction category
- `person_id` (integer, nullable): Related person (for transfers)
- `amount` (decimal): Transaction amount
- `description` (string): Transaction details
- `transaction_date` (datetime): When transaction occurred
- `created_at` (datetime): When record was created

---

### Person
Represents a person (beneficiary, vendor, etc.) for transactions.

```json
{
  "person_id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1-555-0123",
  "address": "123 Main St, City, State 12345"
}
```

**Fields:**
- `person_id` (integer): Unique identifier
- `name` (string): Person's name
- `email` (string, nullable): Email address
- `phone` (string, nullable): Phone number
- `address` (string, nullable): Physical address

---

## API Endpoints

### 1. Initialize Database

**Endpoint:** `POST /api/init`

**Purpose:** Initialize the MoneyPlanner database with schema and default data.

**Required:** First call to set up the system

**Request Body:**

```json
{
  "force_migrate": false,
  "default_wallet_name": "My Wallet",
  "default_wallet_group": "Default",
  "admin_username": "admin",
  "admin_password": "password123",
  "admin_email": "admin@example.com",
  "admin_name": "Administrator"
}
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `force_migrate` | boolean | Yes | Force migration if schema mismatch exists |
| `default_wallet_name` | string | No | Name of primary wallet (default: "My Wallet") |
| `default_wallet_group` | string | No | Name of wallet group (default: "Default") |
| `admin_username` | string | No | Admin account username (default: "admin") |
| `admin_password` | string | No | Admin account password (default: "password123") |
| `admin_email` | string | No | Admin email (default: "admin@example.com") |
| `admin_name` | string | No | Admin full name (default: "Administrator") |

**Response (Success - 200):**

```json
{
  "success": true,
  "message": "Database initialized successfully",
  "data": {
    "database_is_new": true,
    "migration_required": false,
    "admin_user": {
      "user_id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "name": "Administrator",
      "type": "human",
      "default_wallet_id": 1
    },
    "default_wallet": {
      "wallet_id": 1,
      "wallet_group_id": 1,
      "name": "My Wallet",
      "balance": 0,
      "icon": "💰",
      "is_enabled": true
    },
    "default_wallet_group": {
      "wallet_group_id": 1,
      "wallet_group_name": "Default"
    }
  },
  "missing_items": []
}
```

**Response (Migration Required - 400):**

```json
{
  "success": false,
  "message": "Database schema does not match ORM models",
  "error": "Pending migrations detected. Set force_migrate: true to proceed",
  "data": {
    "database_is_new": false,
    "migration_required": true
  },
  "missing_items": [
    "table: wallet_groups",
    "column: users.created_at",
    "column: transactions.updated_at"
  ]
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Operation success status |
| `message` | string | Human-readable message |
| `data` | object | Initialized entities (null on failure) |
| `error` | string | Error details if failed |
| `missing_items` | array | List of missing tables/columns requiring migration |

**Status Codes:**
- `200 OK`: Database initialized successfully
- `400 Bad Request`: Migration needed or invalid input
- `500 Internal Server Error`: Unexpected server error

---

## Sample Commands

All examples use `curl` command-line tool. You can also use Postman, Insomnia, or any HTTP client.

### 1. Initialize New Database

**Command:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{
    "force_migrate": true,
    "default_wallet_name": "My Checking Account",
    "default_wallet_group": "Personal",
    "admin_username": "shivani",
    "admin_password": "securepass123",
    "admin_email": "shivani@example.com",
    "admin_name": "Shivani Sharma"
  }'
```

**Expected Output:**
```json
{
  "success": true,
  "message": "Database initialized successfully",
  "data": {
    "database_is_new": true,
    "migration_required": false,
    "admin_user": {
      "user_id": 1,
      "username": "shivani",
      "email": "shivani@example.com",
      "name": "Shivani Sharma",
      "type": "human",
      "default_wallet_id": 1
    },
    "default_wallet": {
      "wallet_id": 1,
      "wallet_group_id": 1,
      "name": "My Checking Account",
      "balance": 0,
      "icon": "💰",
      "is_enabled": true
    },
    "default_wallet_group": {
      "wallet_group_id": 1,
      "wallet_group_name": "Personal"
    }
  },
  "missing_items": []
}
```

---

### 2. Initialize with Minimal Parameters

**Command:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}'
```

**Description:** Uses all default values for wallet, wallet group, and admin credentials.

---

### 3. Check Database Status Without Initialization

**Command:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}'
```

**Expected Output (if DB already initialized):**
```json
{
  "success": true,
  "message": "Database is already initialized",
  "data": {
    "database_is_new": false,
    "migration_required": false
  },
  "missing_items": []
}
```

---

### 4. Handle Migration with Details

**Command:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}'
```

**Expected Output (if migration needed):**
```json
{
  "success": false,
  "message": "Database schema does not match ORM models",
  "error": "Pending migrations detected. Set force_migrate: true to proceed",
  "data": {
    "database_is_new": false,
    "migration_required": true
  },
  "missing_items": [
    "table: wallet_groups",
    "table: wallets",
    "column: transactions.updated_at"
  ]
}
```

**Follow-up:** Retry with `"force_migrate": true` to apply migrations.

---

## Using PowerShell (Windows)

### 1. Basic Initialization

```powershell
$body = @{
    force_migrate = $true
    default_wallet_name = "My Wallet"
    admin_username = "admin"
    admin_password = "pass123"
    admin_email = "admin@example.com"
    admin_name = "Administrator"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/init" `
  -Method Post `
  -Headers @{"Content-Type"="application/json"} `
  -Body $body
```

### 2. Check Status Only

```powershell
$body = @{force_migrate = $false} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/init" `
  -Method Post `
  -Headers @{"Content-Type"="application/json"} `
  -Body $body

if ($response.success) {
    Write-Host "Database initialized successfully" -ForegroundColor Green
} else {
    Write-Host "Error: $($response.error)" -ForegroundColor Red
    Write-Host "Missing items: $($response.missing_items -join ', ')"
}
```

---

## Using JavaScript/Node.js

### 1. Basic Initialization

```javascript
const initData = {
  force_migrate: true,
  default_wallet_name: "My Wallet",
  admin_username: "admin",
  admin_password: "password123",
  admin_email: "admin@example.com",
  admin_name: "Administrator"
};

fetch('http://localhost:8080/api/init', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(initData)
})
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

### 2. Check Status with Error Handling

```javascript
async function checkDatabaseStatus() {
  try {
    const response = await fetch('http://localhost:8080/api/init', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ force_migrate: false })
    });
    
    const data = await response.json();
    
    if (data.success) {
      console.log('✓ Database is ready');
    } else {
      console.warn('✗ Migration needed:', data.missing_items);
    }
  } catch (error) {
    console.error('Connection failed:', error);
  }
}

checkDatabaseStatus();
```

---

## Using Python

### 1. Basic Initialization

```python
import requests
import json

url = 'http://localhost:8080/api/init'
payload = {
    'force_migrate': True,
    'default_wallet_name': 'My Wallet',
    'admin_username': 'admin',
    'admin_password': 'password123',
    'admin_email': 'admin@example.com',
    'admin_name': 'Administrator'
}

response = requests.post(url, json=payload)
print(json.dumps(response.json(), indent=2))
```

### 2. Status Check with Error Handling

```python
import requests

def check_database_status():
    url = 'http://localhost:8080/api/init'
    payload = {'force_migrate': False}
    
    try:
        response = requests.post(url, json=payload)
        data = response.json()
        
        if data['success']:
            print('✓ Database is ready')
        else:
            print('✗ Migration needed')
            for item in data.get('missing_items', []):
                print(f'  - {item}')
    except requests.exceptions.ConnectionError:
        print('✗ Cannot connect to server')

check_database_status()
```

---

## Error Handling

### Common Error Scenarios

#### 1. Invalid JSON

**Request:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{invalid json}'
```

**Response (400):**
```json
{
  "error": "Invalid request body: unexpected character 'i' looking for beginning of object key string"
}
```

---

#### 2. Connection Refused

**Cause:** Server not running

**Solution:** Start server with:
```bash
cd c:\Users\Shivani\Documents\Shubham\moneyplanner
go run main.go
```

---

#### 3. Database File Already Exists

**Request:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}'
```

**Response (200):**
```json
{
  "success": true,
  "message": "Database is already initialized",
  "data": {
    "database_is_new": false,
    "migration_required": false
  },
  "missing_items": []
}
```

---

## Response Codes

| Code | Status | Meaning |
|------|--------|---------|
| 200 | OK | Request successful |
| 400 | Bad Request | Invalid request or migration needed |
| 405 | Method Not Allowed | Wrong HTTP method (only POST allowed) |
| 500 | Internal Server Error | Server error, check logs |

---

## Default Data Created on Initialization

When `force_migrate: true` and database is new:

### Created User Account
- **Username:** admin
- **Email:** admin@example.com
- **Password:** password123
- **Type:** human
- **Role:** System administrator

### Created Wallet Structure
- **Wallet Group:** "Default"
- **Wallet:** "My Wallet"
- **Initial Balance:** $0.00

### Created Categories
#### Income (Root)
- Salary
- Refund
- Gift
- Interest
- Other Income

#### Expense (Root)
- Groceries
- Transportation
- Utilities
- Entertainment
- Medical
- Shopping
- Food & Dining
- Other Expense

### Dummy Transactions
- One transaction per category with amount = 0.00
- Prevents foreign key constraint violations
- Can be safely deleted

---

## Future Endpoints (Planned)

The following endpoints are planned for future versions:

- `GET /api/users` - List all users
- `POST /api/users` - Create new user
- `GET /api/wallets` - List wallets
- `POST /api/transactions` - Record transaction
- `GET /api/transactions` - List transactions with filters
- `GET /api/categories` - List categories
- `GET /api/reports/summary` - Financial summary
- `GET /api/reports/spending` - Spending analysis

---

## Development Notes

### Database Location
- **Path:** `moneyplanner.db` (SQLite)
- **Location:** Project root directory
- **Size:** Typically < 1MB for normal usage

### API Server
- **Default Port:** 8080
- **Framework:** Go standard `net/http`
- **ORM:** GORM with SQLite driver

### Logging
Server logs all API calls and database operations to console.

---

## Support & Contact

For issues, questions, or feature requests, please refer to the project documentation or contact the development team.

**Last Updated:** December 27, 2025
**API Version:** 1.0
