package site

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

// Template class used for getting templates
type Template struct {
	Templates *template.Template
}

// Render function to implement echo Renderer interface
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
