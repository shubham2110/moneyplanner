# MoneyPlanner API - Complete Examples

This document contains ready-to-use examples for all API endpoints.

---

## Table of Contents

1. [CURL Examples](#curl-examples)
2. [PowerShell Examples](#powershell-examples)
3. [JavaScript/Node.js Examples](#javascriptnodejs-examples)
4. [Python Examples](#python-examples)
5. [Postman Collection](#postman-collection)
6. [Testing Workflows](#testing-workflows)

---

## CURL Examples

All CURL examples use Linux/MacOS/Git Bash syntax. For Windows CMD, use single quotes differently.

### Basic Initialization (Minimal)

**Description:** Initialize database with all default values

**Command:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}'
```

**Output:**
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

---

### Full Initialization with Custom Values

**Description:** Complete setup with custom credentials and wallet names

**Command:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{
    "force_migrate": true,
    "default_wallet_name": "Checking Account",
    "default_wallet_group": "Personal Finance",
    "admin_username": "shivani",
    "admin_password": "MySecurePass123!",
    "admin_email": "shivani.sharma@email.com",
    "admin_name": "Shivani Sharma"
  }'
```

---

### Status Check (Non-destructive)

**Description:** Check if database is initialized and migrations are needed

**Command:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}'
```

**Expected Output (Already Initialized):**
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

### Handling Migration Required

**Description:** When migration is needed but force_migrate is false

**Command:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}'
```

**Expected Output (Migration Needed):**
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
    "table: transactions",
    "column: users.created_at"
  ]
}
```

**Fix:** Retry with `force_migrate: true`
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}'
```

---

### Save Response to File

**Description:** Save API response to a JSON file for analysis

**Command:**
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}' \
  | jq '.' > init_response.json

cat init_response.json
```

---

### Pretty Print Response

**Description:** Format response with proper indentation (requires `jq`)

**Command:**
```bash
curl -s -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}' | jq '.data'
```

---

### Extract Specific Fields

**Description:** Extract just the admin username from response

**Command:**
```bash
curl -s -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}' | jq '.data.admin_user.username'
```

---

## PowerShell Examples

All examples use Windows PowerShell 5.1+

### Basic Initialization

**Description:** Initialize with default settings

**Script:**
```powershell
$body = @{force_migrate = $true} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/init" `
  -Method Post `
  -Headers @{"Content-Type"="application/json"} `
  -Body $body

$response | ConvertTo-Json -Depth 10
```

---

### Full Initialization with Custom Values

**Description:** Custom setup with complete credentials

**Script:**
```powershell
$initRequest = @{
    force_migrate = $true
    default_wallet_name = "My Checking"
    default_wallet_group = "Personal"
    admin_username = "shivani"
    admin_password = "SecurePassword123!"
    admin_email = "shivani@example.com"
    admin_name = "Shivani Sharma"
}

$body = $initRequest | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/init" `
  -Method Post `
  -Headers @{"Content-Type"="application/json"} `
  -Body $body

Write-Host "Initialization Status: $($response.success)" -ForegroundColor Green
Write-Host "Message: $($response.message)"
Write-Host ""
Write-Host "Admin User: $($response.data.admin_user.username)"
Write-Host "Wallet: $($response.data.default_wallet.name)"
```

**Output:**
```
Initialization Status: True
Message: Database initialized successfully

