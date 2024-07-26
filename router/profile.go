package router

import (
	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/handler"
)

func SetupProfile(r *echo.Group) {
	r.GET("/", handler.GetProfiles)

	g := r.Group("profile/")
	g.POST("create/", handler.CreateProfile)
	g.POST("delete/", handler.DeleteProfile)
	g.POST("banip/", handler.BanIP)
}
