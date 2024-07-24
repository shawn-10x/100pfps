package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupPrivacy(r *echo.Group) {
	r.GET("/privacy/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "privacy.html", nil)
	})
}
