package api

import "time"

type Worker struct {
	Id       string    `json:"id" gorm:"primary key; type: varchar(22)"`
	Hostname string    `json:"hostname" gorm:"type: varchar(30)"`
	Country  string    `json:"country" gorm:"type: varchar(15)"`
	IpAddr   string    `json:"ip_addr" gorm:"type: varchar(15)"`
	Os       string    `json:"os" gorm:"type: varchar(15)"`
	IsAdmin  bool      `json:"is_admin" gorm:"type: boolean"`
	Boost    bool      `json:"boost" gorm:"type: boolean"`
	LastSeen time.Time `json:"last_seen" gorm:"type: datetime"`
}

type HardwareInfo struct {
	Gpu      string `json:"gpu" gorm:"type: varchar(50)"`
	Cpu      string `json:"cpu" gorm:"type: varchar(65)"`
	Ram      string `json:"ram" gorm:"type: varchar(10)"`
	WorkerId string `json:"worker_id" gorm:"primary key; type: varchar(22)"`
}

type Job struct {
	Id       string    `json:"id" gorm:"primary key; type: integer"`
	Todo     string    `json:"todo" gorm:"type: varchar(250)"`
	IsDone   bool      `json:"is_done" gorm:"type: boolean"`
	Created  time.Time `json:"created" gorm:"type: datetime"`
	WorkerId string    `json:"worker_id" gorm:"type: varchar(22)"`
}

type JobReport struct {
	JobId  int    `json:"job_id" gorm:"primary key; type: integer"`
	Report string `json:"report" gorm:"type: text"`
}

type Upload struct {
	Id       int       `json:"id" gorm:"primary key; type: integer"`
	Filename string    `json:"filename" gorm:"type: varchar(65)"`
	Type     string    `json:"type" gorm:"type: varchar(16)"`
	Size     uint64    `json:"size" gorm:"type: biginteger"`
	Created  time.Time `json:"created" gorm:"type: datetime"`
	WorkerId string    `json:"worker_id" gorm:"type: varchar(22)"`
}
