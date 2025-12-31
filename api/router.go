package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	initAPI "moneyplanner/api/init"
	usersAPI "moneyplanner/api/users"
	userWalletAPI "moneyplanner/api/userwallet"
	walletAPI "moneyplanner/api/wallet"

	categoriesAPI "moneyplanner/api/categories"
	personsAPI "moneyplanner/api/persons"
	transactionsAPI "moneyplanner/api/transactions"
	userWalletGroupAPI "moneyplanner/api/userwalletgroup"
	walletGroupAPI "moneyplanner/api/walletgroup"
	walletGroupWalletAPI "moneyplanner/api/walletgroupwallet"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(mux *http.ServeMux) {
	log.Println("Registering API routes...")

	// Init API endpoint
	mux.HandleFunc("/api/init", handleInit)

	// Init Done API endpoint (GET)
	mux.HandleFunc("/api/initdone", handleInitDone)

	// User CRUD API endpoints
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("/api/users/", handleUserDetail)

	// Wallet CRUD API endpoints
	mux.HandleFunc("/api/wallets", handleWallets)
	mux.HandleFunc("/api/wallets/", handleWalletDetail)

	// WalletGroup CRUD
	mux.HandleFunc("/api/walletgroups", handleWalletGroups)
	mux.HandleFunc("/api/walletgroups/", handleWalletGroupDetail)

	// Persons CRUD API endpoints
	mux.HandleFunc("/api/persons", handlePersons)
	mux.HandleFunc("/api/persons/", handlePersonDetail)

	// Transactions CRUD API endpoints
	mux.HandleFunc("/api/transactions", handleTransactions)
	mux.HandleFunc("/api/transactions/", handleTransactionDetail)

	log.Println("✓ API routes registered")
}

