package router

import (
	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/handler"
)

func SetupBoard(r *echo.Group) {
	r.GET("/", handler.GetBoard)
	r.POST("/", handler.PostProfile)
}
