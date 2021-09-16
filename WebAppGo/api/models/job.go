package models

import (
	"WebAppGo/api/types"
	"time"
)

const JobSchema = `CREATE TABLE IF NOT EXISTS jobs
  (
     id        INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
     todo      VARCHAR(250),
     is_done   BOOLEAN,
	 created   DATETIME,
	 worker_id VARCHAR(22),
	 FOREIGN KEY(worker_id) REFERENCES workers(id)
  );`

type Job struct {
	Id       int            `json:"id" db:"id"`
	Todo     string         `json:"todo" db:"todo"`
	IsDone   types.NullBool `json:"is_done" db:"is_done"`
	Created  time.Time      `json:"created" db:"created"`
	WorkerId string         `json:"worker_id" db:"worker_id"`
}

func (j *Job) HasEmptyFields() bool {
	return j.Todo == "" || !j.IsDone.Valid || j.WorkerId == ""
}

func (j *Job) Scan(r types.Row) error {
	return r.Scan(&j.Id, &j.Todo, &j.IsDone, &j.Created, &j.WorkerId)
}