Admin User: shivani
Wallet: My Checking
```

---

### Status Check with Error Handling

**Description:** Check database status with proper error handling

**Script:**
```powershell
try {
    $body = @{force_migrate = $false} | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/init" `
      -Method Post `
      -Headers @{"Content-Type"="application/json"} `
      -Body $body
    
    if ($response.success) {
        Write-Host "✓ Database is ready" -ForegroundColor Green
        Write-Host "New DB: $($response.data.database_is_new)"
        Write-Host "Migration Required: $($response.data.migration_required)"
    } else {
        Write-Host "✗ Error: $($response.message)" -ForegroundColor Red
        Write-Host "Missing items:"
        $response.missing_items | ForEach-Object {
            Write-Host "  - $_"
        }
    }
}
catch {
    Write-Host "✗ Connection failed: $_" -ForegroundColor Red
}
```

---

### Reusable Function

**Description:** Create a reusable PowerShell function for initialization

**Script:**
```powershell
function Initialize-MoneyPlanner {
    param(
        [string]$WalletName = "My Wallet",
        [string]$AdminUsername = "admin",
        [string]$AdminPassword = "password123",
        [string]$AdminEmail = "admin@example.com",
        [string]$AdminName = "Administrator",
        [switch]$ForceMigrate
    )

    $initRequest = @{
        force_migrate = $ForceMigrate
        default_wallet_name = $WalletName
        admin_username = $AdminUsername
        admin_password = $AdminPassword
        admin_email = $AdminEmail
        admin_name = $AdminName
    }

    $body = $initRequest | ConvertTo-Json
    
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/api/init" `
          -Method Post `
          -Headers @{"Content-Type"="application/json"} `
          -Body $body

        return $response
    }
    catch {
        Write-Error "Failed to initialize database: $_"
    }
}

# Usage:
# Initialize-MoneyPlanner -WalletName "Checking" -AdminUsername "john" -ForceMigrate
```

---

### Save to CSV/Excel

**Description:** Export response data in CSV format

**Script:**
```powershell
$body = @{force_migrate = $false} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/init" `
  -Method Post `
  -Headers @{"Content-Type"="application/json"} `
  -Body $body

if ($response.data.admin_user) {
    [PSCustomObject]@{
        Username = $response.data.admin_user.username
        Email = $response.data.admin_user.email
        Name = $response.data.admin_user.name
        Type = $response.data.admin_user.type
        DefaultWalletId = $response.data.admin_user.default_wallet_id
    } | Export-Csv -Path "admin_user.csv" -NoTypeInformation
    
    Write-Host "Exported to admin_user.csv"
}
```

---

## JavaScript/Node.js Examples

### Browser Fetch (ES6)

**Description:** Simple fetch in browser console

**Code:**
```javascript
const data = {
  force_migrate: true,
  default_wallet_name: "My Wallet",
  admin_username: "admin"
};

fetch('http://localhost:8080/api/init', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(data)
})
  .then(res => res.json())
  .then(data => console.log(JSON.stringify(data, null, 2)))
  .catch(err => console.error('Error:', err));
```

---

### Node.js with Error Handling

**Description:** Complete Node.js implementation with error handling

**Code:**
```javascript
const http = require('http');

function initializeDatabase(options = {}) {
  const defaultOptions = {
    force_migrate: true,
    default_wallet_name: 'My Wallet',
    admin_username: 'admin',
    admin_password: 'password123',
    admin_email: 'admin@example.com',
    admin_name: 'Administrator'
  };

  const data = JSON.stringify({ ...defaultOptions, ...options });

  const requestOptions = {
    hostname: 'localhost',
    port: 8080,
    path: '/api/init',
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': data.length
    }
  };

  return new Promise((resolve, reject) => {
    const req = http.request(requestOptions, (res) => {
      let body = '';

      res.on('data', chunk => body += chunk);
      res.on('end', () => {
        try {
          const response = JSON.parse(body);
          resolve(response);
        } catch (error) {
          reject(error);
        }
      });
    });

    req.on('error', reject);
    req.write(data);
    req.end();
  });
}

// Usage
initializeDatabase({ 
  admin_username: 'shivani',
  admin_password: 'SecurePass123'
})
  .then(response => {
    if (response.success) {
      console.log('✓ Database initialized');
      console.log('Admin:', response.data.admin_user.username);
    } else {
      console.error('✗ Failed:', response.error);
    }
  })
  .catch(error => console.error('Error:', error));