// handleInit handles database initialization
func handleInit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req initAPI.InitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call initialization
	resp := initAPI.InitializeDatabase(&req, "moneyplanner.db")

	// Set appropriate status code
	statusCode := http.StatusOK
	if !resp.Success {
		statusCode = http.StatusBadRequest
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// handleInitDone handles initialization status check
func handleInitDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call status check
	resp := initAPI.CheckInitStatus("moneyplanner.db")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleUsers handles user list and creation (POST /api/users)
func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		handleUserCreate(w, r)
	case http.MethodGet:
		handleUserList(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleUserCreate handles POST /api/users - Create a new user
func handleUserCreate(w http.ResponseWriter, r *http.Request) {
	var req usersAPI.UserCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Create the user
	result, err := usersAPI.SetupUserWithDefaults(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User created successfully",
		"data": map[string]interface{}{
			"user":   result.User,
			"wallet": result.Wallet,
		},
	})
}

// handleUserList handles GET /api/users - List all users
func handleUserList(w http.ResponseWriter, r *http.Request) {
	users, err := usersAPI.ListAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

// handleUserDetail handles user detail operations (GET, PUT, DELETE /api/users/{id})
// and subroutes like /api/users/{id}/wallets
func handleUserDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user ID from path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	userIDStr := parts[3]
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid user ID: " + err.Error(),
		})
		return
	}

	// Subroute: /api/users/{id}/wallets ...
	if len(parts) >= 5 && parts[4] == "wallets" {
		// /api/users/{id}/wallets
		if len(parts) == 5 || parts[5] == "" {
			switch r.Method {
			case http.MethodGet:
				handleUserWalletList(w, r, uint(userID))
			case http.MethodPut:
				handleUserWalletReplace(w, r, uint(userID))
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// /api/users/{id}/wallets/{walletId}
		walletIDStr := parts[5]
		walletID64, err := strconv.ParseUint(walletIDStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid wallet ID: " + err.Error(),
			})
			return
		}
		walletID := uint(walletID64)

		switch r.Method {
		case http.MethodPost:
			handleUserWalletAttach(w, r, uint(userID), walletID)
		case http.MethodDelete:
			handleUserWalletDetach(w, r, uint(userID), walletID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Subroute: /api/users/{id}/walletgroups
	if len(parts) >= 5 && parts[4] == "walletgroups" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handleUserWalletGroupList(w, r, uint(userID))
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleUserGet(w, r, uint(userID))
	case http.MethodPut:
		handleUserUpdate(w, r, uint(userID))
	case http.MethodDelete:
		handleUserDelete(w, r, uint(userID))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleUserGet handles GET /api/users/{id} - Get user by ID
func handleUserGet(w http.ResponseWriter, r *http.Request, userID uint) {
	user, err := usersAPI.GetUserByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// handleUserUpdate handles PUT /api/users/{id} - Update user
func handleUserUpdate(w http.ResponseWriter, r *http.Request, userID uint) {
	var req usersAPI.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	user, err := usersAPI.UpdateUser(userID, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

// handleUserDelete handles DELETE /api/users/{id} - Delete user
func handleUserDelete(w http.ResponseWriter, r *http.Request, userID uint) {
	err := usersAPI.DeleteUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}

// handleWallets handles wallet list and creation (POST /api/wallets)
func handleWallets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		handleWalletCreate(w, r)
	case http.MethodGet:
		handleWalletList(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleWalletCreate(w http.ResponseWriter, r *http.Request) {
	var req walletAPI.WalletCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	wallet, err := walletAPI.CreateWallet(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Wallet created successfully",
		"data":    wallet,
	})
}

func handleWalletList(w http.ResponseWriter, r *http.Request) {
	wallets, err := walletAPI.ListAllWallets()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Wallets retrieved successfully",
		"data":    wallets,
	})
}

// handleWalletDetail handles wallet detail operations (GET, PUT, DELETE /api/wallets/{id})
func handleWalletDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	walletIDStr := parts[3]
	walletID64, err := strconv.ParseUint(walletIDStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid wallet ID: " + err.Error()})
		return
	}
	walletID := uint(walletID64)

	// Subroute: /api/wallets/{walletId}/categories...
	if len(parts) >= 5 && parts[4] == "categories" {
		// /api/wallets/{walletId}/categories
		if len(parts) == 5 || parts[5] == "" {
			switch r.Method {
			case http.MethodGet:
				handleWalletCategoryList(w, r, walletID)
			case http.MethodPost:
				handleWalletCategoryCreate(w, r, walletID)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// /api/wallets/{walletId}/categories/tree
		if len(parts) == 6 && parts[5] == "tree" && r.Method == http.MethodGet {
			handleWalletCategoryTree(w, r, walletID)
			return
		}

		// /api/wallets/{walletId}/categories/{categoryId}
		if len(parts) == 6 {
			categoryIDStr := parts[5]
			categoryID64, err := strconv.ParseUint(categoryIDStr, 10, 32)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid category ID: " + err.Error()})
				return
			}
			categoryID := uint(categoryID64)

			switch r.Method {
			case http.MethodGet:
				handleWalletCategoryGet(w, r, categoryID)
			case http.MethodPut:
				handleWalletCategoryUpdate(w, r, categoryID)
			case http.MethodDelete:
				handleWalletCategoryDelete(w, r, categoryID)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// /api/wallets/{walletId}/categories/{categoryId}/sync-global
		if len(parts) == 7 && parts[6] == "sync-global" && r.Method == http.MethodPost {
			categoryIDStr := parts[5]
			categoryID64, err := strconv.ParseUint(categoryIDStr, 10, 32)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid category ID: " + err.Error()})
				return
			}
			categoryID := uint(categoryID64)
			handleWalletCategorySyncGlobal(w, r, categoryID)
			return
		}
	}

	// Subroute: /api/wallets/{walletId}/transactions...
	if len(parts) >= 5 && parts[4] == "transactions" {
		// /api/wallets/{walletId}/transactions
		if len(parts) == 5 || parts[5] == "" {
			switch r.Method {
			case http.MethodGet:
				handleWalletTransactionList(w, r, walletID)
			case http.MethodPost:
				handleWalletTransactionCreate(w, r, walletID)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// /api/wallets/{walletId}/transactions/{transactionId}
		if len(parts) == 6 {
			transactionIDStr := parts[5]
			transactionID64, err := strconv.ParseUint(transactionIDStr, 10, 32)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid transaction ID: " + err.Error()})
				return
			}
			transactionID := uint(transactionID64)

			switch r.Method {
			case http.MethodGet:
				handleWalletTransactionGet(w, r, walletID, transactionID)
			case http.MethodPut:
				handleWalletTransactionUpdate(w, r, walletID, transactionID)
			case http.MethodDelete:
				handleWalletTransactionDelete(w, r, walletID, transactionID)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		handleWalletGet(w, r, walletID)
	case http.MethodPut:
		handleWalletUpdate(w, r, walletID)
	case http.MethodDelete:
		handleWalletDelete(w, r, walletID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleWalletGet(w http.ResponseWriter, r *http.Request, walletID uint) {
	wallet, err := walletAPI.GetWalletByID(walletID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Wallet retrieved successfully",
		"data":    wallet,
	})
}

func handleWalletUpdate(w http.ResponseWriter, r *http.Request, walletID uint) {
	var req walletAPI.WalletUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	wallet, err := walletAPI.UpdateWallet(walletID, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Wallet updated successfully",
		"data":    wallet,
	})
}

func handleWalletDelete(w http.ResponseWriter, r *http.Request, walletID uint) {
	if err := walletAPI.DeleteWallet(walletID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Wallet deleted successfully",
	})
}

func handleUserWalletAttach(w http.ResponseWriter, r *http.Request, userID, walletID uint) {
	if err := userWalletAPI.AttachWalletToUser(userID, walletID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Wallet attached to user",
	})
}

func handleUserWalletDetach(w http.ResponseWriter, r *http.Request, userID, walletID uint) {
	if err := userWalletAPI.DetachWalletFromUser(userID, walletID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Wallet detached from user",
	})
}

func handleUserWalletList(w http.ResponseWriter, r *http.Request, userID uint) {
	wallets, err := userWalletAPI.ListUserWallets(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User wallets retrieved successfully",
		"data":    wallets,
	})
}

func handleUserWalletReplace(w http.ResponseWriter, r *http.Request, userID uint) {
	var req userWalletAPI.ReplaceWalletsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := userWalletAPI.ReplaceUserWallets(userID, req.WalletIDs); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User wallets replaced successfully",
	})
}

func handleWalletGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		handleWalletGroupCreate(w, r)
	case http.MethodGet:
		handleWalletGroupList(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleWalletGroupCreate(w http.ResponseWriter, r *http.Request) {
	var req walletGroupAPI.WalletGroupCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	wg, err := walletGroupAPI.CreateWalletGroup(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Wallet group created successfully",
		"data":    wg,
	})
}

func handleWalletGroupList(w http.ResponseWriter, r *http.Request) {
	groups, err := walletGroupAPI.ListAllWalletGroups()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Wallet groups retrieved successfully",
		"data":    groups,
	})
}

func handleWalletGroupDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	groupIDStr := parts[3]
	groupID64, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid wallet group ID: " + err.Error()})
		return
	}
	groupID := uint(groupID64)

	// Subroute: /api/walletgroups/{id}/wallets...
	if len(parts) >= 5 && parts[4] == "wallets" {
		// /api/walletgroups/{id}/wallets
		if len(parts) == 5 || parts[5] == "" {
			switch r.Method {
			case http.MethodGet:
				handleWalletGroupWalletList(w, r, groupID)
			case http.MethodPut:
				handleWalletGroupWalletReplace(w, r, groupID)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// /api/walletgroups/{id}/wallets/{walletId}
		walletIDStr := parts[5]
		walletID64, err := strconv.ParseUint(walletIDStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid wallet ID: " + err.Error()})
			return
		}
		walletID := uint(walletID64)

		switch r.Method {
		case http.MethodPost:
			handleWalletGroupWalletAttach(w, r, groupID, walletID)
		case http.MethodDelete:
			handleWalletGroupWalletDetach(w, r, groupID, walletID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// /api/walletgroups/{id}
	switch r.Method {
	case http.MethodGet:
		wg, err := walletGroupAPI.GetWalletGroupByID(groupID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Wallet group retrieved successfully", "data": wg})

	case http.MethodPut:
		var req walletGroupAPI.WalletGroupUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
			return
		}
		wg, err := walletGroupAPI.UpdateWalletGroup(groupID, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Wallet group updated successfully", "data": wg})

	case http.MethodDelete:
		if err := walletGroupAPI.DeleteWalletGroup(groupID); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Wallet group deleted successfully"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleWalletGroupWalletAttach(w http.ResponseWriter, r *http.Request, groupID, walletID uint) {
	if err := walletGroupWalletAPI.AttachWalletToGroup(groupID, walletID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Wallet attached to wallet group"})
}

func handleWalletGroupWalletDetach(w http.ResponseWriter, r *http.Request, groupID, walletID uint) {
	if err := walletGroupWalletAPI.DetachWalletFromGroup(groupID, walletID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Wallet detached from wallet group"})
}

func handleWalletGroupWalletList(w http.ResponseWriter, r *http.Request, groupID uint) {
	wallets, err := walletGroupWalletAPI.ListWalletsInGroup(groupID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Wallets in group retrieved successfully", "data": wallets})
}

func handleWalletGroupWalletReplace(w http.ResponseWriter, r *http.Request, groupID uint) {
	var req walletGroupWalletAPI.ReplaceWalletsInGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := walletGroupWalletAPI.ReplaceWalletsInGroup(groupID, req.WalletIDs); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Wallet group wallets replaced successfully"})
}

func handleUserWalletGroupList(w http.ResponseWriter, r *http.Request, userID uint) {
	groups, err := userWalletGroupAPI.ListWalletGroupsForUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User wallet groups retrieved successfully",
		"data":    groups,
	})
}

// Category handlers

func handleWalletCategoryCreate(w http.ResponseWriter, r *http.Request, walletID uint) {
	var req categoriesAPI.CategoryCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Set the wallet ID from path
	req.WalletID = walletID

	category, err := categoriesAPI.CreateCategory(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Category created successfully",
		"data":    category,
	})
}

func handleWalletCategoryList(w http.ResponseWriter, r *http.Request, walletID uint) {
	categories, err := categoriesAPI.ListCategoriesByWallet(walletID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Categories retrieved successfully",
		"data":    categories,
	})
}

func handleWalletCategoryTree(w http.ResponseWriter, r *http.Request, walletID uint) {
	tree, err := categoriesAPI.GetCategoryTree(walletID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Category tree retrieved successfully",
		"data":    tree,
	})
}

func handleWalletCategoryGet(w http.ResponseWriter, r *http.Request, categoryID uint) {
	category, err := categoriesAPI.GetCategoryByID(categoryID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Category retrieved successfully",
		"data":    category,
	})
}

func handleWalletCategoryUpdate(w http.ResponseWriter, r *http.Request, categoryID uint) {
	var req categoriesAPI.CategoryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	category, err := categoriesAPI.UpdateCategory(categoryID, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Category updated successfully",
		"data":    category,
	})
}

func handleWalletCategoryDelete(w http.ResponseWriter, r *http.Request, categoryID uint) {
	if err := categoriesAPI.DeleteCategory(categoryID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Category deleted successfully",
	})
}

func handleWalletCategorySyncGlobal(w http.ResponseWriter, r *http.Request, categoryID uint) {
	if err := categoriesAPI.SyncGlobalCategoryToAllWallets(categoryID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Global category synced to all wallets",
	})
}

