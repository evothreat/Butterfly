package controllers

import (
	"WebAppGo/api"
	"WebAppGo/api/models"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func GetAllJobs(c echo.Context) error {
	rows, err := db.Query("SELECT * FROM jobs WHERE worker_id=?", c.Param("wid"))
	if err != nil {
		return err
	}
	//defer rows.Close()
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

func GetUndoneJobs(c echo.Context) error {
	rows, err := db.Query("SELECT * FROM jobs WHERE worker_id=? AND is_done=0", c.Param("wid"))
	if err != nil {
		return err // TODO: first check if worker exist? Also select worker_id??
	}
	//defer rows.Close()
	jobs := make([]*models.Job, 0, api.MIN_LIST_CAP)
	for rows.Next() {
		job := &models.Job{}
		if err := job.Scan(rows); err != nil {
			return err
		}
		jobs = append(jobs, job)
	}
	db.Exec("UPDATE workers SET last_seen=? WHERE id=?", time.Now(), c.Param("wid")) // TODO: check for error? No rows?
	return c.JSON(http.StatusOK, &jobs)
}

func GetJob(c echo.Context) error {
	var job models.Job
	row := db.QueryRow("SELECT * FROM jobs WHERE worker_id=? AND id=?", c.Param("wid"), c.Param("jid"))
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
		job.Todo, job.IsDone, time.Now(), job.WorkerId) // TODO: add is_done in js files!
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func DeleteJob(c echo.Context) error {
	res, err := db.Exec("DELETE FROM jobs WHERE worker_id=? AND id=?", c.Param("wid"), c.Param("jid"))
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	return c.NoContent(http.StatusOK)
}
