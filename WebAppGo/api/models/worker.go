package models

import (
	"WebAppGo/api/types"
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

func (w *Worker) HasEmptyFields() bool {
	return w.Id == "" || w.Hostname == "" || w.Country == "" || w.IpAddr == "" ||
		w.Os == "" || !w.IsAdmin.Valid || !w.Boost.Valid
}

func (w *Worker) Scan(r types.Row) error {
	return r.Scan(&w.Id, &w.Hostname, &w.Country, &w.IpAddr, &w.Os, &w.IsAdmin, &w.Boost, &w.LastSeen)
}

func (w *Worker) ScanColumns(r types.Row, colsStr string) (map[string]interface{}, error) {
	valuesMap := make(map[string]interface{})
	for _, c := range strings.Split(colsStr, ",") {
		valuesMap[c] = nil
	}
	values := make([]interface{}, len(valuesMap))
	i := 0
	if _, ok := valuesMap["id"]; ok {
		values[i] = &w.Id // is valuesMap["id"] = values[i] faster?
		valuesMap["id"] = &w.Id
		i++
	}
	if _, ok := valuesMap["hostname"]; ok {
		values[i] = &w.Hostname // try valuesMap["id"] = values[i]?
		valuesMap["hostname"] = &w.Hostname
		i++
	}
	if _, ok := valuesMap["country"]; ok {
		values[i] = &w.Country
		valuesMap["country"] = &w.Country
		i++
	}
	if _, ok := valuesMap["ip_addr"]; ok {
		values[i] = &w.IpAddr
		valuesMap["ip_addr"] = &w.IpAddr
		i++
	}
	if _, ok := valuesMap["os"]; ok {
		values[i] = &w.Os
		valuesMap["os"] = &w.Os
		i++
	}
	if _, ok := valuesMap["is_admin"]; ok {
		values[i] = &w.IsAdmin
		valuesMap["is_admin"] = &w.IsAdmin
		i++
	}
	if _, ok := valuesMap["boost"]; ok {
		values[i] = &w.Boost
		valuesMap["boost"] = &w.Boost
	}
	if _, ok := valuesMap["last_seen"]; ok {
		values[i] = &w.LastSeen
		valuesMap["last_seen"] = &w.LastSeen
	}
	return valuesMap, r.Scan(values...)
}

func (w *Worker) AsStmt() (string, []interface{}) {
	stmt := ""
	values := make([]interface{}, 0, 8)
	if w.Id != "" {
		stmt += "id=?,"
		values = append(values, &w.Id)
	}
	if w.Hostname != "" {
		stmt += "hostname=?,"
		values = append(values, &w.Hostname)
	}
	if w.Country != "" {
		stmt += "country=?,"
		values = append(values, &w.Country)
	}
	if w.IpAddr != "" {
		stmt += "ip_addr=?,"
		values = append(values, &w.IpAddr)
	}
	if w.Os != "" {
		stmt += "os=?,"
		values = append(values, &w.Os)
	}
	if w.IsAdmin.Valid {
		stmt += "is_admin=?,"
		values = append(values, &w.IsAdmin.Bool) // pass NullBool struct?
	}
	if w.Boost.Valid {
		stmt += "boost=?"
		values = append(values, &w.Boost.Bool)
	}
	return strings.TrimSuffix(stmt, ","), values
}
