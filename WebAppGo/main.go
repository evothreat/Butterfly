package main

import (
	"WebAppGo/api/controllers"
	"github.com/labstack/echo/v4"
)

func main() {
	controllers.SetupDatabase("root:root@tcp(localhost:3306)/data?parseTime=true")
	//controllers.AddTestData()
	e := echo.New()
	e.Debug = true

	controllers.SetupRoutes(e)

	e.Start("localhost:8080")
}
