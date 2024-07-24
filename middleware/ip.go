package middleware

import (
	"math/rand"
	"net"
	"os"

	"github.com/labstack/echo/v4"
)

func GetIP(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := net.IP{}
		if os.Getenv("MODE") == "Debug" {
			ip = make(net.IP, 4)
			for i := 0; i < 4; i++ {
				ip[i] = uint8(rand.Intn(256))
			}
		} else {
			ip = net.ParseIP(c.RealIP())
		}
		c.Set("ip", ip)
		return next(c)
	}
}
