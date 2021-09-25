package main

import (
	"WebAppGo/api"
	"WebAppGo/cnc"
	"encoding/gob"
	"github.com/labstack/echo/v4"
	"time"
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
	gob.Register(time.Time{})
	e.Renderer = cnc.ParseTemplates("resources/templates")
	e.Static("/static", "resources/static")
	e.Start("localhost:8080")
}
