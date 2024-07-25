package router

import (
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	g := e.Group("")
	SetupBoard(g)
	SetupPrivacy(g)
	SetupRules(g)
	SetupDetails(g)
	SetupNotFound()
	return e
}
