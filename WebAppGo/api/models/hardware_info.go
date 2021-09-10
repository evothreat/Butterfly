package models

const HardwareInfoSchema = `CREATE TABLE IF NOT EXISTS hardware_infos
  (
     gpu       VARCHAR(50),
     cpu       VARCHAR(65),
     ram       VARCHAR(10),
     worker_id VARCHAR(22) PRIMARY KEY,
	 FOREIGN KEY(worker_id) REFERENCES workers(id)
  );`

type HardwareInfo struct {
	Gpu      string `json:"gpu" db:"gpu"`
	Cpu      string `json:"cpu" db:"cpu"`
	Ram      string `json:"ram" db:"ram"`
	WorkerId string `json:"worker_id" db:"worker_id"`
}