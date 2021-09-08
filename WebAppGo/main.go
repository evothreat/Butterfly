package main

import (
	"WebAppGo/api/controllers"
	"WebAppGo/api/models"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := models.SetupDatabase("root:root@tcp(localhost:3306)/data?parseTime=true"); err != nil {
		panic(err)
	}
	models.AddTestData()
	e := echo.New()
	e.Debug = true

	e.GET("/api/workers", controllers.GetWorkers)
	e.POST("/api/workers", controllers.CreateWorker)
	e.GET("/api/workers/:wid", controllers.GetWorker)
	e.DELETE("/api/workers/:wid", controllers.DeleteWorker)

	e.Start("localhost:8080")
}
