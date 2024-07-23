package router

import (
	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	middleware.SetupMiddlewares(e)

	g := e.Group("")
	BoardSetup(g)
	return e
}
