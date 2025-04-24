package main

import (
	"log"

	"notes-rest-app/internal/db"
	"notes-rest-app/internal/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Connect to database
	if err := db.SetupDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.DB.Close()

	// Start Echo server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Setup routes
	routes.SetupRoutes(e)

	// Start server
	log.Println("Server running at http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
