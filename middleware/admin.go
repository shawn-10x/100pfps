package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/model"
)

func CheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		adminTokenCookie, err := c.Cookie("admin_token")
		var admin *model.Admin
		if err == nil && adminTokenCookie != nil {
			admin = model.GetAdminByToken(adminTokenCookie.Value)
		}
		c.Set("admin", admin)
		return next(c)
	}
}
