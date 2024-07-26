package router

import (
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	g := e.Group("")
	SetupProfile(g)
	SetupPrivacy(g)
	SetupRules(g)
	SetupDetails(g)
	SetupAdmin(g)
	SetupPublic(g)
	SetupNotFound()

	return e
}