```

---

### Using Axios (Node.js)

**Description:** Using popular axios library

**Code:**
```javascript
const axios = require('axios');

async function initMoneyPlanner() {
  try {
    const response = await axios.post('http://localhost:8080/api/init', {
      force_migrate: true,
      default_wallet_name: 'My Wallet',
      admin_username: 'admin'
    });

    console.log('Success:', response.data.success);
    console.log('Message:', response.data.message);
    console.log('Admin User:', response.data.data?.admin_user?.username);
  } catch (error) {
    if (error.response) {
      console.error('API Error:', error.response.status, error.response.data);
    } else {
      console.error('Connection Error:', error.message);
    }
  }
}

initMoneyPlanner();
```

---

### Async/Await Pattern

**Description:** Modern async/await syntax

**Code:**
```javascript
async function checkDatabaseStatus() {
  const url = 'http://localhost:8080/api/init';
  
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ force_migrate: false })
    });

    const data = await response.json();

    if (data.success) {
      console.log('✓ Database Status:');
      console.log(`  Is New: ${data.data.database_is_new}`);
      console.log(`  Migration Needed: ${data.data.migration_required}`);
    } else {
      console.error('✗ Issues Found:');
      data.missing_items?.forEach(item => console.log(`  - ${item}`));
    }
  } catch (error) {
    console.error('✗ Connection Failed:', error.message);
  }
}

checkDatabaseStatus();
```

---

## Python Examples

### Basic Request

**Description:** Simple Python script to initialize database

**Code:**
```python
import requests
import json

url = 'http://localhost:8080/api/init'
payload = {'force_migrate': True}

response = requests.post(url, json=payload)
data = response.json()

print(json.dumps(data, indent=2))
```

---

### Full Initialization

**Description:** Complete initialization with custom values

**Code:**
```python
import requests

url = 'http://localhost:8080/api/init'

init_data = {
    'force_migrate': True,
    'default_wallet_name': 'My Checking Account',
    'default_wallet_group': 'Personal',
    'admin_username': 'shivani',
    'admin_password': 'SecurePassword123!',
    'admin_email': 'shivani@example.com',
    'admin_name': 'Shivani Sharma'
}

response = requests.post(url, json=init_data)
result = response.json()

if result['success']:
    print('✓ Database initialized successfully')
    print(f"Admin user: {result['data']['admin_user']['username']}")
    print(f"Wallet: {result['data']['default_wallet']['name']}")
else:
    print(f'✗ Error: {result["error"]}')
```

---

### Status Check with Logging

**Description:** Check database status with detailed logging

**Code:**
```python
import requests
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def check_database_status():
    url = 'http://localhost:8080/api/init'
    
    try:
        response = requests.post(url, json={'force_migrate': False})
        data = response.json()
        
        logger.info(f"Status Code: {response.status_code}")
        
        if data['success']:
            logger.info('✓ Database is ready')
            logger.info(f"  Is New: {data['data']['database_is_new']}")
            logger.info(f"  Migration Required: {data['data']['migration_required']}")
        else:
            logger.error(f"✗ {data['message']}")
            if data.get('missing_items'):
                logger.error("Missing items:")
                for item in data['missing_items']:
                    logger.error(f"  - {item}")
                    
    except requests.exceptions.ConnectionError:
        logger.error('✗ Cannot connect to server')
    except Exception as e:
        logger.error(f'✗ Unexpected error: {e}')

check_database_status()
```

---

### Class-Based Implementation

**Description:** Object-oriented approach for reusability

**Code:**
```python
import requests
from typing import Dict, Optional

