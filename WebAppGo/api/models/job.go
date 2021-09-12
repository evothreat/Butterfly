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
     worker_id VARCHAR(22),
	 created   DATETIME,
	 FOREIGN KEY(worker_id) REFERENCES workers(id)
  );`

type Job struct {
	Id       int            `json:"id" db:"id"`
	Todo     string         `json:"todo" db:"todo"`
	IsDone   types.NullBool `json:"is_done" db:"is_done"`
	WorkerId string         `json:"worker_id" db:"worker_id"`
	Created  time.Time      `json:"created" db:"created"`
}

func (j *Job) HasEmptyFields() bool {
	return j.Id == 0 || j.Todo == "" || !j.IsDone.Valid || j.WorkerId == ""
}

func (j *Job) Scan(r types.Row) error {
	return r.Scan(&j.Id, &j.Todo, &j.IsDone, &j.Created, &j.WorkerId)
}
