package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// HomeHandler handles the home page
func HomeHandler(c echo.Context) error {
	// Render the home page
	return c.Render(http.StatusOK, "home.html", nil)
}

// LoginWebHandler handles the login page
func LoginWebHandler(c echo.Context) error {
	// Render the login page
	return c.Render(http.StatusOK, "login.html", nil)
}

// SignUpWebHandler handles the signup page
func SignUpWebHandler(c echo.Context) error {
	// Render the signup page
	return c.Render(http.StatusOK, "signup.html", nil)
}

// LogoutWebHandler handles the logout page
func LogoutWebHandler(c echo.Context) error {
	// Clear the session or token here if needed clear cookie
	cookie := new(http.Cookie)
	cookie.Name = "user_id"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/login")
}
