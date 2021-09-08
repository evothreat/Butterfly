package api

import "time"

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

type HardwareInfo struct {
	Gpu      string `json:"gpu" db:"gpu"`
	Cpu      string `json:"cpu" db:"cpu"`
	Ram      string `json:"ram" db:"ram"`
	WorkerId string `json:"worker_id" db:"worker_id"`
}

type Job struct {
	Id       string    `json:"id" db:"id"`
	Todo     string    `json:"todo" db:"todo"`
	IsDone   bool      `json:"is_done" db:"is_done"`
	Created  time.Time `json:"created" db:"created"`
	WorkerId string    `json:"worker_id" db:"worker_id"`
}

type JobReport struct {
	JobId  int    `json:"job_id" db:"job_id"`
	Report string `json:"report" db:"report"`
}

type Upload struct {
	Id       int       `json:"id" db:"id"`
	Filename string    `json:"filename" db:"filename"`
	Type     string    `json:"type" db:"type"`
	Size     uint64    `json:"size" db:"size"`
	Created  time.Time `json:"created" db:"created"`
	WorkerId string    `json:"worker_id" db:"worker_id"`
}
