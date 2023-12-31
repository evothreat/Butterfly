package win

import (
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"unsafe"
)

type MemoryStatusEx struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	ullAvailPhys uint64
	restInfo     [5]uint64
}

func GetOsName() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
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

func GetTotalRam() (uint64, error) {
	user32dll, err := windows.LoadDLL("kernel32.dll") // TODO: export in own file
	if err != nil {
		return 0, err
	}
	msx := &MemoryStatusEx{
		dwLength: 64,
	}
	r, _, err := user32dll.MustFindProc("GlobalMemoryStatusEx").Call(uintptr(unsafe.Pointer(msx)))
	if r == 0 {
		return 0, err
	}
	return msx.ullTotalPhys, nil
}

// Maybe implement Disc Drives/Free Space (hint: GetLogicalDrives, GetDiskFreeSpaceEx)
