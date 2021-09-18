package controllers

import (
	"WebAppGo/api"
	"WebAppGo/api/models"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
)

// TODO: handle errors

func CreateReport(c echo.Context) error {
	jobId, _ := strconv.Atoi(c.Param("jid"))
	report, err := io.ReadAll(io.LimitReader(c.Request().Body, api.MAX_REPORT_LEN))
	if err != nil {
		return err
	}
	if len(report) == 0 {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	_, err = db.Exec("INSERT INTO job_reports(job_id, report) VALUES(?,?)", jobId, report)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE jobs SET is_done=1 WHERE id=?", jobId)	// TODO: check if job already done?
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func GetReport(c echo.Context) error {
	jobId, _ := strconv.Atoi(c.Param("jid"))
	var jr models.JobReport
	row := db.QueryRow("SELECT * FROM job_reports WHERE job_id=?", jobId)
	if err := jr.Scan(row); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &jr)
}

func DeleteReport(c echo.Context) error {
	jobId, _ := strconv.Atoi(c.Param("jid"))
	if _, err := db.Exec("DELETE FROM jobs WHERE id=?", jobId); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}