class MoneyPlannerClient:
    def __init__(self, base_url: str = 'http://localhost:8080'):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({'Content-Type': 'application/json'})
    
    def initialize_database(
        self,
        force_migrate: bool = False,
        wallet_name: str = 'My Wallet',
        admin_username: str = 'admin',
        admin_password: str = 'password123',
        admin_email: str = 'admin@example.com',
        admin_name: str = 'Administrator'
    ) -> Dict:
        """Initialize the MoneyPlanner database"""
        
        payload = {
            'force_migrate': force_migrate,
            'default_wallet_name': wallet_name,
            'admin_username': admin_username,
            'admin_password': admin_password,
            'admin_email': admin_email,
            'admin_name': admin_name
        }
        
        response = self.session.post(
            f'{self.base_url}/api/init',
            json=payload
        )
        
        return response.json()
    
    def check_status(self) -> Dict:
        """Check database status without making changes"""
        return self.initialize_database(force_migrate=False)

# Usage
client = MoneyPlannerClient()

# Initialize
result = client.initialize_database(
    force_migrate=True,
    wallet_name='Checking Account',
    admin_username='shivani'
)

print(f"Success: {result['success']}")
print(f"Message: {result['message']}")
```

---

## Postman Collection

### JSON Collection Export

**Description:** Import this into Postman for easy testing

**Collection File (Copy to file.json):**
```json
{
  "info": {
    "name": "MoneyPlanner API",
    "description": "Collection for MoneyPlanner API endpoints",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Initialize Database",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"force_migrate\": true,\n  \"default_wallet_name\": \"My Wallet\",\n  \"default_wallet_group\": \"Personal\",\n  \"admin_username\": \"admin\",\n  \"admin_password\": \"password123\",\n  \"admin_email\": \"admin@example.com\",\n  \"admin_name\": \"Administrator\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/init",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "init"]
        }
      }
    },
    {
      "name": "Check Database Status",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"force_migrate\": false\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/init",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "init"]
        }
      }
    }
  ]
}
```

### Import Steps:
1. Open Postman
2. Click "Import"
3. Select "Paste Raw Text"
4. Paste the JSON above
5. Click "Import"

---

## Testing Workflows

### Workflow 1: Fresh Setup

**Goal:** Set up database from scratch

**Steps:**
```bash
# 1. Delete old database
rm -f moneyplanner.db

# 2. Initialize with defaults
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}'

# 3. Verify
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}'
```

---

### Workflow 2: Custom Setup

**Goal:** Initialize with custom credentials

**Steps:**
```bash
# 1. Initialize with custom values
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{
    "force_migrate": true,
    "default_wallet_name": "Business Account",
    "admin_username": "owner",
    "admin_password": "BusinessPass123!"
  }'

# 2. Extract and save credentials
ADMIN_USER=$(curl -s -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}' | jq '.data.admin_user')

echo "Admin credentials saved"
```

---

### Workflow 3: Error Recovery

**Goal:** Handle and fix migration errors

**Steps:**
```bash
# 1. Check status (may fail with missing items)
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}' | jq '.missing_items'

# 2. If migration needed, apply it
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}'

# 3. Verify success
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}' | jq '.success'
```

---

## Testing with curl (Batch)

**Description:** Run multiple tests in sequence

**Script (test_api.sh):**
```bash
#!/bin/bash

API="http://localhost:8080/api/init"

echo "=== Test 1: Initialize Database ==="
curl -s -X POST "$API" \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}' | jq '.success'

echo ""
echo "=== Test 2: Check Status ==="
curl -s -X POST "$API" \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}' | jq '.data.migration_required'

echo ""
echo "✓ All tests completed"
```

**Run:**
```bash
chmod +x test_api.sh
./test_api.sh
```

---

## Performance Testing

**Description:** Load test the API

**Using Apache Bench:**
```bash
ab -n 100 -c 10 -p payload.json \
  -T application/json \
  http://localhost:8080/api/init
```

**Using wrk:**
```bash
wrk -t4 -c100 -d30s \
  -s init_test.lua \
  http://localhost:8080/api/init
```

---

Last Updated: December 27, 2025
