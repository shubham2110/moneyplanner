# Check if Server is Initiatated
```
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/initdone"
```

##  Sample Response:
```
 {
   "success": true,
   "message": "Initialization status retrieved",
   "init_done": true,
   "is_new_db": false
 }
```

# Init the Database
```
$body = @{
     force_migrate = $true
     default_wallet_name = "My Wallet"
     admin_username = "admin1"
     admin_password = "pass123"
     admin_email = "admin1@example.com"
     admin_name = "Administrator"
 } | ConvertTo-Json

$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/init" `
   -Method Post `
   -Headers @{"Content-Type"="application/json"} `
   -Body $body
```

##  Sample Response 1:
```
 {
   "success": true,
   "message": "Database initialized successfully",
   "data": {
     "database_is_new": true,
     "migration_required": false,
     "admin_user": {
       "user_id": 1,
       "username": "admin1",
       "name": "Administrator",
       "email": "admin1@example.com",
       "user_type": "human",
       "created_at": "2025-12-31T12:54:00Z",
       "updated_at": "2025-12-31T12:54:00Z"
     },
     "default_wallet": {
       "wallet_id": 1,
       "name": "My Wallet",
       "icon": "💰",
       "is_enabled": true,
       "balance": 0,
       "last_modified_time": "2025-12-31T12:54:00Z"
     },
     "default_wallet_group": {
       "wallet_group_id": 1,
       "wallet_group_name": "Default",
       "created_at": "2025-12-31T12:54:00Z",
       "updated_at": "2025-12-31T12:54:00Z"
     },
     "root_categories": [
       {
         "category_id": 1,
         "name": "Income",
         "icon": "💵",
         "parent_id": null,
         "wallet_id": 1,
         "is_global": false,
         "created_at": "2025-12-31T12:54:00Z",
         "updated_at": "2025-12-31T12:54:00Z"
       },
       {
         "category_id": 2,
         "name": "Expense",
         "icon": "💸",
         "parent_id": null,
         "wallet_id": 1,
         "is_global": false,
         "created_at": "2025-12-31T12:54:00Z",
         "updated_at": "2025-12-31T12:54:00Z"
       }
     ]
   }
 }
 ```

## Sample Response 2
```
 {
  "success": true,
  "message": "Database schema is in sync",
  "data": {
    "database_is_new": false,
    "migration_required": false
  }
 } 
```

#  Create User

```
$uri = "https://ml.xlr.ovh/api/users"

$body = @{
     username           = "johndoe1"
     name               = "John Doe"
     email              = "john.doe1@example.com"
     password           = "StrongPassword123!"
     user_type          = "regular"       # must match models.UserType
     wallet_name        = "Main Wallet"
     wallet_group_name  = "Personal"
     create_categories  = $true
 } | ConvertTo-Json -Depth 3

$result=Invoke-RestMethod `
     -Uri $uri `
     -Method Post `
     -Headers @{
         "Content-Type" = "application/json"
     } `
    -Body $body
```
## Sample Success Response (with create_categories = true):
```
 {
   "success": true,
   "message": "User created successfully",
   "data": {
     "user": {
       "user_id": 2,
       "username": "johndoe1",
       "name": "John Doe",
       "email": "john.doe1@example.com",
       "user_type": "regular",
       "created_at": "2025-12-31T07:30:00Z",
       "updated_at": "2025-12-31T07:30:00Z",
       "default_wallet_id": 2
     },
     "wallet": {
       "wallet_id": 2,
       "name": "Main Wallet",
       "icon": "💰",
       "is_enabled": true,
       "balance": 0,
       "last_modified_time": "2025-12-31T07:30:00Z"
     }
   }
 }
```

## Sample Success Response (with create_categories = false):
```
 {
   "success": true,
   "message": "User created successfully",
   "data": {
     "user": {
       "user_id": 3,
       "username": "janesmith",
       "name": "Jane Smith",
       "email": "jane.smith@example.com",
       "user_type": "regular",
       "created_at": "2025-12-31T07:35:00Z",
       "updated_at": "2025-12-31T07:35:00Z",
       "default_wallet_id": 3
     },
     "wallet": {
       "wallet_id": 3,
       "name": "Personal Wallet",
       "icon": "💰",
       "is_enabled": true,
       "balance": 0,
       "last_modified_time": "2025-12-31T07:35:00Z"
     }
   }
 }
```

## Sample Failure Response (Invalid JSON):
```
 {
   "error": "Invalid request body: invalid character '}' looking for beginning of value"
 }
```
## Sample Failure Response (Missing Required Field):
```
 {
   "error": "username is mandatory"
 }
```

## Sample Failure Response (Duplicate Username):
```
 {
   "error": "UNIQUE constraint failed: users.username"
 }  
```

# Get List of users
```
$result=Invoke-RestMethod -Uri https://ml.xlr.ovh/api/users
```

## Sample Response:
```
 {
   "success": true,
   "message": "Users retrieved successfully",
   "data": [
     {
       "user_id": 1,
       "username": "admin1",
       "name": "Administrator",
       "email": "admin1@example.com",
       "user_type": "human",
       "created_at": "2025-12-31T12:45:00Z",
       "updated_at": "2025-12-31T12:45:00Z"
     }
   ]
 }
```

# Modify a User
```ps1
$uri = "https://ml.xlr.ovh/api/users/5"
$body = @{
     username = "new_username"
 } | ConvertTo-Json

$result=Invoke-RestMethod `
     -Uri $uri `
     -Method Put `
     -ContentType "application/json" `
     -Body $body
$result | Convertto-json -depth 10
```
## Sample Response
```json
{
  "data": {
    "user_id": 1,
    "username": "admin",
    "name": "Administrator",
    "email": "admin@example.com",
    "password": "12345678",
    "type": "human",
    "default_wallet_id": 1
  },
  "message": "User updated successfully",
  "success": true
}
```

# Add user to wallet
```
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/users/1/wallets/2" -Method POST
```

## sample responses

```json
{
  "message": "Wallet attached to user",
  "success": true
}
```

```json
{
  "error": "wallet not found: record not found"
}
```

# Detach user from wallet

```
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/users/1/wallets/2" -Method DELETE
```

## Sample Response
```json
{
message:"Wallet detached from user"
success:true
}
```

```json
{
error:"user not found: record not found"
}
```

# List User1 wallet
```
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/users/1/wallets" -Method GET
```

## Sample Response 
```json
{
  "data": [
    {
      "wallet_id": 1,
      "name": "My Wallet",
      "icon": "💰",
      "is_enabled": true,
      "balance": 0,
      "last_modified_time": "2025-12-31T12:11:46.923282928Z"
    }
  ],
  "message": "User wallets retrieved successfully",
  "success": true
}
```

# Repalce user's wallet

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/users/1/wallets" `
  -Method PUT `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{ "wallet_ids": [1, 3, 5] }'
```

```json
{
  "error": "one or more wallet IDs not found"
}
```
```
{
  "message": "User wallets replaced successfully",
  "success": true
}
```

# Create New Wallet 
```ps1
$result=Invoke-RestMethod `
   -Uri "https://ml.xlr.ovh/api/wallets" `
   -Method POST `
   -Headers @{ "Content-Type" = "application/json" } `
   -Body '{
     "name": "Main Wallet",
     "icon": "X",
     "is_enabled": true,
     "balance": 1000
   }'
```
## Sample Response
```json
{
  "data": {
    "wallet_id": 3,
    "name": "Main Wallet",
    "icon": "X",
    "is_enabled": true,
    "balance": 1000,
    "last_modified_time": "2025-12-31T13:31:29.913958796Z"
  },
  "message": "Wallet created successfully",
  "success": true
}
```

# Update wallet

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1" `
  -Method PUT `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "name": "Updated Wallet Name",
    "icon": "🏦",
    "is_enabled": true,
    "balance": 2500
  }'
```

## Sample Response
```json
{
  "data": {
    "wallet_id": 1,
    "name": "Updated Wallet Name",
    "icon": "🏦",
    "is_enabled": true,
    "balance": 2500,
    "last_modified_time": "2025-12-31T13:31:57.349818811Z"
  },
  "message": "Wallet updated successfully",
  "success": true
}
```


