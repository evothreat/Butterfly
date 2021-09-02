package win

import (
	"golang.org/x/sys/windows"
	"os/exec"
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