// Transaction handlers for wallet-specific endpoints

func handleWalletTransactionCreate(w http.ResponseWriter, r *http.Request, walletID uint) {
	var req transactionsAPI.TransactionCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Set wallet ID from path
	req.WalletID = walletID

	transaction, err := transactionsAPI.CreateTransaction(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Transaction created successfully",
		"data":    transaction,
	})
}

func handleWalletTransactionList(w http.ResponseWriter, r *http.Request, walletID uint) {
	transactions, err := transactionsAPI.ListTransactionsByWallet(walletID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Transactions retrieved successfully",
		"data":    transactions,
	})
}

func handleWalletTransactionGet(w http.ResponseWriter, r *http.Request, walletID, transactionID uint) {
	transaction, err := transactionsAPI.GetTransactionByID(transactionID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Check if transaction belongs to the wallet
	if transaction.WalletID != walletID {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Transaction not found in this wallet"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Transaction retrieved successfully",
		"data":    transaction,
	})
}

func handleWalletTransactionUpdate(w http.ResponseWriter, r *http.Request, walletID, transactionID uint) {
	// First check if transaction exists and belongs to wallet
	transaction, err := transactionsAPI.GetTransactionByID(transactionID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if transaction.WalletID != walletID {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Transaction not found in this wallet"})
		return
	}

	var req transactionsAPI.TransactionUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Prevent changing wallet ID in wallet-specific endpoint
	if req.WalletID != nil && *req.WalletID != walletID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot change wallet ID in wallet-specific endpoint"})
		return
	}

	transaction, err = transactionsAPI.UpdateTransaction(transactionID, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Transaction updated successfully",
		"data":    transaction,
	})
}

