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
	Id       string         `json:"id" db:"id"`
	Todo     string         `json:"todo" db:"todo"`
	IsDone   types.NullBool `json:"is_done" db:"is_done"`
	Created  time.Time      `json:"created" db:"created"`
	WorkerId string         `json:"worker_id" db:"worker_id"`
}
