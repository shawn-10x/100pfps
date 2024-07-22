package views

import (
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func getTempFilesFromFolders(folders ...string) []string {
	var filepaths []string
	for _, folder := range folders {
		files, err := os.ReadDir(folder)
		if err != nil {
			panic(err.Error())
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".html") {
				filepaths = append(filepaths, folder+file.Name())
			}
		}
	}
	return filepaths
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func SetupViews(e *echo.Echo) {
	funcMap := template.FuncMap{
		"hasKey":       hasKey,
		"ternary":      ternary,
		"valueOr":      valueOr,
		"valueOrEmpty": valueOrEmpty,
		"derefStr":     derefStr,
		"stringMap":    stringMap,
	}
	t := &Template{
		templates: template.Must(template.New("base").Funcs(funcMap).ParseFiles(getTempFilesFromFolders("views/")...)),
	}
	e.Renderer = t
}
