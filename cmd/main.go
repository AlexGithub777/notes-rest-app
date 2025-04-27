package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AlexGithub777/notes-rest-app/internal/db"
	"github.com/AlexGithub777/notes-rest-app/internal/routes"
	"github.com/AlexGithub777/notes-rest-app/internal/utils"
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

	// Setup static file serving
	e.Static("/static", "static")

	// Initialize the TemplateRenderer
	tmplRenderer := utils.NewTemplateRenderer()

	// Set the template renderer for Echo
	e.Renderer = tmplRenderer

	port := "8080" // Define the port to listen on

	localhost := "localhost" // Define the localhost address

	// HTTP listener is in a goroutine as it's blocking
	go func() {
		if err := e.Start(localhost + ":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// make clickable link
	log.Printf("Server started at http://%s:%s", localhost, port)

	// Setup a ctrl-c trap to ensure a graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Log the shutdown process
	log.Println("Shutting HTTP service down")
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	log.Println("Shutdown complete")
}
