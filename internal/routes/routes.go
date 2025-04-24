package routes

import (
	"github.com/AlexGithub777/notes-rest-app/internal/handlers"

	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(e *echo.Echo) {

	// group for notes
	notes := e.Group("/notes")

	notes.GET("", handlers.GetAllNotesHandler)
	notes.GET("/:id", handlers.GetNoteByIDHandler)
	notes.POST("", handlers.CreateNoteHandler)
	notes.PUT("/:id", handlers.UpdateNoteHandler)
	notes.DELETE("/:id", handlers.DeleteNoteHandler)
}
