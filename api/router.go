package api

import (
	"encoding/json"
	"log"
	"net/http"

	initAPI "moneyplanner/api/init"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(mux *http.ServeMux) {
	log.Println("Registering API routes...")

	// Init API endpoint
	mux.HandleFunc("/api/init", handleInit)

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
