package cnc

import (
	"WebAppGo/cnc/types"
	"html/template"
)

func ParseTemplates(dir string) *types.TemplateRegistry {
	templates := map[string]types.Template{}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	templates["login"] = types.Template{
		Name: "login",
		Tmpl: template.Must(template.ParseFiles(dir + "login.html")),
	}
	baseHtml := dir + "base.html"
	templates["workers"] = types.Template{
		Name: "base",
		Tmpl: template.Must(template.ParseFiles(dir+"workers.html", baseHtml)),
	}
	templates["interact"] = types.Template{
		Name: "base",
		Tmpl: template.Must(template.ParseFiles(dir+"interact.html", dir+"uploads.html", dir+"history.html", baseHtml)),
	}
	return &types.TemplateRegistry{
		Templates: templates,
	}
}
