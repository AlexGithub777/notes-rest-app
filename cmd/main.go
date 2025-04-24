package main

import (
	"log"

	"github.com/AlexGithub777/notes-rest-app/internal/db"
	"github.com/AlexGithub777/notes-rest-app/internal/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize and test database connection
	if err := db.SetupDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // You can restrict origins if needed
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	// Setup API routes
	routes.SetupRoutes(e)

	// Start server
	log.Println("🚀 Server running at http://localhost:8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
