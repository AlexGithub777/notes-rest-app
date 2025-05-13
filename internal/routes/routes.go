package routes

import (
	"net/http"

	"github.com/AlexGithub777/notes-rest-app/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RequireLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Try to get user_id cookie
		cookie, err := c.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			// Not logged in
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		// You can also add DB lookup here if needed
		c.Set("user_id", cookie.Value) // Store user ID in context if needed
		return next(c)
	}
}

// SetupRoutes configures all the routes for the application
func SetupRoutes(e *echo.Echo) {

	// web routes
	e.GET("/", handlers.LoginWebHandler)
	e.GET("/signup", handlers.SignUpWebHandler)
	e.GET("/login", handlers.LoginWebHandler)
	e.GET("/logout", handlers.LoginWebHandler)

	e.GET("/home", handlers.HomeHandler, RequireLogin)
	e.GET("/all-notes", handlers.AllNotesHandler, RequireLogin)

	// Api routes
	e.POST("api/login", handlers.LoginHandler)
	e.POST("api/signup", handlers.SignUpHandler)
	// group for notes api
	notes := e.Group("/api/notes", RequireLogin)
	notes.GET("/categories", handlers.GetAllCategoriesHandler)
	// gets all notes for logged in user
	notes.GET("", handlers.GetAllNotesHandler)
	notes.GET("/:id", handlers.GetNoteByIDHandler)
	notes.POST("", handlers.CreateNoteHandler)
	notes.PUT("/:id", handlers.UpdateNoteHandler)
	notes.DELETE("/:id", handlers.DeleteNoteHandler)

	// get all notes for all users
	notes.GET("/all", handlers.GetAllNotesForAllUsersHandler)
}