func handleWalletTransactionDelete(w http.ResponseWriter, r *http.Request, walletID, transactionID uint) {
	// First check if transaction exists and belongs to wallet
	transaction, err := transactionsAPI.GetTransactionByID(transactionID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if transaction.WalletID != walletID {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Transaction not found in this wallet"})
		return
	}

	err = transactionsAPI.DeleteTransaction(transactionID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Transaction deleted successfully",
	})
}

// handlePersons handles person list and creation (POST /api/persons)
func handlePersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		handlePersonCreate(w, r)
	case http.MethodGet:
		handlePersonList(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handlePersonCreate handles POST /api/persons - Create a new person
func handlePersonCreate(w http.ResponseWriter, r *http.Request) {
	var req personsAPI.PersonCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	person, err := personsAPI.CreatePerson(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Person created successfully",
		"data":    person,
	})
}

// handlePersonList handles GET /api/persons - List all persons
func handlePersonList(w http.ResponseWriter, r *http.Request) {
	persons, err := personsAPI.ListAllPersons()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Persons retrieved successfully",
		"data":    persons,
	})
}

// handlePersonDetail handles person detail operations (GET, PUT, DELETE /api/persons/{id})
func handlePersonDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	personIDStr := parts[3]
	personID, err := strconv.ParseUint(personIDStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid person ID: " + err.Error(),
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		handlePersonGet(w, r, uint(personID))
	case http.MethodPut:
		handlePersonUpdate(w, r, uint(personID))
	case http.MethodDelete:
		handlePersonDelete(w, r, uint(personID))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handlePersonGet handles GET /api/persons/{id} - Get person by ID
func handlePersonGet(w http.ResponseWriter, r *http.Request, personID uint) {
	person, err := personsAPI.GetPersonByID(personID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Person retrieved successfully",
		"data":    person,
	})
}

