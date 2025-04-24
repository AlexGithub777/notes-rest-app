package routes

import (
	"notes-rest-app/internal/handlers"

	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(e *echo.Echo) {

	// group for notes
	notes := e.Group("/notes")

	notes.GET("", handlers.GetNotes)
	notes.GET("/:id", handlers.GetNote)
	notes.POST("", handlers.CreateNote)
	notes.PUT("/:id", handlers.UpdateNote)
	notes.DELETE("/:id", handlers.DeleteNote)
}
