package win

import (
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"unsafe"
)

func ProcessHasAdminRights() bool { // TODO: check for errors?
	sid, _ := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	token := windows.GetCurrentProcessToken()
	isAdmin, _ := token.IsMember(sid)
	return isAdmin
}

func getProcessPathByPid(processId uint32) (string, error) {
	proc, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, processId)
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(proc)
	pathLen := windows.MAX_PATH
	var pathBuf [windows.MAX_PATH]uint16
	err = windows.QueryFullProcessImageName(proc, windows.PROCESS_NAME_NATIVE, (*uint16)(unsafe.Pointer(&pathBuf)),
		(*uint32)(unsafe.Pointer(&pathLen)))
	if err != nil {
		return "", err
	}
	return windows.UTF16ToString(pathBuf[:]), nil
}

func ProcessAlreadyRunning() (bool, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return false, err
	}
	defer windows.CloseHandle(snapshot)

	currProcId := windows.GetCurrentProcessId()
	currProcPath, _ := os.Executable()
	currProcName := filepath.Base(currProcPath)

	procEntry := &windows.ProcessEntry32{}
	procEntry.Size = uint32(unsafe.Sizeof(*procEntry))
	for err := windows.Process32First(snapshot, procEntry); err == nil; err = windows.Process32Next(snapshot, procEntry) {
		// compare error with ERROR_NO_MORE_FILES
		if windows.UTF16ToString(procEntry.ExeFile[:]) == currProcName && procEntry.ProcessID != currProcId {
			if path, _ := getProcessPathByPid(procEntry.ProcessID); path != currProcPath {
				return true, nil
			}
		}
	}
	return false, nil
}
