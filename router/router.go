package router

import (
	"github.com/labstack/echo/v4"
	m "github.com/labstack/echo/v4/middleware"
	"github.com/shawn-10x/100pfps/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	e.Use(m.Logger())
	e.Use(m.Recover())
	e.Use(middleware.HandleError)

	g := e.Group("")
	BoardSetup(g)
	return e
}
