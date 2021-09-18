package controllers

import (
	"WebAppGo/api/models"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateHardwareInfo(c echo.Context) error {
	var hwi models.HardwareInfo
	if (&echo.DefaultBinder{}).BindBody(c, &hwi) != nil || hwi.HasEmptyFields() {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	_, err := db.Exec("INSERT INTO hardware_infos(gpu,cpu,ram,worker_id) VALUES(?,?,?,?)",
		hwi.Gpu, hwi.Cpu, hwi.Ram, c.Param("wid"))
	if err != nil {
		if IsDuplicateEntry(err) {
			return c.NoContent(http.StatusConflict)
		}
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func GetHardwareInfo(c echo.Context) error {
	var hwi models.HardwareInfo
	row := db.QueryRow("SELECT * FROM hardware_infos WHERE worker_id=?", c.Param("wid"))
	if err := hwi.Scan(row); err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	return c.JSON(http.StatusOK, &hwi)
}
