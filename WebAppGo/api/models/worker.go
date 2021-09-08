package models

import "time"

const WorkerSchema = `CREATE TABLE IF NOT EXISTS workers
  (
     id        VARCHAR(22) NOT NULL PRIMARY KEY,
     hostname  VARCHAR(30),
     country   VARCHAR(15),
     ip_addr   VARCHAR(15),
     os        VARCHAR(15),
     is_admin  BOOLEAN,
     boost     BOOLEAN,
     last_seen DATETIME
  );`

type Worker struct {
	Id       string    `json:"id" db:"id"`
	Hostname string    `json:"hostname" db:"hostname"`
	Country  string    `json:"country" db:"country"`
	IpAddr   string    `json:"ip_addr" db:"ip_addr"`
	Os       string    `json:"os" db:"os"`
	IsAdmin  bool      `json:"is_admin" db:"is_admin"`
	Boost    bool      `json:"boost" db:"boost"`
	LastSeen time.Time `json:"last_seen" db:"last_seen"`
}

func (w *Worker) Save() error {
	const insertStmt = "INSERT INTO workers(id,hostname,country,ip_addr,os,is_admin,boost,last_seen) VALUES(?,?,?,?,?,?,?,?)"
	_, err := db.Exec(insertStmt, w.Id, w.Hostname, w.Country, w.IpAddr, w.Os, w.IsAdmin, w.Boost, time.Now())
	return err
}

func (w *Worker) Load(id string) error {
	const selectStmt = "SELECT id,hostname,country,ip_addr,os,is_admin,boost,last_seen FROM workers WHERE id=?"
	row := db.QueryRow(selectStmt, id)
	return row.Scan(&w.Id, &w.Hostname, &w.Country, &w.IpAddr, &w.Os, &w.IsAdmin, &w.Boost, &w.LastSeen)
}

func ListWorkers() ([]Worker, error) {
	rows, err := db.Query("SELECT * FROM workers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	workers := make([]Worker, 0, 10)
	for rows.Next() {
		w := Worker{} // TODO: use pointer?
		if err := rows.Scan(&w.Id, &w.Hostname, &w.Country, &w.IpAddr, &w.Os, &w.IsAdmin, &w.Boost, &w.LastSeen); err != nil {
			return nil, err
		}
		workers = append(workers, w)
	}
	return workers, nil
}

func DeleteWorker(id string) error {
	if _, err := db.Exec("DELETE FROM workers WHERE id=?", id); err != nil {
		return err
	}
	return nil
}
