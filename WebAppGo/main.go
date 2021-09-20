package main

import (
	"WebAppGo/api"
	"WebAppGo/cnc"
	"WebAppGo/cnc/types"
	"github.com/labstack/echo/v4"
	"html/template"
)

func main() {
	api.SetupDatabase("root:root@tcp(localhost:3306)/data?parseTime=true")
	//api.AddTestData()

	e := echo.New()
	e.Debug = true

	apiGroup := e.Group("/api")
	api.SetupRoutes(apiGroup)

	cncGroup := e.Group("/cnc")
	cnc.SetupRoutes(cncGroup)

	e.Renderer = &types.Template{
		Templates: template.Must(template.ParseGlob("./cnc/templates/*.html")),
	}
	e.Static("/", "cnc/") // TODO: examine later...
	e.Start("localhost:8080")
}
