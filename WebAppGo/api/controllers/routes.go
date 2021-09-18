package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func SetupRoutes(g *echo.Group) {
	g.Use(ResourceExists)

	g.GET("/workers", GetAllWorkers)
	g.POST("/workers", CreateWorker)
	g.GET("/workers/:wid", GetWorker)
	g.DELETE("/workers/:wid", DeleteWorker)
	g.PATCH("/workers/:wid", UpdateWorker)

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

func ResourceExists(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		workerId := c.Param("wid")
		jobId := c.Param("jid")

		if jobId != "" {
			jobIdNum, err := strconv.Atoi(jobId)
			if err != nil || jobIdNum < 1 || !rowExists("SELECT id FROM jobs WHERE worker_id=? AND id=?", workerId, jobIdNum) {
				return c.NoContent(http.StatusNotFound)
			}
			return handler(c)
		}
		if workerId != "" {
			if !rowExists("SELECT id FROM workers WHERE id=?", workerId) {
				return c.NoContent(http.StatusNotFound)
			}
		}
		return handler(c)
	}
}