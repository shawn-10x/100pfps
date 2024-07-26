package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/model"
)

func CheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		adminTokenCookie, err := c.Cookie("admin_token")
		if err == nil && adminTokenCookie != nil {
			c.Set("admin", model.GetAdminByToken(adminTokenCookie.Value))
		}
		return next(c)
	}
}
