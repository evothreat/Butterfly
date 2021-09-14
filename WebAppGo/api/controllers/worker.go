package controllers

import (
	"WebAppGo/api"
	"WebAppGo/api/models"
	"WebAppGo/utils"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// TODO: add default error handler!

func GetAllWorkers(c echo.Context) error {
	rows, err := db.Query("SELECT * FROM workers")
	if err != nil {
		return err
	}
	//defer rows.Close()
	workers := make([]*models.Worker, 0, 15)
	for rows.Next() {
		w := &models.Worker{}
		if err := w.Scan(rows); err != nil {
			return err
		}
		workers = append(workers, w)
	}
	return c.JSON(http.StatusOK, &workers)
}

func GetWorker(c echo.Context) error {
	worker := models.Worker{}
	if cols := c.QueryParam("props"); cols != "" && utils.IsValidListString(cols) { // TODO: change to "fields"
		row := db.QueryRow("SELECT "+cols+" FROM workers WHERE id=?", c.Param("wid"))
		data, err := worker.ScanColumns(row, cols) // use (&models.Worker).ScanColumns()?
		if err != nil {
			return c.NoContent(http.StatusUnprocessableEntity)
		}
		return c.JSON(http.StatusOK, &data)
	}
	row := db.QueryRow("SELECT * FROM workers WHERE id=?", c.Param("wid"))
	if err := worker.Scan(row); err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	return c.JSON(http.StatusOK, &worker)
}

func CreateWorker(c echo.Context) error {
	var w models.Worker
	if (&echo.DefaultBinder{}).BindBody(c, &w) != nil || w.HasEmptyFields() {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	_, err := db.Exec("INSERT INTO workers VALUES(?,?,?,?,?,?,?,?)",
		w.Id, w.Hostname, w.Country, w.IpAddr, w.Os, w.IsAdmin, w.Boost, time.Now())
	if err != nil {
		return err
	}
	os.Mkdir(filepath.Join(api.UPLOADS_DIR, w.Id), os.ModePerm)
	return c.NoContent(http.StatusCreated)
}

func DeleteWorker(c echo.Context) error {
	res, err := db.Exec("DELETE FROM workers WHERE id=?", c.Param("wid"))
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	os.Remove(filepath.Join(api.UPLOADS_DIR, c.Param("wid")))
	return c.NoContent(http.StatusOK)
}

func UpdateWorker(c echo.Context) error {
	valuesMap := map[string]interface{}{}
	if (&echo.DefaultBinder{}).BindBody(c, &valuesMap) != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	cols, vals := utils.ValuesMapToWhere(valuesMap)
	vals = append(vals, c.Param("wid"))
	res, err := db.Exec("UPDATE workers SET "+cols+" WHERE id=?", vals...) // very bad security...
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	return c.NoContent(http.StatusOK)
}
