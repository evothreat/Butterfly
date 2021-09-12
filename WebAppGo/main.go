package main

import (
	"WebAppGo/api/controllers"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := controllers.SetupDatabase("root:root@tcp(localhost:3306)/data?parseTime=true"); err != nil {
		panic(err)
	}
	//models.AddTestData()
	e := echo.New()
	e.Debug = true

	e.GET("/api/workers", controllers.GetAllWorkers)
	e.POST("/api/workers", controllers.CreateWorker)
	e.GET("/api/workers/:wid", controllers.GetWorker)
	e.DELETE("/api/workers/:wid", controllers.DeleteWorker)
	e.PATCH("/api/workers/:wid", controllers.UpdateWorker)

	e.Start("localhost:8080")
}
