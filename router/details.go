package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupDetails(e *echo.Group) {
	e.GET("/details/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "details.html", nil)
	})
}
