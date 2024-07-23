package middleware

import (
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/model"
	"github.com/shawn-10x/100pfps/utils"
)

func CheckBan(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isBanned, err := model.IsBanned(c.Get("ip").(net.IP))
		if err != nil {
			return err
		}
		if isBanned {
			return c.Render(http.StatusUnauthorized, "ban.html", utils.M{})
		}
		return next(c)
	}
}
