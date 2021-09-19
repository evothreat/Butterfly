package controllers

import (
	"WebAppGo/api"
	"database/sql"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
)

func CreateReport(c echo.Context) error {
	jobId, _ := strconv.Atoi(c.Param("jid"))
	if !rowExists("SELECT id FROM jobs WHERE worker_id=? AND id=? AND is_done=0", c.Param("wid"), jobId) {
		return c.NoContent(http.StatusNotFound)
	}
	report, err := io.ReadAll(io.LimitReader(c.Request().Body, api.MAX_REPORT_LEN))
	if err != nil {
		return err
	}
	if len(report) == 0 {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO job_reports(job_id, report) VALUES(?,?)", jobId, report)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("UPDATE jobs SET is_done=1 WHERE worker_id=? AND id=?", c.Param("wid"), jobId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func GetReport(c echo.Context) error {
	jobId, _ := strconv.Atoi(c.Param("jid"))
	var report string
	stmt := "SELECT job_reports.report FROM job_reports INNER JOIN jobs ON job_reports.job_id=jobs.id WHERE jobs.worker_id=? AND job_reports.job_id=?"
	row := db.QueryRow(stmt, c.Param("wid"), jobId)
	if err := row.Scan(&report); err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	return c.String(http.StatusOK, report)
}

func DeleteReport(c echo.Context) error {
	jobId, _ := strconv.Atoi(c.Param("jid"))
	stmt := "DELETE job_reports FROM job_reports INNER JOIN jobs ON job_reports.job_id=jobs.id WHERE jobs.worker_id=? AND job_reports.job_id=?"
	if _, err := db.Exec(stmt, c.Param("wid"), jobId); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
