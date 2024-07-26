package router

import (
	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/handler"
)

func SetupAdmin(r *echo.Group) {
	r.GET("/admin/", handler.AdminGet)
	r.POST("/admin/signin/", handler.AdminSignIn)
	r.POST("/admin/create/", handler.CreateAdmin)
}
