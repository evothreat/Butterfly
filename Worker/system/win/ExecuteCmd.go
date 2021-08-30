package win

import (
	"golang.org/x/sys/windows"
	"os/exec"
)

func ExecuteCommand(args... string) (string, error) {
	cmd := exec.Command("cmd", append([]string{"/C"}, args...)...)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	return string(output), err
}