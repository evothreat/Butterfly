package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func SetupRoutes(g *echo.Group) {
	g.GET("/workers", GetAllWorkers)        // OK
	g.POST("/workers", CreateWorker)        // OK
	g.GET("/workers/:wid", GetWorker)       // OK
	g.DELETE("/workers/:wid", DeleteWorker) // OK
	g.PATCH("/workers/:wid", UpdateWorker)  // OK

	g.Use(WorkerExists)

	g.GET("/workers/:wid/jobs", GetAllJobs)
	g.POST("/workers/:wid/jobs", CreateJob)        // OK
	g.GET("/workers/:wid/jobs/:jid", GetJob)       // OK
	g.DELETE("/workers/:wid/jobs/:jid", DeleteJob) // OK

	g.POST("/workers/:wid/hardware", CreateHardwareInfo) // OK
	g.GET("/workers/:wid/hardware", GetHardwareInfo)     // OK

	g.POST("/workers/:wid/uploads", CreateUpload)           // OK
	g.GET("/workers/:wid/uploads/:uid", GetUpload)          // OK
	g.DELETE("/workers/:wid/uploads/:uid", DeleteUpload)    // OK
	g.GET("/workers/:wid/uploads/:uid/info", GetUploadInfo) // OK

	g.POST("/workers/:wid/jobs/:jid/report", CreateReport) // USE SQL JOINS!
	g.GET("/workers/:wid/jobs/:jid/report", GetReport)
	g.DELETE("/workers/:wid/jobs/:jid/report", DeleteReport)
}

func WorkerExists(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !rowExists("SELECT id FROM workers WHERE id=?", c.Param("wid")) {
			return c.NoContent(http.StatusNotFound)
		}
		return handler(c)
	}
}