# Create walletgroup

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/walletgroups" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "wallet_group_name": "Family Expenses"
  }'
```

## sample Response
```
{
  "data": {
    "wallet_group_id": 2,
    "wallet_group_name": "Family Expenses"
  },
  "message": "Wallet group created successfully",
  "success": true
}
```
```
{
  "error": "wallet_group_name is required"
}
```

# List all wallet group

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/walletgroups" `
  -Method GET
```

##  Sample Response 
```json
{
  "data": [
    {
      "wallet_group_id": 1,
      "wallet_group_name": "Default"
    },
    {
      "wallet_group_id": 2,
      "wallet_group_name": "Family Expenses"
    }
  ],
  "message": "Wallet groups retrieved successfully",
  "success": true
}
```

# Get wallet group by ID 

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/walletgroups/1" `
  -Method GET
```

```json
{
  "data": {
    "wallet_group_id": 1,
    "wallet_group_name": "Default"
  },
  "message": "Wallet group retrieved successfully",
  "success": true
}
```

# Update Wallet Group Name

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/walletgroups/1" `
  -Method PUT `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "wallet_group_name": "Updated Family Group"
  }'
```

## Response
```json
{
  "data": {
    "wallet_group_id": 1,
    "wallet_group_name": "Updated Family Group"
  },
  "message": "Wallet group updated successfully",
  "success": true
}
```

# Delete wallet group
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/walletgroups/1" `
  -Method DELETE
```

```json
{
  "message": "Wallet group deleted successfully",
  "success": true
}
```

```json
{
  "error": "wallet group not found: record not found"
}
```


# Attach Wallet to Wallet Group
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/walletgroups/1/wallets/2" `
  -Method POST
```

## Response

```json
{
  "message": "Wallet attached to wallet group",
  "success": true
}
```

# Detach Wallet from Wallet Group
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/walletgroups/1/wallets/2" `
  -Method DELETE
```

## Response
```json
{
  "message": "Wallet detached from wallet group",
  "success": true
}
```

# List Wallets in a Wallet Group
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/walletgroups/1/wallets" `
  -Method GET
```
## Response
```
{
  "data": [
    {
      "wallet_id": 1,
      "name": "Updated Wallet Name",
      "icon": "🏦",
      "is_enabled": true,
      "balance": 2500,
      "last_modified_time": "2025-12-31T13:31:57.349818811Z"
    }
  ],
  "message": "Wallets in group retrieved successfully",
  "success": true
}
```
## Sample Response

```json
{
  "success": true,
  "message": "Wallets in group retrieved successfully",
  "data": [
    {
      "wallet_id": 1,
      "name": "Main Wallet",
      "icon": "🏦",
      "is_enabled": true,
      "balance": 1500.50,
      "last_modified_time": "2025-12-31T07:34:26.630Z"
    },
    {
      "wallet_id": 2,
      "name": "Savings Account",
      "icon": "💰",
      "is_enabled": true,
      "balance": 5000.00,
      "last_modified_time": "2025-12-30T10:15:00.000Z"
    }
  ]
}
```


# Replace Wallets in a Wallet Group (attach multiple at once)

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/walletgroups/1/wallets" `
  -Method PUT `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "wallet_ids": [1, 2, 3]
  }'
```

## Response
```json
{
  "message": "Wallet group wallets replaced successfully",
  "success": true
}
```
# List Wallet Groups for a User (via ORM preload)

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/users/1/walletgroups" `
  -Method GET
```

## Response
```json
{
  "data": [
    {
      "wallet_group_id": 1,
      "wallet_group_name": "Updated Family Group"
    }
  ],
  "message": "User wallet groups retrieved successfully",
  "success": true
}

```

# Create a Category for a Wallet
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/categories" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "name": "Groceries",
    "icon": "🛒",
    "parent_id": null,
    "is_global": false
  }'

## Response
```json
{
  "data": {
    "category_id": 13,
    "icon": "🛒",
    "name": "Groceries2",
    "parent_id": null,
    "root_id": 13,
    "wallet_id": 1,
    "is_global": false,
    "wallet": {
      "wallet_id": 0,
      "name": "",
      "icon": "",
      "is_enabled": false,
      "balance": 0,
      "last_modified_time": "0001-01-01T00:00:00Z"
    }
  },
  "message": "Category created successfully",
  "success": true
}
```

```json
{
  "error": "failed to create category: constraint failed: UNIQUE constraint failed: categories.name, categories.wallet_id (2067)"
}
```
# Create a Global Category (syncs to all wallets)

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/categories" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "name": "Transportation",
    "icon": "🚗",
    "parent_id": null,
    "is_global": true
  }'
```

## Response
```json
{
  "data": {
    "category_id": 14,
    "icon": "🚗",
    "name": "Transportation",
    "parent_id": null,
    "root_id": 14,
    "wallet_id": 1,
    "is_global": true,
    "wallet": {
      "wallet_id": 0,
      "name": "",
      "icon": "",
      "is_enabled": false,
      "balance": 0,
      "last_modified_time": "0001-01-01T00:00:00Z"
    }
  },
  "message": "Category created successfully",
  "success": true
}
```

# List Categories for a Wallet (flat list)
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/categories" `
  -Method GET
