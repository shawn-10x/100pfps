package middleware

import (
	"net"

	"github.com/labstack/echo/v4"
)

func GetIP(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("ip", net.ParseIP(c.RealIP()))
		return next(c)
	}
}
