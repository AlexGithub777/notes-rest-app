// handlers/auth_handlers.go
package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/AlexGithub777/notes-rest-app/internal/db"
	"github.com/AlexGithub777/notes-rest-app/internal/models"
	"github.com/AlexGithub777/notes-rest-app/internal/utils"
	"github.com/labstack/echo/v4"
)

func SignUpHandler(c echo.Context) error {
	user := new(models.User)

	// get the user info form the form dont use bind
	if err := c.FormValue("username"); err == "" {
		return utils.JSONError(c, http.StatusBadRequest, "Username is required")
	}

	if err := c.FormValue("password"); err == "" {
		return utils.JSONError(c, http.StatusBadRequest, "Password is required")
	}

	user.Username = c.FormValue("username")
	user.Password = c.FormValue("password")

	// Check if the username already exists
	var existingUser models.User
	err := db.DB.QueryRow("SELECT id FROM users WHERE username = $1", user.Username).Scan(&existingUser.ID)
	if err != nil && err != sql.ErrNoRows {
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to check username")
	}

	if existingUser.ID != 0 {
		return utils.JSONError(c, http.StatusConflict, "Username already exists")
	}

	// DONT HASH PASSWORD
	if len(user.Password) < 6 {
		return utils.JSONError(c, http.StatusBadRequest, "Password must be at least 6 characters long")
	}

	// Validate the username (this should be hashed and checked in a real app)
	if len(user.Username) < 3 {
		return utils.JSONError(c, http.StatusBadRequest, "Username must be at least 3 characters long")
	}
	// Validate the username (this should be hashed and checked in a real app)
	if len(user.Username) > 20 {
		return utils.JSONError(c, http.StatusBadRequest, "Username must be at most 20 characters long")
	}

	// Validate the password (this should be hashed and checked in a real app)
	if len(user.Password) > 20 {
		return utils.JSONError(c, http.StatusBadRequest, "Password must be at most 20 characters long")
	}

	// Validate the password (this should be hashed and checked in a real app)
	if len(user.Password) < 6 {
		return utils.JSONError(c, http.StatusBadRequest, "Password must be at least 6 characters long")

	}

	// Insert the user into the database
	var userID int
	err = db.DB.QueryRow(
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
		user.Username, user.Password,
	).Scan(&userID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to create user")
	}

	// Redirect to the login page
	return c.Redirect(http.StatusSeeOther, "/login")
}

// Login handles user login
func LoginHandler(c echo.Context) error {
	user := new(models.User)

	// get the user info form the form dont use bind
	if err := c.FormValue("username"); err == "" {
		return utils.JSONError(c, http.StatusBadRequest, "Username is required")
	}

	if err := c.FormValue("password"); err == "" {
		return utils.JSONError(c, http.StatusBadRequest, "Password is required")
	}

	user.Username = c.FormValue("username")
	user.Password = c.FormValue("password")

	// log the user
	fmt.Println("User: ", user)

	// Check if the user exists
	var dbUser models.User
	err := db.DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", user.Username).
		Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "Invalid credentials")
	}

	// Check the password (this should be hashed and checked in a real app)
	if dbUser.Password != user.Password {
		return utils.JSONError(c, http.StatusUnauthorized, "Invalid credentials")
	}

	// Set the user in the session
	// In a real app, you would use a session management library
	// and set the user ID in the session
	// For this example, we'll just set it in the context
	c.Set("user", dbUser)
	// Set a cookie to remember the user
	cookie := new(http.Cookie)
	cookie.Name = "user_id"
	cookie.Value = fmt.Sprintf("%d", dbUser.ID)
	cookie.Path = "/"
	cookie.MaxAge = 3600 // 1 hour
	c.SetCookie(cookie)

	// redirect to home page
	return c.Redirect(http.StatusSeeOther, "/home")
}