```
## Response
```json
{
  "data": [
    {
      "category_id": 1,
      "icon": "💵",
      "name": "Income",
      "parent_id": null,
      "root_id": 1,
      "wallet_id": 1,
      "is_global": true,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 2,
      "icon": "💸",
      "name": "Expense",
      "parent_id": null,
      "root_id": 2,
      "wallet_id": 1,
      "is_global": true,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 3,
      "icon": "📊",
      "name": "Salary",
      "parent_id": 1,
      "root_id": 1,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 4,
      "icon": "📊",
      "name": "Refund",
      "parent_id": 1,
      "root_id": 1,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 5,
      "icon": "📊",
      "name": "Bonus",
      "parent_id": 1,
      "root_id": 1,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 6,
      "icon": "📊",
      "name": "Interest",
      "parent_id": 1,
      "root_id": 1,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 7,
      "icon": "📉",
      "name": "Groceries",
      "parent_id": 2,
      "root_id": 2,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 8,
      "icon": "📉",
      "name": "House Maintenance",
      "parent_id": 2,
      "root_id": 2,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 9,
      "icon": "📉",
      "name": "Investment",
      "parent_id": 2,
      "root_id": 2,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 10,
      "icon": "📉",
      "name": "Utilities",
      "parent_id": 2,
      "root_id": 2,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 11,
      "icon": "📉",
      "name": "Transport",
      "parent_id": 2,
      "root_id": 2,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 12,
      "icon": "📉",
      "name": "Entertainment",
      "parent_id": 2,
      "root_id": 2,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 13,
      "icon": "🛒",
      "name": "Groceries2",
      "parent_id": null,
      "root_id": 13,
      "wallet_id": 1,
      "is_global": false,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    {
      "category_id": 14,
      "icon": "🚗",
      "name": "Transportation",
      "parent_id": null,
      "root_id": 14,
      "wallet_id": 1,
      "is_global": true,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    }
  ],
  "message": "Categories retrieved successfully",
  "success": true
}

```
# Get Category Tree for a Wallet (hierarchical)
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/categories/tree" `
  -Method GET
```
## Response
```json

{
  "data": {
    "roots": [
      {
        "category": {
          "category_id": 1,
          "icon": "💵",
          "name": "Income",
          "parent_id": null,
          "root_id": 1,
          "wallet_id": 1,
          "is_global": true,
          "wallet": {
            "wallet_id": 1,
            "name": "Updated Wallet Name",
            "icon": "🏦",
            "is_enabled": true,
            "balance": 2500,
            "last_modified_time": "2025-12-31T13:31:57.349818811Z"
          }
        },
        "children": [
          {
            "category": {
              "category_id": 3,
              "icon": "📊",
              "name": "Salary",
              "parent_id": 1,
              "root_id": 1,
              "wallet_id": 1,
              "is_global": false,
              "wallet": {
                "wallet_id": 1,
                "name": "Updated Wallet Name",
                "icon": "🏦",
                "is_enabled": true,
                "balance": 2500,
                "last_modified_time": "2025-12-31T13:31:57.349818811Z"
              }
            },
            "children": null
          },
          {
            "category": {
              "category_id": 4,
              "icon": "📊",
              "name": "Refund",
              "parent_id": 1,
              "root_id": 1,
              "wallet_id": 1,
              "is_global": false,
              "wallet": {
                "wallet_id": 1,
                "name": "Updated Wallet Name",
                "icon": "🏦",
                "is_enabled": true,
                "balance": 2500,
                "last_modified_time": "2025-12-31T13:31:57.349818811Z"
              }
            },
            "children": null
          },
          {
            "category": {
              "category_id": 5,
              "icon": "📊",
              "name": "Bonus",
              "parent_id": 1,
              "root_id": 1,
              "wallet_id": 1,
              "is_global": false,
              "wallet": {
                "wallet_id": 1,
                "name": "Updated Wallet Name",
                "icon": "🏦",
                "is_enabled": true,
                "balance": 2500,
                "last_modified_time": "2025-12-31T13:31:57.349818811Z"
              }
            },
            "children": null
          },
          {
            "category": {
              "category_id": 6,
              "icon": "📊",
              "name": "Interest",
              "parent_id": 1,
              "root_id": 1,
              "wallet_id": 1,
              "is_global": false,
              "wallet": {
                "wallet_id": 1,
                "name": "Updated Wallet Name",
                "icon": "🏦",
                "is_enabled": true,
                "balance": 2500,
                "last_modified_time": "2025-12-31T13:31:57.349818811Z"
              }
            },
            "children": null
          }
        ]
      },
      {
        "category": {
          "category_id": 2,
          "icon": "💸",
          "name": "Expense",
          "parent_id": null,
          "root_id": 2,
          "wallet_id": 1,
          "is_global": true,
          "wallet": {
            "wallet_id": 1,
            "name": "Updated Wallet Name",
            "icon": "🏦",
            "is_enabled": true,
            "balance": 2500,
            "last_modified_time": "2025-12-31T13:31:57.349818811Z"
          }
        },
        "children": [
          {
            "category": {
              "category_id": 7,
              "icon": "📉",
              "name": "Groceries",
              "parent_id": 2,
              "root_id": 2,
              "wallet_id": 1,
              "is_global": false,
              "wallet": {
                "wallet_id": 1,
                "name": "Updated Wallet Name",
                "icon": "🏦",
                "is_enabled": true,
                "balance": 2500,
                "last_modified_time": "2025-12-31T13:31:57.349818811Z"
              }
            },
            "children": null
          },
          {
            "category": {
              "category_id": 8,
              "icon": "📉",
              "name": "House Maintenance",
              "parent_id": 2,
              "root_id": 2,
              "wallet_id": 1,
              "is_global": false,
              "wallet": {
                "wallet_id": 1,
                "name": "Updated Wallet Name",
                "icon": "🏦",
                "is_enabled": true,
                "balance": 2500,
                "last_modified_time": "2025-12-31T13:31:57.349818811Z"
              }
            },
            "children": null
          },
          {
            "category": {
              "category_id": 9,
              "icon": "📉",
              "name": "Investment",
              "parent_id": 2,
              "root_id": 2,
              "wallet_id": 1,
              "is_global": false,
              "wallet": {
                "wallet_id": 1,
                "name": "Updated Wallet Name",
                "icon": "🏦",
                "is_enabled": true,
                "balance": 2500,
                "last_modified_time": "2025-12-31T13:31:57.349818811Z"
              }
            },
            "children": null
          },
          {
            "category": {
              "category_id": 10,
              "icon": "📉",
              "name": "Utilities",
              "parent_id": 2,
              "root_id": 2,
              "wallet_id": 1,
              "is_global": false,
              "wallet": {
                "wallet_id": 1,
                "name": "Updated Wallet Name",
                "icon": "🏦",
                "is_enabled": true,
                "balance": 2500,
                "last_modified_time": "2025-12-31T13:31:57.349818811Z"
              }
            },
            "children": null
          },
          {
            "category": {
              "category_id": 11,
              "icon": "📉",
              "name": "Transport",
              "parent_id": 2,
              "root_id": 2,
              "wallet_id": 1,
              "is_global": false,
              "wallet": {
                "wallet_id": 1,
                "name": "Updated Wallet Name",
                "icon": "🏦",
                "is_enabled": true,
                "balance": 2500,
                "last_modified_time": "2025-12-31T13:31:57.349818811Z"
              }
            },
            "children": null
          },
          {
            "category": {
              "category_id": 12,
              "icon": "📉",
              "name": "Entertainment",
              "parent_id": 2,
              "root_id": 2,
              "wallet_id": 1,
              "is_global": false,
              "wallet": {
                "wallet_id": 1,
                "name": "Updated Wallet Name",
                "icon": "🏦",
                "is_enabled": true,
                "balance": 2500,
                "last_modified_time": "2025-12-31T13:31:57.349818811Z"
              }
            },
            "children": null
          }
        ]
      },
      {
        "category": {
          "category_id": 13,
          "icon": "🛒",
          "name": "Groceries2",
          "parent_id": null,
          "root_id": 13,
          "wallet_id": 1,
          "is_global": false,
          "wallet": {
            "wallet_id": 1,
            "name": "Updated Wallet Name",
            "icon": "🏦",
            "is_enabled": true,
            "balance": 2500,
            "last_modified_time": "2025-12-31T13:31:57.349818811Z"
          }
        },
        "children": null
      },
      {
        "category": {
          "category_id": 14,
          "icon": "🚗",
          "name": "Transportation",
          "parent_id": null,
          "root_id": 14,
          "wallet_id": 1,
          "is_global": true,
          "wallet": {
            "wallet_id": 1,
            "name": "Updated Wallet Name",
            "icon": "🏦",
            "is_enabled": true,
            "balance": 2500,
            "last_modified_time": "2025-12-31T13:31:57.349818811Z"
          }
        },
        "children": null
      }
    ]
  },
  "message": "Category tree retrieved successfully",
  "success": true
}
```

# Get a Specific Category
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/categories/5" `
  -Method GET
```
## Response
```json
{
  "data": {
    "category_id": 5,
    "icon": "📊",
    "name": "Bonus",
    "parent_id": 1,
    "root_id": 1,
    "wallet_id": 1,
    "is_global": false,
    "wallet": {
      "wallet_id": 1,
      "name": "Updated Wallet Name",
      "icon": "🏦",
      "is_enabled": true,
      "balance": 2500,
      "last_modified_time": "2025-12-31T13:31:57.349818811Z"
    }
  },
  "message": "Category retrieved successfully",
  "success": true
}

```

# Update a Category
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/categories/5" `
  -Method PUT `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "name": "Updated Groceries",
    "icon": "🛒",
    "is_global": true
  }'
```

## Response
```json
{
  "data": {
    "category_id": 5,
    "icon": "🛒",
    "name": "Updated Groceries",
    "parent_id": 1,
    "root_id": 1,
    "wallet_id": 1,
    "is_global": true,
    "wallet": {
      "wallet_id": 1,
      "name": "Updated Wallet Name",
      "icon": "🏦",
      "is_enabled": true,
      "balance": 2500,
      "last_modified_time": "2025-12-31T13:31:57.349818811Z"
    }
  },
  "message": "Category updated successfully",
  "success": true
}

```
# Delete a Category
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/categories/5" `
  -Method DELETE
```

## Response
```json
{
  "message": "Category deleted successfully",
  "success": true
}
```

# Sync Global Category to All Wallets

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/categories/5/sync-global" `
  -Method POST
```
## Response
```json
{
  "message": "Global category synced to all wallets",
  "success": true
}
```
```json
{
  "error": "category not found: record not found"
}
```

# Create a Person

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/persons" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "person_name": "John Smith",
    "alias": "Johnny"
  }'
```

## Response
```json
{
  "data": {
    "person_id": 1,
    "person_name": "John Smith",
    "alias": "Johnny"
  },
  "message": "Person created successfully",
  "success": true
}
```

# List All Persons
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/persons" `
  -Method GET
```

## Response
```json
{
  "data": [
    {
      "person_id": 1,
      "person_name": "John Smith",
      "alias": "Johnny"
    },
    {
      "person_id": 2,
      "person_name": "John Smith1",
      "alias": "Johnny1"
    }
  ],
  "message": "Persons retrieved successfully",
  "success": true
}
```

```json
{
  "data": [],
  "message": "Persons retrieved successfully",
  "success": true
}
```


# Get a Specific Person
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/persons/1" `
  -Method GET
```

## Response

```json
{
  "data": {
    "person_id": 1,
    "person_name": "John Smith",
    "alias": "Johnny"
  },
  "message": "Person retrieved successfully",
  "success": true
}
```
# Update a Person
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/persons/1" `
  -Method PUT `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "person_name": "Updated Name",
    "alias": "Updated Alias"
  }'
```
## Response

```
{
  "data": {
    "person_id": 1,
    "person_name": "Updated Name",
    "alias": "Updated Alias"
  },
  "message": "Person updated successfully",
  "success": true
}
```

# Delete a Person
```
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/persons/1" `
  -Method DELETE
```

## Response
```json
{
  "message": "Person deleted successfully",
  "success": true
}
```

```json
{
  "error": "person not found: record not found"
}
```
# Create a Transaction
```
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/transactions" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "wallet_id": 1,
    "category_id": 1,
    "amount": 50.00,
    "person_name": "New Person",
    "note": "Lunch expense",
    "user_id": 1
  }'