// handlePersonUpdate handles PUT /api/persons/{id} - Update person
func handlePersonUpdate(w http.ResponseWriter, r *http.Request, personID uint) {
	var req personsAPI.PersonUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	person, err := personsAPI.UpdatePerson(personID, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Person updated successfully",
		"data":    person,
	})
}

// handlePersonDelete handles DELETE /api/persons/{id} - Delete person
func handlePersonDelete(w http.ResponseWriter, r *http.Request, personID uint) {
	err := personsAPI.DeletePerson(personID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Person deleted successfully",
	})
}

// handleTransactions handles transaction list and creation (POST /api/transactions)
func handleTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		handleTransactionCreate(w, r)
	case http.MethodGet:
		handleTransactionList(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTransactionCreate handles POST /api/transactions - Create a new transaction
func handleTransactionCreate(w http.ResponseWriter, r *http.Request) {
	var req transactionsAPI.TransactionCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	transaction, err := transactionsAPI.CreateTransaction(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Transaction created successfully",
		"data":    transaction,
	})
}

// handleTransactionList handles GET /api/transactions - List all transactions
func handleTransactionList(w http.ResponseWriter, r *http.Request) {
	transactions, err := transactionsAPI.ListAllTransactions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Transactions retrieved successfully",
		"data":    transactions,
	})
}

// handleTransactionDetail handles transaction detail operations (GET, PUT, DELETE /api/transactions/{id})
func handleTransactionDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	transactionIDStr := parts[3]
	transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid transaction ID: " + err.Error(),
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleTransactionGet(w, r, uint(transactionID))
	case http.MethodPut:
		handleTransactionUpdate(w, r, uint(transactionID))
	case http.MethodDelete:
		handleTransactionDelete(w, r, uint(transactionID))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTransactionGet handles GET /api/transactions/{id} - Get transaction by ID
func handleTransactionGet(w http.ResponseWriter, r *http.Request, transactionID uint) {
	transaction, err := transactionsAPI.GetTransactionByID(transactionID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Transaction retrieved successfully",
		"data":    transaction,
	})
}

// handleTransactionUpdate handles PUT /api/transactions/{id} - Update transaction
func handleTransactionUpdate(w http.ResponseWriter, r *http.Request, transactionID uint) {
	var req transactionsAPI.TransactionUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	transaction, err := transactionsAPI.UpdateTransaction(transactionID, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Transaction updated successfully",
		"data":    transaction,
	})
}

// handleTransactionDelete handles DELETE /api/transactions/{id} - Delete transaction
func handleTransactionDelete(w http.ResponseWriter, r *http.Request, transactionID uint) {
	err := transactionsAPI.DeleteTransaction(transactionID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Transaction deleted successfully",
	})
}
