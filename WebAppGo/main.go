package main

import (
	"WebAppGo/api"
	"WebAppGo/cnc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	api.SetupDatabase("root:root@tcp(localhost:3306)/data?parseTime=true")
	//api.AddTestData()

	e := echo.New()
	e.Debug = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	api.SetupRoutes(e)

	cnc.Setup(e)

	e.Start("localhost:8080")
}
