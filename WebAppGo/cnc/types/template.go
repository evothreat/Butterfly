package types

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type TemplateRegistry struct {
	Templates  map[string]*template.Template
	EntryPoint string
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.Templates[name]
	if !ok {
		return fmt.Errorf("template not found: %s", name)
	}
	return tmpl.ExecuteTemplate(w, t.EntryPoint, data)
}
