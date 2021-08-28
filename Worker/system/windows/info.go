package windows

import (
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"unsafe"
)

type MemStatusEx struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	ullAvailPhys uint64
	unused       [5]uint64
}

func GetOsName() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return "", err // TODO: handle all errors! Return error message! Write to report
	}
	defer k.Close()
	res, _, _ := k.GetStringValue("ProductName")
	return res, nil
}

func GetCpuName() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `HARDWARE\DESCRIPTION\System\CentralProcessor\0`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	res, _, _ := k.GetStringValue("ProcessorNameString")
	return res, nil
}

func GetGpuName() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\WinSAT`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	res, _, _ := k.GetStringValue("PrimaryAdapterString")
	return res, nil
}

func GetMachineGuid() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	res, _, _ := k.GetStringValue("MachineGuid")
	return res, nil
}

func GetTotalRam() uint64 {
	user32dll := windows.NewLazyDLL("kernel32.dll")
	globalMemoryStatusEx := user32dll.NewProc("GlobalMemoryStatusEx")
	msx := &MemStatusEx{
		dwLength: 64,
	}
	r, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(msx)))
	if r == 0 {
		return 0
	}
	return msx.ullTotalPhys
}

func IsAdmin() bool {
	sid, _ := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	token := windows.GetCurrentProcessToken()
	isAdmin, _ := token.IsMember(sid)
	return isAdmin
}

/*func IsAdmin() bool {
	fd, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil && os.IsPermission(err) {
		return false
	}
	defer fd.Close()
	return true
}*/

// Maybe implement Disc Drives/Free Space (hint: GetLogicalDrives, GetDiskFreeSpaceEx)
// add MessageBox show
