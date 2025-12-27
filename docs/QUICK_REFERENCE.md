# MoneyPlanner API - Quick Reference

## Quick Start

### 1. Start the Server
```bash
cd c:\Users\Shivani\Documents\Shubham\moneyplanner
go run main.go
```

Server runs on: **http://localhost:8080**

---

## API Endpoints Summary

| Method | Endpoint | Purpose | Status |
|--------|----------|---------|--------|
| POST | `/api/init` | Initialize database | ✅ Active |

---

## Most Common Commands

### Initialize New Database (Complete Setup)
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{
    "force_migrate": true,
    "default_wallet_name": "My Wallet",
    "default_wallet_group": "Personal",
    "admin_username": "admin",
    "admin_password": "password123",
    "admin_email": "admin@example.com",
    "admin_name": "Administrator"
  }'
```

### Initialize with Defaults (Minimal)
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": true}'
```

### Check Database Status (Non-destructive)
```bash
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"force_migrate": false}'
```

---

## PowerShell Examples

### Full Initialization
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
  -Body $body | ConvertTo-Json -Depth 10
```

### Quick Status Check
```powershell
$body = @{force_migrate = $false} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/init" `
  -Method Post -Headers @{"Content-Type"="application/json"} -Body $body
```

---

## JavaScript Examples

### Browser Fetch (Copy-Paste Ready)
```javascript
const data = {
  force_migrate: true,
  default_wallet_name: "My Wallet",
  admin_username: "admin",
  admin_password: "password123"
};

fetch('http://localhost:8080/api/init', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(data)
})
  .then(r => r.json())
  .then(d => console.log(d));
```

---

## Python Examples

### Quick Init Script
```python
import requests

requests.post(
  'http://localhost:8080/api/init',
  json={'force_migrate': True, 'default_wallet_name': 'My Wallet'}
).json()
```

---

## Default Credentials

**Auto-created on first initialization:**
- **Username:** admin
- **Password:** password123
- **Email:** admin@example.com

---

## Default Data Created

✅ **Wallet Group:** Default  
✅ **Wallet:** My Wallet ($0.00)  
✅ **Categories:** Income & Expense (10+ subcategories)  
✅ **Admin User:** Created with full access  

---

## Response Examples

### Success Response (200)
```json
{
  "success": true,
  "message": "Database initialized successfully",
  "data": { ... },
  "missing_items": []
}
```

### Migration Needed (400)
```json
{
  "success": false,
  "message": "Database schema does not match ORM models",
  "error": "Set force_migrate: true to proceed",
  "missing_items": ["table: wallet_groups", "column: users.created_at"]
}
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Connection refused | Start server: `go run main.go` |
| 405 Method Not Allowed | Use `POST` not `GET` |
| Invalid JSON error | Check JSON syntax at jsonlint.com |
| Missing columns error | Set `force_migrate: true` |

---

## Data Models (Quick Reference)

### User
```json
{
  "user_id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "name": "Administrator",
  "type": "human"
}
```

### Wallet
```json
{
  "wallet_id": 1,
  "name": "My Wallet",
  "balance": 0.00,
  "is_enabled": true
}
```

### Category
```json
{
  "category_id": 1,
  "name": "Income",
  "type": "income",
  "is_enabled": true
}
```

---

## Database Info

- **Location:** `moneyplanner.db` (in project root)
- **Type:** SQLite
- **Size:** < 1MB
- **Auto-created:** Yes, on first `/api/init` call

---

## Swagger/API Documentation

Currently preparing Swagger UI integration at:
- `http://localhost:8080/swagger/index.html` (planned)

See `API_DOCUMENTATION.md` for complete details.

---

## Need More Details?

See **API_DOCUMENTATION.md** for:
- Full endpoint documentation
- All request/response formats
- Complete data model definitions
- Error handling guide
- Future planned endpoints