```

## Response
```
{
  "data": {
    "transaction_id": 13,
    "category_id": 1,
    "amount": 50,
    "note": "Lunch expense",
    "person_id": 1,
    "wallet_id": 1,
    "transaction_time": "2025-12-31T14:01:13.828584581Z",
    "entry_time": "2025-12-31T14:01:13.828584581Z",
    "last_modified_time": "2025-12-31T14:01:13.828584581Z",
    "user_id": 1,
    "category": {
      "category_id": 1,
      "icon": "💵",
      "name": "Income",
      "parent_id": null,
      "root_id": 1,
      "wallet_id": 1,
      "is_global": true,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2500,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    "person": {
      "person_id": 1,
      "person_name": "New Person",
      "alias": ""
    },
    "wallet": {
      "wallet_id": 1,
      "name": "Updated Wallet Name",
      "icon": "🏦",
      "is_enabled": true,
      "balance": 2500,
      "last_modified_time": "2025-12-31T13:31:57.349818811Z"
    },
    "user": {
      "user_id": 1,
      "username": "admin",
      "name": "Administrator",
      "email": "admin@example.com",
      "password": "12345678",
      "type": "human",
      "default_wallet_id": 1
    }
  },
  "message": "Transaction created successfully",
  "success": true
}
```

# List All Transactions

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/transactions" `
  -Method GET
```

## Response

```json
{
  "data": [
    {
      "transaction_id": 1,
      "category_id": 1,
      "amount": 0,
      "note": null,
      "person_id": null,
      "wallet_id": 1,
      "transaction_time": "2025-12-31T12:11:49.179424478Z",
      "entry_time": "2025-12-31T12:11:49.179424558Z",
      "last_modified_time": "2025-12-31T12:11:49.179424632Z",
      "user_id": 1,
      "category": {
        "category_id": 1,
        "icon": "💵",
        "name": "Income",
        "parent_id": null,
        "root_id": 1,
        "wallet_id": 1,
        "is_global": true,
        "wallet": {
          "wallet_id": 1,
          "name": "Updated Wallet Name",
          "icon": "🏦",
          "is_enabled": true,
          "balance": 2550,
          "last_modified_time": "2025-12-31T13:31:57.349818811Z"
        }
      },
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2550,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      },
      "user": {
        "user_id": 1,
        "username": "admin",
        "name": "Administrator",
        "email": "admin@example.com",
        "password": "12345678",
        "type": "human",
        "default_wallet_id": 1
      }
    },
    {
      "transaction_id": 2,
      "category_id": 2,
      "amount": 0,
      "note": null,
      "person_id": null,
      "wallet_id": 1,
      "transaction_time": "2025-12-31T12:11:49.343650622Z",
      "entry_time": "2025-12-31T12:11:49.343650712Z",
      "last_modified_time": "2025-12-31T12:11:49.343650786Z",
      "user_id": 1,
      "category": {
        "category_id": 2,
        "icon": "💸",
        "name": "Expense",
        "parent_id": null,
        "root_id": 2,
        "wallet_id": 1,
        "is_global": true,
        "wallet": {
          "wallet_id": 1,
          "name": "Updated Wallet Name",
          "icon": "🏦",
          "is_enabled": true,
          "balance": 2550,
          "last_modified_time": "2025-12-31T13:31:57.349818811Z"
        }
      },
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2550,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      },
      "user": {
        "user_id": 1,
        "username": "admin",
        "name": "Administrator",
        "email": "admin@example.com",
        "password": "12345678",
        "type": "human",
        "default_wallet_id": 1
      }
    },
    {
      "transaction_id": 3,
      "category_id": 3,
      "amount": 0,
      "note": null,
      "person_id": null,
      "wallet_id": 1,
      "transaction_time": "2025-12-31T12:11:49.594060642Z",
      "entry_time": "2025-12-31T12:11:49.594060723Z",
      "last_modified_time": "2025-12-31T12:11:49.594060797Z",
      "user_id": 1,
      "category": {
        "category_id": 3,
        "icon": "📊",
        "name": "Salary",
        "parent_id": 1,
        "root_id": 1,
        "wallet_id": 1,
        "is_global": false,
        "wallet": {
          "wallet_id": 1,
          "name": "Updated Wallet Name",
          "icon": "🏦",
          "is_enabled": true,
          "balance": 2550,
          "last_modified_time": "2025-12-31T13:31:57.349818811Z"
        }
      },
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2550,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      },
      "user": {
        "user_id": 1,
        "username": "admin",
        "name": "Administrator",
        "email": "admin@example.com",
        "password": "12345678",
        "type": "human",
        "default_wallet_id": 1
      }
    },
    {
      "transaction_id": 4,
      "category_id": 4,
      "amount": 0,
      "note": null,
      "person_id": null,
      "wallet_id": 1,
      "transaction_time": "2025-12-31T12:11:49.775723979Z",
      "entry_time": "2025-12-31T12:11:49.775724124Z",
      "last_modified_time": "2025-12-31T12:11:49.775724233Z",
      "user_id": 1,
      "category": {
        "category_id": 4,
        "icon": "📊",
        "name": "Refund",
        "parent_id": 1,
        "root_id": 1,
        "wallet_id": 1,
        "is_global": false,
        "wallet": {
          "wallet_id": 1,
          "name": "Updated Wallet Name",
          "icon": "🏦",
          "is_enabled": true,
          "balance": 2550,
          "last_modified_time": "2025-12-31T13:31:57.349818811Z"
        }
      },
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2550,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      },
      "user": {
        "user_id": 1,
        "username": "admin",
        "name": "Administrator",
        "email": "admin@example.com",
        "password": "12345678",
        "type": "human",
        "default_wallet_id": 1
      }
    },
    {
      "transaction_id": 5,
      "category_id": 5,
      "amount": 0,
      "note": null,
      "person_id": null,
      "wallet_id": 1,
      "transaction_time": "2025-12-31T12:11:49.95044152Z",
      "entry_time": "2025-12-31T12:11:49.950441653Z",
      "last_modified_time": "2025-12-31T12:11:49.95044176Z",
      "user_id": 1,
      "category": {
        "category_id": 0,
        "icon": "",
        "name": "",
        "parent_id": null,
        "root_id": 0,
        "wallet_id": 0,
        "is_global": false,
        "wallet": {
          "wallet_id": 0,
          "name": "",
          "icon": "",
          "is_enabled": false,
          "balance": 0,
          "last_modified_time": "0001-01-01T00:00:00Z"
        }
      },
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2550,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      },
      "user": {
        "user_id": 1,
        "username": "admin",
        "name": "Administrator",
        "email": "admin@example.com",
        "password": "12345678",
        "type": "human",
        "default_wallet_id": 1
      }
    }
  ],
  "message": "Transactions retrieved successfully",
  "success": true
}
```


# Get a Specific Transaction

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/transactions/1" `
  -Method GET
