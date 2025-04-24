package utils

import (
	"github.com/labstack/echo/v4"
)

func JSONError(c echo.Context, status int, msg string) error {
	return c.JSON(status, map[string]string{"error": msg})
}
