package models

import "WebAppGo/api/types"

const JobReportSchema = `CREATE TABLE IF NOT EXISTS job_reports
  (
     job_id INTEGER PRIMARY KEY,
     report TEXT,
     FOREIGN KEY(job_id) REFERENCES jobs(id) ON DELETE CASCADE
  );`

type JobReport struct {
	JobId  int    `json:"job_id" db:"job_id"`
	Report string `json:"report" db:"report"`
}

func (jr *JobReport) HasEmptyFields() bool {
	return jr.JobId < 1 || jr.Report == ""
}

func (jr *JobReport) Scan(r types.Row) error {
	return r.Scan(&jr.JobId, &jr.Report)
}
