package models

import (
	"WebAppGo/api/types"
	"database/sql"
	"errors"
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
	Id       string         `json:"id" db:"id"`
	Hostname string         `json:"hostname" db:"hostname"`
	Country  string         `json:"country" db:"country"`
	IpAddr   string         `json:"ip_addr" db:"ip_addr"`
	Os       string         `json:"os" db:"os"`
	IsAdmin  types.NullBool `json:"is_admin" db:"is_admin"`
	Boost    types.NullBool `json:"boost" db:"boost"`
	LastSeen time.Time      `json:"last_seen" db:"last_seen"`
}

func (w *Worker) hasEmptyFields() bool {
	return w.Id == "" || w.Hostname == "" || w.Country == "" || w.IpAddr == "" ||
		w.Os == "" || !w.IsAdmin.Valid || !w.Boost.Valid
}

func (w *Worker) Save() error {
	if w.hasEmptyFields() {
		return errors.New("not all fields are set")
	}
	_, err := db.Exec("INSERT INTO workers VALUES(?,?,?,?,?,?,?,?)",
		w.Id, w.Hostname, w.Country, w.IpAddr, w.Os, w.IsAdmin, w.Boost, time.Now())
	return err
}

type WorkerWhereStmt types.WhereStmt

func scanWorkers(rows *sql.Rows) ([]*Worker, error) {
	workers := make([]*Worker, 0, 15)
	for rows.Next() {
		w := &Worker{}
		if err := rows.Scan(&w.Id, &w.Hostname, &w.Country, &w.IpAddr, &w.Os, &w.IsAdmin, &w.Boost, &w.LastSeen); err != nil {
			return nil, err
		}
		workers = append(workers, w)
	}
	return workers, nil
}

func FilterWorkers(cols string, values ...interface{}) *WorkerWhereStmt {
	return &WorkerWhereStmt{
		Cols:   cols,
		Values: values,
	}
}

func (wws *WorkerWhereStmt) Get() ([]*Worker, error) {
	rows, err := db.Query("SELECT * FROM workers WHERE "+wws.Cols, wws.Values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWorkers(rows)
}

func (wws *WorkerWhereStmt) GetFirst() (*Worker, error) {
	w := &Worker{}
	row := db.QueryRow("SELECT * FROM workers WHERE "+wws.Cols+" LIMIT 1", wws.Values...)
	if err := row.Scan(&w.Id, &w.Hostname, &w.Country, &w.IpAddr, &w.Os, &w.IsAdmin, &w.Boost, &w.LastSeen); err != nil {
		return nil, err
	}
	return w, nil
}

func (wws *WorkerWhereStmt) Delete() (int64, error) {
	res, err := db.Exec("DELETE FROM workers WHERE "+wws.Cols, wws.Values...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (wws *WorkerWhereStmt) Update(w *Worker) (int64, error) {
	stmt := "UPDATE workers SET "
	values := make([]interface{}, 0, 8+len(wws.Values))
	if w.Hostname != "" {
		stmt += "hostname=?,"
		values = append(values, w.Hostname)
	}
	if w.Country != "" {
		stmt += "country=?,"
		values = append(values, w.Country)
	}
	if w.IpAddr != "" {
		stmt += "ip_addr=?,"
		values = append(values, w.IpAddr)
	}
	if w.Os != "" {
		stmt += "os=?,"
		values = append(values, w.Os)
	}
	if w.IsAdmin.Valid {
		stmt += "is_admin=?,"
		values = append(values, w.IsAdmin.Bool) // pass NullBool struct?
	}
	if w.Boost.Valid {
		stmt += "boost=?"
		values = append(values, w.Boost.Bool)
	}
	if !w.LastSeen.IsZero() {
		stmt += "last_seen=?"
		values = append(values, w.LastSeen)
	}
	stmt = strings.TrimSuffix(stmt, ",") + " WHERE " + wws.Cols
	values = append(values, wws.Values...)
	res, err := db.Exec(stmt, values...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func GetAllWorkers() ([]*Worker, error) {
	rows, err := db.Query("SELECT * FROM workers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWorkers(rows)
}
