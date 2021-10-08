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
