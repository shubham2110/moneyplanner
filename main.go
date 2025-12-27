package main

import (
	"log"
	"net/http"

	"moneyplanner/api"
)

func main() {
	// Setup HTTP router
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	// Start server
	log.Println("MoneyPlanner server running on :8080")
	log.Println("Initialize database by calling:")
	log.Println("  POST http://localhost:8080/api/init")
	log.Println("  Body: {\"force_migrate\": false}")
	log.Println("")
	log.Println("For a new database with defaults:")
	log.Println("  Body: {\"force_migrate\": true, \"default_wallet_name\": \"My Wallet\"}")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
