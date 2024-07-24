package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRules(r *echo.Group) {
	r.GET("/rules/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "rules.html", nil)
	})
}
