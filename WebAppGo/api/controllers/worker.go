package controllers

import (
	"WebAppGo/api"
	"WebAppGo/api/models"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"path/filepath"
)

// TODO: add default error handler!

func GetAllWorkers(c echo.Context) error {
	workers, err := models.GetAllWorkers()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, workers)
}

func GetWorker(c echo.Context) error {
	worker, err := models.FilterWorkers("id=?", c.Param("wid")).GetFirst()
	if err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	return c.JSON(http.StatusOK, worker)
}

// TODO: check if all fields are filled!
// TODO: check for error type?
// TODO: handle os.Mkdir error?

func CreateWorker(c echo.Context) error {
	var worker models.Worker
	if (&echo.DefaultBinder{}).BindBody(c, &worker) != nil || worker.Save() != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	os.Mkdir(filepath.Join(api.UPLOADS_DIR, worker.Id), os.ModePerm)
	return c.NoContent(http.StatusCreated)
}

func DeleteWorker(c echo.Context) error {
	workerId := c.Param("wid")
	n, err := models.FilterWorkers("id=?", workerId).Delete()
	if err != nil {
		return err
	}
	if n == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	os.Remove(filepath.Join(api.UPLOADS_DIR, workerId))
	return c.NoContent(http.StatusOK)
}

func UpdateWorker(c echo.Context) error {
	var worker models.Worker
	if (&echo.DefaultBinder{}).BindBody(c, &worker) != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	n, err := models.FilterWorkers("id=?", c.Param("wid")).Update(&worker)
	if err != nil {
		return err
	}
	if n == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	return c.NoContent(http.StatusOK)
}
