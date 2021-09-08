package controllers

import (
	"WebAppGo/api/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetWorkers(c echo.Context) error {
	workers, err := models.ListWorkers()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &workers)
}

func GetWorker(c echo.Context) error {
	var worker models.Worker
	if err := worker.Load(c.Param("wid")); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &worker)
}

func CreateWorker(c echo.Context) error {
	var worker models.Worker
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &worker); err != nil {
		// TODO: instead of returning error, return error code
		return err
	}
	if err := worker.Save(); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusCreated, []byte(`{"msg": "worker created"}`))
}

func DeleteWorker(c echo.Context) error {
	if err := models.DeleteWorker(c.Param("wid")); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, []byte(`{"msg": "worker deleted"}`))
}

func UpdateWorker(c echo.Context) error {
	return nil
}
