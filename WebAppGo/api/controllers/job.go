package controllers

import (
	"WebAppGo/api"
	"WebAppGo/api/models"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func GetAllJobs(c echo.Context) error {
	workerId := c.Param("wid")
	stmt := "SELECT * FROM jobs WHERE worker_id=?"

	if _, ok := c.QueryParams()["undone"]; ok {
		stmt += " AND is_done=0" // or assign new statement?
		_, err := db.Exec("UPDATE workers SET last_seen=? WHERE id=?", time.Now(), workerId)
		if err != nil {
			return err
		}
	}
	rows, err := db.Query(stmt, workerId)
	if err != nil {
		return err
	}
	jobs := make([]*models.Job, 0, api.MIN_LIST_CAP)
	for rows.Next() {
		job := &models.Job{}
		if err := job.Scan(rows); err != nil {
			return err
		}
		jobs = append(jobs, job)
	}
	return c.JSON(http.StatusOK, &jobs)
}

func GetJob(c echo.Context) error {
	jobId, _ := strconv.Atoi(c.Param("jid"))
	var job models.Job
	row := db.QueryRow("SELECT * FROM jobs WHERE worker_id=? AND id=?", c.Param("wid"), jobId)
	if err := job.Scan(row); err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	return c.JSON(http.StatusOK, &job)
}

func CreateJob(c echo.Context) error {
	var job models.Job
	if (&echo.DefaultBinder{}).BindBody(c, &job) != nil || job.HasEmptyFields() {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	_, err := db.Exec("INSERT INTO jobs(todo,is_done,created,worker_id) VALUES(?,?,?,?)",
		job.Todo, job.IsDone, time.Now(), c.Param("wid")) // TODO: add is_done in js files!
	if err != nil {
		if isNoReferencedRowErr(err) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func DeleteJob(c echo.Context) error {
	jobId, _ := strconv.Atoi(c.Param("jid")) // default id value will be 0
	res, err := db.Exec("DELETE FROM jobs WHERE worker_id=? AND id=?", c.Param("wid"), jobId)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	return c.NoContent(http.StatusOK)
}
