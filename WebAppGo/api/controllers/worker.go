package controllers

import (
	"WebAppGo/api/models"
	"github.com/labstack/echo/v4"
	"net/http"
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
		return err
	}
	return c.JSON(http.StatusOK, worker)
}

func CreateWorker(c echo.Context) error {
	var worker models.Worker
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &worker); err != nil {
		return err
	}
	if err := worker.Save(); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusCreated, []byte(`{"msg": "worker created"}`))
}

func DeleteWorker(c echo.Context) error {
	if _, err := models.DeleteWorker(c.Param("wid")); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, []byte(`{"msg": "worker deleted"}`))
}

func UpdateWorker(c echo.Context) error {
	var worker models.Worker
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &worker); err != nil {
		return err
	}
	worker.Id = c.Param("wid")
	if _, err := worker.Update(); err != nil {
		return err
	}
	return nil
}
