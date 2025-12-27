package main

import (
	"log"
	"moneyplanner/database"
)

func main() {
	// Initialize the database
	err := database.InitDB("moneyplanner.db")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MoneyPlanner started successfully!")

	// Database is ready to use via database.DB
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
}
