package cnc

import (
	"WebAppGo/cnc/types"
	"html/template"
	"path"
)

func parseTemplates() *types.TemplateRegistry {
	tmplDir := path.Join(RESOURCE_DIR, "/templates") + "/"
	templates := map[string]*template.Template{}
	templates["login"] = template.Must(template.ParseFiles(tmplDir + "login.html"))
	baseHtml := tmplDir + "base.html"
	templates["workers"] = template.Must(template.ParseFiles(tmplDir+"workers.html", baseHtml))
	templates["interact"] = template.Must(template.ParseFiles(tmplDir+"interact.html", tmplDir+"uploads.html", tmplDir+"history.html", baseHtml))
	return &types.TemplateRegistry{
		Templates:  templates,
		EntryPoint: "main",
	}
}
