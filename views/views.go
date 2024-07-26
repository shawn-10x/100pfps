package views

import (
	"embed"
	"encoding/base64"
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/utils"
)

//go:embed *.html
var FS embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if data == nil {
		data = utils.M{}
	}
	data.(utils.M)["view"] = name
	data.(utils.M)["path"] = c.Request().URL.Path
	return t.templates.ExecuteTemplate(w, name, data)
}

func SetupViews(e *echo.Echo) {
	funcMap := template.FuncMap{
		"hasKey":       hasKey,
		"ternary":      ternary,
		"valueOr":      valueOr,
		"valueOrEmpty": valueOrEmpty,
		"derefStr":     derefStr,
		"strMap":       strMap,
		"strMapSet":    strMapSet,
		"base64":       base64.StdEncoding.EncodeToString,
	}
	t := &Template{
		templates: template.Must(template.New("base").Funcs(funcMap).ParseFS(FS, "*")),
	}
	e.Renderer = t
}
