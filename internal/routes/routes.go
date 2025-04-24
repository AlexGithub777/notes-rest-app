package routes

import (
	"github.com/AlexGithub777/notes-rest-app/internal/handlers"

	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(e *echo.Echo) {

	// web routes
	e.GET("/", handlers.LoginWebHandler)
	e.GET("/signup", handlers.SignUpWebHandler)
	e.GET("/login", handlers.LoginWebHandler)
	e.GET("/logout", handlers.LoginWebHandler)

	e.GET("/home", handlers.HomeHandler)

	// Api routes
	e.POST("api/login", handlers.LoginHandler)
	e.POST("api/signup", handlers.SignUpHandler)
	// group for notes api
	notes := e.Group("/api/notes")
	notes.GET("/categories", handlers.GetAllCategoriesHandler)
	notes.GET("", handlers.GetAllNotesHandler)
	notes.GET("/:id", handlers.GetNoteByIDHandler)
	notes.POST("", handlers.CreateNoteHandler)
	notes.PUT("/:id", handlers.UpdateNoteHandler)
	notes.DELETE("/:id", handlers.DeleteNoteHandler)
}
