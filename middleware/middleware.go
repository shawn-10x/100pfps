package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddlewares(e *echo.Echo) {
	if os.Getenv("mode") == "Debug" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.Use(HandleError)
	e.Use(GetIP)
	e.Use(CheckBan)
}
