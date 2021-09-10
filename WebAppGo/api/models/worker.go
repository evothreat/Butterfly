package models

import (
	"strings"
	"time"
)

// TODO: implement Scan() and Value() Methods for time!

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
	IsAdmin  NullBool  `json:"is_admin" db:"is_admin"`
	Boost    NullBool  `json:"boost" db:"boost"`
	LastSeen time.Time `json:"last_seen" db:"last_seen"`
}

func (w *Worker) Save() error {
	const stmt = "INSERT INTO workers(id,hostname,country,ip_addr,os,is_admin,boost,last_seen) VALUES(?,?,?,?,?,?,?,?)"
	_, err := db.Exec(stmt, w.Id, w.Hostname, w.Country, w.IpAddr, w.Os, w.IsAdmin, w.Boost, time.Now())
	return err
}

func (w *Worker) Update() (int64, error) {
	stmt := "UPDATE workers SET "
	fields := make([]interface{}, 0, 7)
	if w.Hostname != "" {
		stmt += "hostname=?,"
		fields = append(fields, w.Hostname)
	}
	if w.Country != "" {
		stmt += "country=?,"
		fields = append(fields, w.Country)
	}
	if w.IpAddr != "" {
		stmt += "ip_addr=?,"
		fields = append(fields, w.IpAddr)
	}
	if w.Os != "" {
		stmt += "os=?,"
		fields = append(fields, w.Os)
	}
	if w.IsAdmin.Valid {
		stmt += "is_admin=?,"
		fields = append(fields, w.IsAdmin.Bool) // pass NullBool struct?
	}
	if w.Boost.Valid {
		stmt += "boost=?"
		fields = append(fields, w.Boost.Bool)
	}
	stmt = strings.TrimSuffix(stmt, ",") + " WHERE id=?"
	fields = append(fields, w.Id) // check whether id is set?
	res, err := db.Exec(stmt, fields...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func GetWorker(id string) (*Worker, error) {
	const stmt = "SELECT hostname,country,ip_addr,os,is_admin,boost,last_seen FROM workers WHERE id=?"
	w := &Worker{Id: id}
	err := db.QueryRow(stmt, id).Scan(&w.Hostname, &w.Country, &w.IpAddr, &w.Os, &w.IsAdmin, &w.Boost, &w.LastSeen)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func GetAllWorkers() ([]*Worker, error) {
	rows, err := db.Query("SELECT * FROM workers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	workers := make([]*Worker, 0, 15)
	for rows.Next() {
		w := &Worker{} // TODO: use pointer?
		err = rows.Scan(&w.Id, &w.Hostname, &w.Country, &w.IpAddr, &w.Os, &w.IsAdmin, &w.Boost, &w.LastSeen)
		if err != nil {
			return nil, err
		}
		workers = append(workers, w)
	}
	return workers, nil
}

func DeleteWorker(id string) (int64, error) {
	res, err := db.Exec("DELETE FROM workers WHERE id=?", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
