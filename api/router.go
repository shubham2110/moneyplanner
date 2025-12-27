package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	initAPI "moneyplanner/api/init"
	usersAPI "moneyplanner/api/users"
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
