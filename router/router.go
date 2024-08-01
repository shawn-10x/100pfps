package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/static"
)

func New() *echo.Echo {
	e := echo.New()

	g := e.Group("")
	SetupProfile(g)
	SetupAdmin(g)

	g.GET("/privacy/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "privacy.html", nil)
	})
	g.GET("/rules/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "rules.html", nil)
	})
	g.GET("/details/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "details.html", nil)
	})
	g.GET("/robots.txt", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/static/robots.txt")
	})
	g.GET("/sitemap.xml", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/static/sitemap.xml")
	})
	g.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(static.FS)))))

	echo.NotFoundHandler = func(c echo.Context) error {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}

	return e
}
