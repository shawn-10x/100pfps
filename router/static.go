package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/static"
)

func SetupPublic(g *echo.Group) {
	g.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(static.FS)))))
}
