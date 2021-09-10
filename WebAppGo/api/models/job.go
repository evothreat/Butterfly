package models

import (
	"WebAppGo/api/types"
	"database/sql"
	"errors"
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

func (j *Job) hasEmptyFields() bool {
	return j.Id == 0 || j.Todo == "" || !j.IsDone.Valid || j.WorkerId == ""
}

func (j *Job) Save() error {
	if j.hasEmptyFields() {
		return errors.New("not all fields are set")
	}
	const stmt = "INSERT INTO jobs(todo,is_done,worker_id,created) VALUES(?,?,?,?)"
	_, err := db.Exec(stmt, j.Todo, j.IsDone, j.WorkerId, time.Now()) // TODO: dont pass is_done? set default?
	return err
}

func scanJobs(rows *sql.Rows) ([]*Job, error) {
	jobs := make([]*Job, 0, 15)
	for rows.Next() {
		j := &Job{}
		if err := rows.Scan(&j.Id, &j.Todo, &j.IsDone, &j.Created, &j.WorkerId); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func FilterJobs(cols string, values ...interface{}) ([]*Job, error) {
	rows, err := db.Query("SELECT * FROM jobs WHERE "+cols, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJobs(rows)
}

func GetAllJobs() ([]*Job, error) {
	rows, err := db.Query("SELECT * FROM jobs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJobs(rows)
}
