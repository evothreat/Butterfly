package main

const (
	SERVER_ADDR  = "http://127.0.0.1:8080"
	REGISTER_URL = SERVER_ADDR + "/api/workers"
	JOBS_URL     = SERVER_ADDR + "/api/workers/%s/jobs?undone"
	REPORT_URL   = SERVER_ADDR + "/api/workers/%s/jobs/%d/report"
	UPLOADS_URL  = SERVER_ADDR + "/api/workers/%s/uploads"
	HARDWARE_URL = SERVER_ADDR + "/api/workers/%s/hardware"

	MAX_RETRIES = 10

	MIN_DELAY = 3
	MAX_DELAY = 60
)
