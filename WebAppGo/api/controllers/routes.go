package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func SetupRoutes(g *echo.Group) {
	g.GET("/workers", GetAllWorkers)
	g.POST("/workers", CreateWorker)
	g.GET("/workers/:wid", GetWorker)
	g.DELETE("/workers/:wid", DeleteWorker)
	g.PATCH("/workers/:wid", UpdateWorker)

	g.Use(WorkerExists)

	g.GET("/workers/:wid/jobs", GetAllJobs)
	g.POST("/workers/:wid/jobs", CreateJob)
	g.GET("/workers/:wid/jobs/undone", GetUndoneJobs)
	g.GET("/workers/:wid/jobs/:jid", GetJob)
	g.DELETE("/workers/:wid/jobs/:jid", DeleteJob)

	g.POST("/workers/:wid/hardware", CreateHardwareInfo)
	g.GET("/workers/:wid/hardware", GetHardwareInfo)

	g.POST("/workers/:wid/uploads", CreateUpload)
	g.GET("/workers/:wid/uploads/:uid", GetUpload)
	g.DELETE("/workers/:wid/uploads/:uid", DeleteUpload)
	g.GET("/workers/:wid/uploads/:uid/info", GetUploadInfo)

	g.POST("/workers/:wid/jobs/:jid/report", CreateReport)
}

func WorkerExists(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !rowExists("SELECT id FROM workers WHERE id=?", c.Param("wid")) {
			return c.NoContent(http.StatusNotFound)
		}
		return handler(c)
	}
}
