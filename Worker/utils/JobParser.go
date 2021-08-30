package utils

import "strings"

type JobType int

const (
	UNKNOWN JobType = iota
	SHELL_CMD
	UPLOAD
	DOWNLOAD
	DDOS
	SLEEP
	//CHDIR
)

func ParseJob(jobStr string) (JobType, []string) {
	values := strings.Fields(jobStr)
	valuesN := len(values)
	if valuesN == 0 {
		return UNKNOWN, nil
	}
	jobTypeStr := values[0]
	values = values[1:]
	valuesN--
	if jobTypeStr == "cmd" {
		return SHELL_CMD, values
	} else if jobTypeStr == "upload" && valuesN == 1 {
		return UPLOAD, values
	} else if jobTypeStr == "download" && valuesN == 2 {
		return DOWNLOAD, values
	} else if jobTypeStr == "sleep" && valuesN == 1 {
		return SLEEP, values
	}
	return UNKNOWN, nil
}
