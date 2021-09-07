package api

import "time"

type Worker struct {
	Id       string    `json:"id"`
	Hostname string    `json:"hostname"`
	Country  string    `json:"country"`
	Os       string    `json:"os"`
	IpAddr   string    `json:"ip_addr"`
	IsAdmin  bool      `json:"is_admin"`
	Boost    bool      `json:"boost"`
	LastSeen time.Time `json:"last_seen"`
}

type HardwareInfo struct {
	Gpu      string `json:"gpu"`
	Cpu      string `json:"cpu"`
	Ram      string `json:"ram"`
	WorkerId string `json:"worker_id"`
}

type Job struct {
	Id       string    `json:"id"`
	Todo     string    `json:"todo"`
	IsDone   bool      `json:"is_done"`
	Created  time.Time `json:"created"`
	WorkerId string    `json:"worker_id"`
}

type JobReport struct {
	JobId  int    `json:"job_id"`
	Report string `json:"report"`
}

type Upload struct {
	Id       int       `json:"id"`
	Filename string    `json:"filename"`
	Type     string    `json:"type"`
	Size     uint64    `json:"size"`
	Created  time.Time `json:"created"`
	WorkerId string    `json:"worker_id"`
}
