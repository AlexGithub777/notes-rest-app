package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"notes-rest-app/db"
	// "notes-rest-app/routes"
)

func main() {
	// Connect to database
	if err := db.SetupDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Start Echo server
	e := echo.New()
	// use routes

	log.Println("Server running at http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
