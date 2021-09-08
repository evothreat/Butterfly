package models

const JobReportSchema = `CREATE TABLE IF NOT EXISTS job_reports
  (
     job_id INTEGER PRIMARY KEY,
     report TEXT,
     FOREIGN KEY(job_id) REFERENCES jobs(id)
  );`

type JobReport struct {
	JobId  int    `json:"job_id" db:"job_id"`
	Report string `json:"report" db:"report"`
}
