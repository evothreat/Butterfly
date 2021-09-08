package api

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

const workersSchema = `CREATE TABLE IF NOT EXISTS workers
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

const hardwareInfosSchema = `CREATE TABLE IF NOT EXISTS hardware_infos
  (
     gpu       VARCHAR(50),
     cpu       VARCHAR(65),
     ram       VARCHAR(10),
     worker_id VARCHAR(22) PRIMARY KEY,
	 FOREIGN KEY(worker_id) REFERENCES workers(id)
  );`

const jobsSchema = `CREATE TABLE IF NOT EXISTS jobs
  (
     id        INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
     todo      VARCHAR(250),
     is_done   BOOLEAN,
     created   DATETIME,
     worker_id VARCHAR(22),
	 FOREIGN KEY(worker_id) REFERENCES workers(id)
  );`

const jobReportsSchema = `CREATE TABLE IF NOT EXISTS job_reports
  (
     job_id INTEGER PRIMARY KEY,
     report TEXT,
     FOREIGN KEY(job_id) REFERENCES jobs(id)
  );`

const uploadsSchema = `CREATE TABLE IF NOT EXISTS uploads
  (
     id        INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
     filename  VARCHAR(65),
     type      VARCHAR(16),
     size      BIGINT,
     created   DATETIME,
     worker_id VARCHAR(22),
	 FOREIGN KEY(worker_id) REFERENCES workers(id)
  );`

var db *sqlx.DB

func SetupDatabase(dbPath string) {
	db = sqlx.MustConnect("mysql", dbPath)
	db.MustExec(workersSchema)
	db.MustExec(hardwareInfosSchema)
	db.MustExec(jobsSchema)
	db.MustExec(jobReportsSchema)
	db.MustExec(uploadsSchema)
}

func AddTestData() {
	w1 := Worker{
		Id:       "C1vHa4fB9kukvA6ILps0kQ",
		Hostname: "Predator",
		Os:       "Windows 10",
		Country:  "Germany",
		IpAddr:   "234.145.222.1",
		IsAdmin:  true,
		Boost:    true,
		LastSeen: time.Now(),
	}
	w2 := Worker{
		Id:       "K8_RNoHLL0S-UELe3WqhSw",
		Hostname: "Helios 300",
		Os:       "Windows 7",
		Country:  "England",
		IpAddr:   "127.2.4.1",
		IsAdmin:  false,
		Boost:    false,
		LastSeen: time.Now(),
	}
	w3 := Worker{
		Id:       "8GYEaE8G5E2oZtVyxY8nxg",
		Hostname: "Acer Nexus",
		Os:       "Windows 8",
		Country:  "USA",
		IpAddr:   "127.11.55.1",
		IsAdmin:  true,
		Boost:    false,
		LastSeen: time.Now(),
	}
	w4 := Worker{
		Id:       "cxi5YsNdk020NNMYhqZ78g",
		Hostname: "KNIGHT",
		Os:       "Windows 10",
		Country:  "Russia",
		IpAddr:   "127.66.33.1",
		IsAdmin:  false,
		Boost:    false,
		LastSeen: time.Now(),
	}

	/*const insertQuery1 = "INSERT INTO workers VALUES(:id, :hostname, :country, :ip_addr, :os, :is_admin, :boost, :last_seen)"

	db.NamedExec(insertQuery1, w1)
	db.NamedExec(insertQuery1, w2)
	db.NamedExec(insertQuery1, w3)
	db.NamedExec(insertQuery1, w4)*/

	const insertQuery1 = "INSERT INTO workers(id,hostname,country,ip_addr,os,is_admin,boost,last_seen) VALUES(?,?,?,?,?,?,?,?)"

	db.Exec(insertQuery1, w1.Id, w1.Hostname, w1.Country, w1.IpAddr, w1.Os, w1.IsAdmin, w1.Boost, w1.LastSeen)
	db.Exec(insertQuery1, w2.Id, w2.Hostname, w2.Country, w2.IpAddr, w2.Os, w2.IsAdmin, w2.Boost, w2.LastSeen)
	db.Exec(insertQuery1, w3.Id, w3.Hostname, w3.Country, w3.IpAddr, w3.Os, w3.IsAdmin, w3.Boost, w3.LastSeen)
	db.Exec(insertQuery1, w4.Id, w4.Hostname, w4.Country, w4.IpAddr, w4.Os, w4.IsAdmin, w4.Boost, w4.LastSeen)

	hw1 := HardwareInfo{
		Gpu:      "NVIDIA GeForce GTX 1060",
		Cpu:      "Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz 2.21 GHz",
		Ram:      "8,0 GB",
		WorkerId: w1.Id,
	}
	hw2 := HardwareInfo{
		Gpu:      "AMD Radeon RX 5600 XT",
		Cpu:      "Intel Core i9-10900K Processor",
		Ram:      "4,0 GB",
		WorkerId: w2.Id,
	}
	hw3 := HardwareInfo{
		Gpu:      "Nvidia GeForce RTX 3080",
		Cpu:      "AMD Ryzen 9 5900X",
		Ram:      "12,0 GB",
		WorkerId: w3.Id,
	}
	hw4 := HardwareInfo{
		Gpu:      "Intel Core i9-10980XE Extreme Edition Processor",
		Cpu:      "Intel(R) UHD Graphics 630",
		Ram:      "32,0 GB",
		WorkerId: w4.Id,
	}
	/*const insertQuery2 = "INSERT INTO hardware_infos VALUES(:gpu, :cpu, :ram, :worker_id)"

	db.NamedExec(insertQuery2, hw1)
	db.NamedExec(insertQuery2, hw2)
	db.NamedExec(insertQuery2, hw3)
	db.NamedExec(insertQuery2, hw4)*/

	const insertQuery2 = "INSERT INTO hardware_infos(gpu,cpu,ram,worker_id) VALUES(?,?,?,?)"

	db.Exec(insertQuery2, hw1.Gpu, hw1.Cpu, hw1.Ram, hw1.WorkerId)
	db.Exec(insertQuery2, hw2.Gpu, hw2.Cpu, hw2.Ram, hw2.WorkerId)
	db.Exec(insertQuery2, hw3.Gpu, hw3.Cpu, hw3.Ram, hw3.WorkerId)
	db.Exec(insertQuery2, hw4.Gpu, hw4.Cpu, hw4.Ram, hw4.WorkerId)

	/*hw5 := HardwareInfo{
		Gpu: "Intel Core i9-10980XE Extreme Edition Processor",
		Cpu: "Intel(R) UHD Graphics 630",
		Ram: "32,0 GB",
		WorkerId: "wfwofhwfwhfiuwh",
	}
	if _, err := db.NamedExec(insertQuery2, hw5); err != nil {
		fmt.Println(err)
	}*/

	const insertQuery3 = "INSERT INTO jobs(todo,is_done,created,worker_id) VALUES(?,?,?,?)"

	j1 := Job{
		Todo:     "ddos fbi.gov",
		IsDone:   false,
		Created:  time.Now(),
		WorkerId: w1.Id,
	}

	j2 := Job{
		Todo:     "upload passwords.txt",
		IsDone:   true,
		Created:  time.Now(),
		WorkerId: w1.Id,
	}

	db.Exec(insertQuery3, j1.Todo, j1.IsDone, j1.Created, j1.WorkerId)
	db.Exec(insertQuery3, j2.Todo, j2.IsDone, j2.Created, j2.WorkerId)
}
