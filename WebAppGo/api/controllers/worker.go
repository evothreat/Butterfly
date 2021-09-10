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

// TODO: instead of returning error, return status code!

func GetAllWorkers(c echo.Context) error {
	workers, err := models.GetAllWorkers()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, workers)
}

func GetWorker(c echo.Context) error {
	worker, err := models.GetWorker(c.Param("wid"))
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return err
	}
	return c.JSON(http.StatusOK, worker)
}

func CreateWorker(c echo.Context) error {
	var worker models.Worker
	binder := &echo.DefaultBinder{}
	if binder.BindBody(c, &worker) != nil || worker.Save() != nil { // TODO: check for error type?
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}
	os.Mkdir(filepath.Join(api.UPLOADS_DIR, worker.Id), os.ModePerm) // TODO: handle error?
	return c.JSONBlob(http.StatusCreated, []byte(`{"msg": "worker created"}`))
}

func DeleteWorker(c echo.Context) error {
	workerId := c.Param("wid")
	n, err := models.DeleteWorker(workerId)
	if err != nil {
		return err
	}
	if n == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	os.Remove(filepath.Join(api.UPLOADS_DIR, workerId))
	return c.JSONBlob(http.StatusOK, []byte(`{"msg": "worker deleted"}`))
}

func UpdateWorker(c echo.Context) error {
	var worker models.Worker
	binder := &echo.DefaultBinder{}
	if binder.BindBody(c, &worker) != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}
	worker.Id = c.Param("wid")
	n, err := worker.Update()
	if err != nil {
		return err
	}
	if n == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return nil
}