```

## Response
```json
{
  "data": {
    "transaction_id": 1,
    "category_id": 1,
    "amount": 0,
    "note": null,
    "person_id": null,
    "wallet_id": 1,
    "transaction_time": "2025-12-31T12:11:49.179424478Z",
    "entry_time": "2025-12-31T12:11:49.179424558Z",
    "last_modified_time": "2025-12-31T12:11:49.179424632Z",
    "user_id": 1,
    "category": {
      "category_id": 1,
      "icon": "💵",
      "name": "Income",
      "parent_id": null,
      "root_id": 1,
      "wallet_id": 1,
      "is_global": true,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2550,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    "wallet": {
      "wallet_id": 1,
      "name": "Updated Wallet Name",
      "icon": "🏦",
      "is_enabled": true,
      "balance": 2550,
      "last_modified_time": "2025-12-31T13:31:57.349818811Z"
    },
    "user": {
      "user_id": 1,
      "username": "admin",
      "name": "Administrator",
      "email": "admin@example.com",
      "password": "12345678",
      "type": "human",
      "default_wallet_id": 1
    }
  },
  "message": "Transaction retrieved successfully",
  "success": true
}
```

```json
{
  "error": "transaction not found: record not found"
}
```
# Update a Transaction
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/transactions/1" `
  -Method PUT `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "amount": 60.00,
    "note": "Updated note"
  }'
```
## Response
```json
{
  "data": {
    "transaction_id": 1,
    "category_id": 1,
    "amount": 60,
    "note": "Updated note",
    "person_id": null,
    "wallet_id": 1,
    "transaction_time": "2025-12-31T12:11:49.179424478Z",
    "entry_time": "2025-12-31T12:11:49.179424558Z",
    "last_modified_time": "2025-12-31T14:04:31.617300145Z",
    "user_id": 1,
    "category": {
      "category_id": 1,
      "icon": "💵",
      "name": "Income",
      "parent_id": null,
      "root_id": 1,
      "wallet_id": 1,
      "is_global": true,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2550,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    "wallet": {
      "wallet_id": 1,
      "name": "Updated Wallet Name",
      "icon": "🏦",
      "is_enabled": true,
      "balance": 2550,
      "last_modified_time": "2025-12-31T13:31:57.349818811Z"
    },
    "user": {
      "user_id": 1,
      "username": "admin",
      "name": "Administrator",
      "email": "admin@example.com",
      "password": "12345678",
      "type": "human",
      "default_wallet_id": 1
    }
  },
  "message": "Transaction updated successfully",
  "success": true
}
```

# Delete a Transaction
```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/transactions/1" `
  -Method DELETE
```
## Response
```json
{
  "message": "Transaction deleted successfully",
  "success": true
}
```
# Create a Transaction for a Specific Wallet

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/transactions" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "category_id": 1,
    "amount": 50.00,
    "person_name": "New Person",
    "note": "Lunch expense",
    "user_id": 1
  }'
```

## Response
```json
{
  "data": {
    "transaction_id": 13,
    "category_id": 1,
    "amount": 50,
    "note": "Lunch expense",
    "person_id": 1,
    "wallet_id": 1,
    "transaction_time": "2025-12-31T14:05:39.583390747Z",
    "entry_time": "2025-12-31T14:05:39.583390747Z",
    "last_modified_time": "2025-12-31T14:05:39.583390747Z",
    "user_id": 1,
    "category": {
      "category_id": 1,
      "icon": "💵",
      "name": "Income",
      "parent_id": null,
      "root_id": 1,
      "wallet_id": 1,
      "is_global": true,
      "wallet": {
        "wallet_id": 1,
        "name": "Updated Wallet Name",
        "icon": "🏦",
        "is_enabled": true,
        "balance": 2560,
        "last_modified_time": "2025-12-31T13:31:57.349818811Z"
      }
    },
    "person": {
      "person_id": 1,
      "person_name": "New Person",
      "alias": ""
    },
    "wallet": {
      "wallet_id": 1,
      "name": "Updated Wallet Name",
      "icon": "🏦",
      "is_enabled": true,
      "balance": 2560,
      "last_modified_time": "2025-12-31T13:31:57.349818811Z"
    },
    "user": {
      "user_id": 1,
      "username": "admin",
      "name": "Administrator",
      "email": "admin@example.com",
      "password": "12345678",
      "type": "human",
      "default_wallet_id": 1
    }
  },
  "message": "Transaction created successfully",
  "success": true
}
```

# List Transactions for a Specific Wallet

```ps1
$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/transactions" `
  -Method GET
```
## Response
Response is similar to "List All Transactions"

# Detailed Transaction Filtering Examples
Response for all these are also similar to "List all transactions" since these are are sending filtered transactions. 

## Filter by User ID
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?user_id=1" -Method GET

## Filter by Date Ranges (Transaction Time)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?start_transaction_time=2025-01-01T00:00:00Z&end_transaction_time=2025-12-31T23:59:59Z" -Method GET

## Filter by Entry Time Range
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?start_entry_time=2025-01-01T00:00:00Z" -Method GET

## Filter by Last Modified Time
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?end_last_modified_time=2025-12-31T23:59:59Z" -Method GET

## Filter by Categories (Multiple)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=1,2,3" -Method GET

## Filter by Single Category
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=5" -Method GET

## Filter by Wallet ID (Global Endpoint)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?wallet_id=1" -Method GET

## Filter by Person ID
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?person_id=2" -Method GET

## Filter by Amount (Greater Than)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?amount_op=gt&amount_value=100" -Method GET

## Filter by Amount (Less Than or Equal)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?amount_op=le&amount_value=50" -Method GET

## Filter by Amount (Equal)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?amount_op=eq&amount_value=25.50" -Method GET

## Filter by Note (Fuzzy Search)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=lunch" -Method GET

## Combined Filters Example
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?user_id=1&start_transaction_time=2025-01-01T00:00:00Z&category_ids=1,2&amount_op=gt&amount_value=10&fuzzy_note=expense" -Method GET

## Wallet-Specific Filters
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/transactions?start_entry_time=2025-01-01T00:00:00Z&person_id=1" -Method GET

# Get a Specific Transaction for a Wallet

$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/transactions/13" `
  -Method GET

# Update a Transaction for a Wallet

$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/transactions/13" `
  -Method PUT `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "amount": 60.00,
    "note": "Updated note"
  }'

# Delete a Transaction for a Wallet

$result=Invoke-RestMethod `
  -Uri "https://ml.xlr.ovh/api/wallets/1/transactions/13" `
  -Method DELETE

-----------------------------------------------------------------------------------------
# User Journey Workflows

## Journey 1: Complete Setup and Basic Transaction Management
# 1. Initialize Database
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/initdone"
$body = @{ force_migrate = $true; default_wallet_name = "Personal Wallet"; admin_username = "admin"; admin_password = "pass123"; admin_email = "admin@example.com"; admin_name = "Administrator" } | ConvertTo-Json
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/init" -Method Post -Headers @{"Content-Type"="application/json"} -Body $body

# 2. Create User with Wallet and Categories
$body = @{ username = "johndoe"; name = "John Doe"; email = "john@example.com"; password = "SecurePass123!"; user_type = "regular"; wallet_name = "Main Wallet"; wallet_group_name = "Personal"; create_categories = $true } | ConvertTo-Json -Depth 3
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/users" -Method Post -Headers @{"Content-Type"="application/json"} -Body $body

# 3. Create Additional Categories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"name": "Food & Dining", "icon": "🍽️", "is_global": false}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"name": "Transportation", "icon": "🚗", "is_global": false}'

# 4. Create Persons
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/persons" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"person_name": "Grocery Store", "alias": "Local Market"}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/persons" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"person_name": "Gas Station"}'

# 5. Add Transactions
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 45.50, "person_name": "Grocery Store", "note": "Weekly groceries", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 25.00, "person_name": "Gas Station", "note": "Fuel for car", "user_id": 1}'

# 6. View All Transactions
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method GET

# 7. Filter Transactions by Category
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=1" -Method GET

## Journey 2: Multi-Wallet Management and Transfers
# 1. Create Additional Wallet
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"name": "Savings Account", "icon": "💰", "is_enabled": true, "balance": 1000}'

# 2. Create Transfer Category
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"name": "Transfer", "icon": "↔️", "is_global": false}'

# 3. Add Transfer Transactions (Note: In real app, this would be paired)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 3, "amount": -200.00, "note": "Transfer to Savings", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 2, "category_id": 3, "amount": 200.00, "note": "Transfer from Main", "user_id": 1}'

# 4. Check Wallet Balances
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1" -Method GET
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/2" -Method GET

# 5. List Transactions by Wallet
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/transactions" -Method GET

## Journey 3: Monthly Expense Tracking and Reporting
# 1. Add Monthly Expenses
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 120.00, "note": "Rent payment", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 85.30, "note": "Electricity bill", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 60.00, "note": "Internet bill", "user_id": 1}'

# 2. Filter by Date Range (Current Month)
$startDate = (Get-Date -Day 1 -Hour 0 -Minute 0 -Second 0).ToString("yyyy-MM-ddTHH:mm:ssZ")
$endDate = (Get-Date -Day (Get-Date).Day -Hour 23 -Minute 59 -Second 59).ToString("yyyy-MM-ddTHH:mm:ssZ")
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?start_transaction_time=$startDate&end_transaction_time=$endDate" -Method GET

# 3. Filter by Amount Range (High expenses)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?amount_op=gt&amount_value=50" -Method GET

# 4. Search by Note Content
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=bill" -Method GET

## Journey 4: Category Management and Global Categories
# 1. Create Global Categories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"name": "Income", "icon": "💼", "is_global": true}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"name": "Entertainment", "icon": "🎬", "is_global": true}'

# 2. Check if Global Categories Synced to Other Wallets
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/2/categories" -Method GET

# 3. Add Income Transaction
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 4, "amount": 3000.00, "note": "Monthly salary", "user_id": 1}'

# 4. Add Entertainment Expense
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 75.00, "note": "Movie night", "user_id": 1}'

# 5. Filter by Global Category
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=5" -Method GET

## Journey 5: Data Cleanup and Maintenance
# 1. List All Transactions (Before Cleanup)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method GET

# 2. Update Incorrect Transaction
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions/1" -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"amount": 50.00, "note": "Corrected grocery amount"}'

# 3. Delete Duplicate/Stale Transaction
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions/2" -Method DELETE

# 4. Update Person Information
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/persons/1" -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"person_name": "Updated Store Name"}'

# 5. Check Final State
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method GET
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/persons" -Method GET
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1" -Method GET

-----------
Start writing from here: Additional Comprehensive API Usage Samples

# Additional Transaction Management Use Cases

## Use Case: Daily Expense Tracking
# Morning coffee
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 5.50, "note": "Morning coffee at cafe", "user_id": 1}'

# Lunch
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 12.75, "person_name": "Restaurant ABC", "note": "Business lunch", "user_id": 1}'

# Evening groceries
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 67.30, "person_name": "SuperMart", "note": "Weekly groceries", "user_id": 1}'

# Filter today's expenses
$today = Get-Date -Format "yyyy-MM-dd"
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?start_transaction_time=${today}T00:00:00Z&end_transaction_time=${today}T23:59:59Z" -Method GET

## Use Case: Budget Monitoring
# Set monthly budget categories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Monthly Budget", "icon": "📊", "is_global": false}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Savings Goal", "icon": "🎯", "is_global": false}'

# Track budget vs actual
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=1,2&start_transaction_time=2025-01-01T00:00:00Z&end_transaction_time=2025-01-31T23:59:59Z" -Method GET

## Use Case: Recurring Payments
# Monthly subscriptions
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 15.99, "person_name": "Netflix", "note": "Monthly subscription", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 12.99, "person_name": "Spotify", "note": "Premium subscription", "user_id": 1}'

# Filter subscriptions
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=subscription" -Method GET

## Use Case: Business Expense Tracking
# Client meetings
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 45.00, "person_name": "Client XYZ Corp", "note": "Business lunch meeting", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 120.00, "person_name": "Conference Center", "note": "Conference registration", "user_id": 1}'

# Business expense report
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=business&amount_op=gt&amount_value=20" -Method GET

## Use Case: Travel Expenses
# Flight booking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 450.00, "person_name": "Airline Co", "note": "Flight to New York", "user_id": 1}'

# Hotel
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 280.00, "person_name": "Hotel Chain", "note": "3 nights accommodation", "user_id": 1}'

# Travel per diem
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 75.00, "note": "Daily meal allowance", "user_id": 1}'

# Travel expense summary
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=travel&start_transaction_time=2025-01-01T00:00:00Z" -Method GET

## Use Case: Family Budget Sharing
# Family member expenses
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 25.00, "person_name": "Family Member A", "note": "Shared dinner", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 30.00, "person_name": "Family Member B", "note": "Movie tickets", "user_id": 1}'

# Filter by person
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?person_id=1" -Method GET

## Use Case: Investment Tracking
# Stock purchases
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 2, "category_id": 4, "amount": -500.00, "person_name": "Brokerage Firm", "note": "Stock purchase - TECH Corp", "user_id": 1}'

# Dividend income
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 2, "category_id": 4, "amount": 25.50, "person_name": "Brokerage Firm", "note": "Dividend payment", "user_id": 1}'

# Investment performance
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=4&wallet_id=2" -Method GET

## Use Case: Emergency Fund Management
# Emergency withdrawals
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 2, "category_id": 3, "amount": -300.00, "note": "Emergency car repair", "user_id": 1}'

# Emergency deposits
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 2, "category_id": 4, "amount": 150.00, "note": "Emergency fund contribution", "user_id": 1}'

## Use Case: Seasonal Expense Planning
# Holiday shopping
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 89.99, "person_name": "Online Store", "note": "Holiday gifts", "user_id": 1}'

# Seasonal decorations
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 45.00, "person_name": "Decor Store", "note": "Christmas decorations", "user_id": 1}'

# Holiday entertainment
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=5&fuzzy_note=holiday" -Method GET

## Use Case: Tax Preparation
# Tax-deductible expenses
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 200.00, "person_name": "Charity Org", "note": "Tax-deductible donation", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 350.00, "person_name": "Medical Center", "note": "Medical expenses", "user_id": 1}'

# Tax year summary
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?start_transaction_time=2024-01-01T00:00:00Z&end_transaction_time=2024-12-31T23:59:59Z&fuzzy_note=tax" -Method GET

## Use Case: Multi-User Household
# Shared household expenses
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 150.00, "note": "Monthly utilities split", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 150.00, "note": "Monthly utilities split", "user_id": 2}'

# Filter by user
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?user_id=1" -Method GET

## Use Case: Goal-Oriented Saving
# Savings transfers
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 3, "amount": -100.00, "note": "Vacation fund transfer", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 2, "category_id": 4, "amount": 100.00, "note": "Vacation savings deposit", "user_id": 1}'

# Goal progress
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?wallet_id=2&fuzzy_note=vacation" -Method GET

## Use Case: Cash Flow Analysis
# Income tracking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 4, "amount": 2500.00, "person_name": "Employer", "note": "Monthly salary", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 4, "amount": 150.00, "person_name": "Freelance Client", "note": "Project payment", "user_id": 1}'

# Expense tracking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 800.00, "person_name": "Landlord", "note": "Monthly rent", "user_id": 1}'

# Monthly cash flow
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?start_transaction_time=2025-01-01T00:00:00Z&end_transaction_time=2025-01-31T23:59:59Z" -Method GET

## Use Case: Debt Management
# Loan payments
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 250.00, "person_name": "Bank", "note": "Car loan payment", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 180.00, "person_name": "Credit Card", "note": "Minimum payment", "user_id": 1}'

# Debt payoff progress
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=loan&fuzzy_note=payment" -Method GET

## Use Case: Gift and Special Occasions
# Gift purchases
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 125.00, "person_name": "Gift Shop", "note": "Birthday gift", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 75.00, "person_name": "Florist", "note": "Anniversary flowers", "user_id": 1}'

# Special occasion spending
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=5&start_transaction_time=2025-01-01T00:00:00Z" -Method GET

## Use Case: Home Improvement Projects
# DIY supplies
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 89.99, "person_name": "Hardware Store", "note": "Paint and brushes for bedroom", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 45.50, "person_name": "Home Depot", "note": "Garden tools", "user_id": 1}'

# Project tracking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=home&fuzzy_note=improvement" -Method GET

## Use Case: Educational Expenses
# Course fees
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 299.00, "person_name": "Online Learning Platform", "note": "Programming course", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 75.00, "person_name": "Bookstore", "note": "Technical books", "user_id": 1}'

# Education budget
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=course&fuzzy_note=book" -Method GET

## Use Case: Health and Wellness
# Medical expenses
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 120.00, "person_name": "Doctor Office", "note": "Annual checkup", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 45.00, "person_name": "Pharmacy", "note": "Prescription medication", "user_id": 1}'

# Gym membership
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 50.00, "person_name": "Fitness Center", "note": "Monthly membership", "user_id": 1}'

# Health expense tracking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=medical&fuzzy_note=health" -Method GET

## Use Case: Vehicle Maintenance
# Car maintenance
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 85.00, "person_name": "Auto Shop", "note": "Oil change and tire rotation", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 150.00, "person_name": "Mechanic", "note": "Brake pad replacement", "user_id": 1}'

# Insurance
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 120.00, "person_name": "Insurance Co", "note": "Auto insurance premium", "user_id": 1}'

# Vehicle expense summary
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=car&fuzzy_note=auto&fuzzy_note=insurance" -Method GET

## Use Case: Pet Care Expenses
# Pet food and supplies
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 45.99, "person_name": "Pet Store", "note": "Dog food and treats", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 200.00, "person_name": "Vet Clinic", "note": "Annual checkup and vaccinations", "user_id": 1}'

# Pet care tracking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=pet&fuzzy_note=vet&fuzzy_note=dog" -Method GET

## Use Case: Technology and Gadgets
# Device purchases
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 299.99, "person_name": "Tech Store", "note": "New smartphone", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 49.99, "person_name": "Accessory Shop", "note": "Phone case and screen protector", "user_id": 1}'

# Software subscriptions
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 9.99, "person_name": "Cloud Service", "note": "Monthly storage subscription", "user_id": 1}'

# Tech spending analysis
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=phone&fuzzy_note=tech&fuzzy_note=software" -Method GET

## Use Case: Travel and Vacation
# Vacation planning
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 850.00, "person_name": "Travel Agency", "note": "Flight and hotel package", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 120.00, "person_name": "Restaurant", "note": "Pre-vacation dinner", "user_id": 1}'

# Vacation expenses
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 65.00, "person_name": "Gift Shop", "note": "Souvenirs", "user_id": 1}'

# Vacation budget tracking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=vacation&fuzzy_note=travel&fuzzy_note=flight" -Method GET

## Use Case: Home Office Setup
# Office supplies
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 199.99, "person_name": "Office Supply", "note": "Ergonomic chair", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 149.99, "person_name": "Computer Store", "note": "External monitor", "user_id": 1}'

# Internet and utilities for home office
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 75.00, "person_name": "ISP", "note": "High-speed internet", "user_id": 1}'

# Home office expense report
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=office&fuzzy_note=home&fuzzy_note=monitor" -Method GET

## Use Case: Charitable Donations
# Regular donations
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 50.00, "person_name": "Local Charity", "note": "Monthly donation", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 25.00, "person_name": "Animal Shelter", "note": "Pet adoption support", "user_id": 1}'

# Special campaigns
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 100.00, "person_name": "Red Cross", "note": "Disaster relief donation", "user_id": 1}'

# Donation tracking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=donation&fuzzy_note=charity" -Method GET

## Use Case: Insurance Premiums
# Various insurance payments
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 95.00, "person_name": "Insurance Company", "note": "Home insurance quarterly", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 45.00, "person_name": "Dental Insurance", "note": "Monthly dental coverage", "user_id": 1}'

# Life insurance
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 75.00, "person_name": "Life Insurance Co", "note": "Term life insurance", "user_id": 1}'

# Insurance expense summary
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=insurance" -Method GET

## Use Case: Child Education
# School supplies
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 85.50, "person_name": "School Store", "note": "Back to school supplies", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 120.00, "person_name": "Bookstore", "note": "Textbooks and materials", "user_id": 1}'

# Extracurricular activities
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 150.00, "person_name": "Music School", "note": "Piano lessons", "user_id": 1}'

# Education expense tracking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=school&fuzzy_note=education&fuzzy_note=book" -Method GET

## Use Case: Home Security
# Security system installation
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 499.99, "person_name": "Security Company", "note": "Home security system installation", "user_id": 1}'

# Monthly monitoring
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 39.99, "person_name": "Security Company", "note": "Monthly monitoring service", "user_id": 1}'

# Security expenses
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=security" -Method GET

## Use Case: Professional Development
# Conference fees
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 350.00, "person_name": "Conference Org", "note": "Industry conference registration", "user_id": 1}'

# Professional memberships
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 150.00, "person_name": "Professional Association", "note": "Annual membership fee", "user_id": 1}'

# Training materials
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 2, "amount": 89.99, "person_name": "Online Platform", "note": "Professional certification course", "user_id": 1}'

# Professional development tracking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=conference&fuzzy_note=professional&fuzzy_note=certification" -Method GET

## Use Case: Gardening and Outdoor
# Gardening supplies
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 67.89, "person_name": "Garden Center", "note": "Plants and soil", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 45.99, "person_name": "Hardware Store", "note": "Gardening tools", "user_id": 1}'

# Outdoor equipment
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 299.99, "person_name": "Outdoor Store", "note": "Camping gear", "user_id": 1}'

# Outdoor expenses
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?fuzzy_note=garden&fuzzy_note=outdoor&fuzzy_note=camping" -Method GET

## Use Case: Hobby and Recreation
# Hobby supplies
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 125.00, "person_name": "Art Supply Store", "note": "Painting materials", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 89.99, "person_name": "Craft Store", "note": "Knitting supplies", "user_id": 1}'

# Recreation activities
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 55.00, "person_name": "Recreation Center", "note": "Bowling night", "user_id": 1}'

# Hobby tracking
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=5&fuzzy_note=hobby&fuzzy_note=art&fuzzy_note=craft" -Method GET

# Error Examples (20 samples)

## Error: Invalid JSON
# Missing closing brace
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 50.00'  # Returns 400 Bad Request

## Error: Missing Required Fields
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"note": "Test transaction"}'  # Returns 400 - wallet_id, category_id, amount, user_id required

## Error: Invalid Wallet ID
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 999, "category_id": 1, "amount": 50.00, "user_id": 1}'  # Returns 400 - wallet not found

## Error: Invalid Category ID
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 999, "amount": 50.00, "user_id": 1}'  # Returns 400 - category not found

## Error: Invalid User ID
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 50.00, "user_id": 999}'  # Returns 400 - user not found

## Error: Transaction Not Found
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions/999" -Method GET  # Returns 404 Not Found

## Error: Person Not Found
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 50.00, "person_id": 999, "user_id": 1}'  # Returns 400 - person not found

## Error: Invalid Amount
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": "invalid", "user_id": 1}'  # Returns 400 - invalid amount format

## Error: Empty Person Name
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 50.00, "person_name": "", "user_id": 1}'  # Person creation ignored, transaction created without person

## Error: Invalid Date Format
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?start_transaction_time=invalid-date" -Method GET  # Returns 400 - invalid date format

## Error: Wallet-Specific Transaction Access
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/transactions/5" -Method GET  # Returns 404 if transaction doesn't belong to wallet 1

## Error: Invalid Amount Operator
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?amount_op=invalid&amount_value=100" -Method GET  # Ignores invalid operator, returns all transactions

## Error: Category Not Belonging to Wallet
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 999, "amount": 50.00, "user_id": 1}'  # Returns 400 - category not found or doesn't belong to wallet

## Error: Duplicate Category Name
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Food & Dining", "icon": "🍽️"}'  # Returns 400 - category name must be unique per wallet

## Error: Invalid Parent Category
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Subcategory", "icon": "📁", "parent_id": 999}'  # Returns 400 - parent category not found

## Error: Cannot Delete Category With Children
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories/1" -Method DELETE  # Returns 400 if category has subcategories

## Error: Cannot Delete Category With Transactions
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories/1" -Method DELETE  # Returns 400 if category has associated transactions

## Error: Invalid Wallet Group Name
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/walletgroups" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_group_name": ""}'  # Returns 400 - name required

## Error: Wallet Already in Group
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/walletgroups/1/wallets/1" -Method POST  # Returns 400 if wallet is already in the group

## Error: Self-Transfer in Wallet-Specific Endpoint
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/transactions/1" -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 2}'  # Returns 400 - cannot change wallet in wallet-specific endpoint

## Error: Database Connection Issues
# If database is down or connection fails, all endpoints return 500 Internal Server Error with connection error message

## Error: Concurrent Modification
# If multiple users modify the same transaction simultaneously, one will succeed and others may get 409 Conflict or database constraint errors

## Error: Invalid Filter Combinations
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?amount_op=gt&amount_value=invalid" -Method GET  # Returns 400 - invalid amount value format

# Category Hierarchy and Transaction Assignment Journeys

## Category Journey 1: Food and Dining Hierarchy
# 1. Create root category
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Food & Dining", "icon": "🍽️", "is_global": false}'

# 2. Create subcategories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Restaurants", "icon": "🍽️", "parent_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Fast Food", "icon": "🍔", "parent_id": 2}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Fine Dining", "icon": "🥂", "parent_id": 2}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Groceries", "icon": "🛒", "parent_id": 1}'

# 3. View category tree
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories/tree" -Method GET

# 4. Add transactions to different levels
# Root level transaction
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 1, "amount": 25.00, "note": "General food expense", "user_id": 1}'

# Subcategory transactions
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 3, "amount": 15.99, "person_name": "Burger Joint", "note": "Lunch burger", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 4, "amount": 85.50, "person_name": "Fine Restaurant", "note": "Dinner date", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 5, "amount": 67.30, "person_name": "SuperMart", "note": "Weekly groceries", "user_id": 1}'

# 5. Filter by specific subcategory
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=3" -Method GET

# 6. Filter by parent category (should include all subcategories)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=1,2" -Method GET

## Category Journey 2: Transportation Expense Hierarchy
# 1. Create transportation root
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Transportation", "icon": "🚗", "is_global": false}'

# 2. Create transportation subcategories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Fuel", "icon": "⛽", "parent_id": 6}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Public Transport", "icon": "🚇", "parent_id": 6}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Parking", "icon": "🅿️", "parent_id": 6}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Maintenance", "icon": "🔧", "parent_id": 6}'

# 3. Add maintenance sub-subcategories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Oil Change", "icon": "🛢️", "parent_id": 10}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Tire Service", "icon": "🛞", "parent_id": 10}'

# 4. View full tree
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories/tree" -Method GET

# 5. Add transactions at various levels
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 6, "amount": 45.00, "note": "General transport", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 7, "amount": 55.00, "person_name": "Gas Station", "note": "Fuel", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 8, "amount": 12.00, "person_name": "Transit Authority", "note": "Bus pass", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 9, "amount": 8.00, "note": "Downtown parking", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 11, "amount": 75.00, "person_name": "Auto Shop", "note": "Oil change", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 12, "amount": 150.00, "person_name": "Tire Store", "note": "New tires", "user_id": 1}'

# 6. Filter by maintenance category (includes subcategories)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=10" -Method GET

# 7. Filter by specific leaf categories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=7,8,9" -Method GET

# 8. Monthly transportation summary
$startOfMonth = (Get-Date -Day 1 -Hour 0 -Minute 0 -Second 0).ToString("yyyy-MM-ddTHH:mm:ssZ")
$endOfMonth = (Get-Date -Day (Get-Date).Day -Hour 23 -Minute 59 -Second 59).ToString("yyyy-MM-ddTHH:mm:ssZ")
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=6&start_transaction_time=$startOfMonth&end_transaction_time=$endOfMonth" -Method GET

## Category Journey 3: Global Categories and Cross-Wallet Usage
# 1. Create global entertainment category
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Entertainment", "icon": "🎬", "is_global": true}'

# 2. Create entertainment subcategories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Movies", "icon": "🎥", "parent_id": 13}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Games", "icon": "🎮", "parent_id": 13}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Music", "icon": "🎵", "parent_id": 13}'

# 3. Check if global categories synced to wallet 2
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/2/categories" -Method GET

# 4. Add transactions using global categories from different wallets
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 14, "amount": 25.00, "note": "Movie ticket", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 2, "category_id": 15, "amount": 59.99, "note": "New video game", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 16, "amount": 9.99, "note": "Music streaming", "user_id": 1}'

# 5. Filter entertainment across all wallets
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=13" -Method GET

# 6. Sync global categories manually if needed
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories/13/sync-global" -Method POST

## Category Journey 4: Category Management and Restructuring
# 1. Create initial categories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Shopping", "icon": "🛍️", "is_global": false}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Clothing", "icon": "👕", "parent_id": 17}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Electronics", "icon": "📱", "parent_id": 17}'

# 2. Add transactions
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 18, "amount": 89.99, "note": "New shirt", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 19, "amount": 299.99, "note": "Smartphone", "user_id": 1}'

# 3. Restructure categories - move electronics to new parent
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Technology", "icon": "💻", "is_global": false}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories/19" -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"parent_id": 20}'

# 4. Update category names
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories/18" -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"name": "Fashion", "icon": "👗"}'

# 5. View updated tree
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories/tree" -Method GET

# 6. Transactions still accessible under new structure
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=17,20" -Method GET

## Category Journey 5: Budget Categories and Expense Tracking
# 1. Create budget-oriented categories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Monthly Budget", "icon": "📊", "is_global": false}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Under Budget", "icon": "✅", "parent_id": 21}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Over Budget", "icon": "⚠️", "parent_id": 21}'

# 2. Create spending categories
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Dining Out", "icon": "🍽️", "is_global": false}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/wallets/1/categories" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name": "Coffee Shops", "icon": "☕", "is_global": false}'

# 3. Add budgeted vs actual transactions
# Budget allocations (income category might be better, but using existing)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 22, "amount": 300.00, "note": "Monthly dining budget", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 23, "amount": 100.00, "note": "Monthly coffee budget", "user_id": 1}'

# Actual spending
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 24, "amount": 45.67, "note": "Restaurant dinner", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 24, "amount": 78.90, "note": "Weekend dining", "user_id": 1}'
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"wallet_id": 1, "category_id": 25, "amount": 12.50, "note": "Morning coffee", "user_id": 1}'

# 4. Monthly budget analysis
$startOfMonth = (Get-Date -Day 1 -Hour 0 -Minute 0 -Second 0).ToString("yyyy-MM-ddTHH:mm:ssZ")
$endOfMonth = (Get-Date -Day (Get-Date).Day -Hour 23 -Minute 59 -Second 59).ToString("yyyy-MM-ddTHH:mm:ssZ")
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=24,25&start_transaction_time=$startOfMonth&end_transaction_time=$endOfMonth" -Method GET

# 5. Budget vs actual comparison (would need frontend calculation)
$result=Invoke-RestMethod -Uri "https://ml.xlr.ovh/api/transactions?category_ids=21,22,23,24,25&start_transaction_time=$startOfMonth&end_transaction_time=$endOfMonth" -Method GET

