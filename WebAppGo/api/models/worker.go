package models

import (
	"WebAppGo/api/types"
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
	Id       string         `json:"id" db:"id"`
	Hostname string         `json:"hostname" db:"hostname"`
	Country  string         `json:"country" db:"country"`
	IpAddr   string         `json:"ip_addr" db:"ip_addr"`
	Os       string         `json:"os" db:"os"`
	IsAdmin  types.NullBool `json:"is_admin" db:"is_admin"`
	Boost    types.NullBool `json:"boost" db:"boost"`
	LastSeen time.Time      `json:"last_seen" db:"last_seen"`
}

func (w *Worker) HasEmptyFields() bool {
	return w.Id == "" || w.Hostname == "" || w.Country == "" || w.IpAddr == "" ||
		w.Os == "" || !w.IsAdmin.Valid || !w.Boost.Valid
}

func (w *Worker) Scan(r types.Row) error {
	return r.Scan(&w.Id, &w.Hostname, &w.Country, &w.IpAddr, &w.Os, &w.IsAdmin, &w.Boost, &w.LastSeen)
}
