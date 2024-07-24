package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupNotFound() {
	echo.NotFoundHandler = func(c echo.Context) error {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
}
