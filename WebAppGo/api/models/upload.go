package models

import "time"

const UploadSchema = `CREATE TABLE IF NOT EXISTS uploads
  (
     id        INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
     filename  VARCHAR(65),
     type      VARCHAR(16),
     size      BIGINT,
     created   DATETIME,
     worker_id VARCHAR(22),
	 FOREIGN KEY(worker_id) REFERENCES workers(id)
  );`

type Upload struct {
	Id       int       `json:"id" db:"id"`
	Filename string    `json:"filename" db:"filename"`
	Type     string    `json:"type" db:"type"`
	Size     uint64    `json:"size" db:"size"`
	Created  time.Time `json:"created" db:"created"`
	WorkerId string    `json:"worker_id" db:"worker_id"`
}
