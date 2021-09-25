package types

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type Template struct {
	Name string
	Tmpl *template.Template
	//Single	bool
}

type TemplateRegistry struct {
	Templates map[string]Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.Templates[name]
	if !ok {
		return fmt.Errorf("template not found: %s", name)
	}
	return tmpl.Tmpl.ExecuteTemplate(w, tmpl.Name, data)
}
