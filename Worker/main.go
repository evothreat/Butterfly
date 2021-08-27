package main

import (
	"Worker/system/windows"
	"Worker/utils"
	"fmt"
	"os"
	"runtime"
)

func main() {
	hostname, _ := os.Hostname()
	osName, _ := windows.GetOsName()
	cpuName, _ := windows.GetCpuName()
	cpuCores := runtime.NumCPU()
	gpuName, _ := windows.GetGpuName()
	guid, _ := windows.GetMachineGuid()
	fmt.Printf("Hostname: %s\n", hostname)
	fmt.Printf("OS Name: %s\n", osName)
	fmt.Printf("CPU Name: %s\n", cpuName)
	fmt.Printf("CPU Cores: %d\n", cpuCores)
	fmt.Printf("GPU Name: %s\n", gpuName)
	fmt.Printf("Machine Guid: %s\n", guid)
	fmt.Printf("Total RAM: %s\n", utils.ToReadableSize(windows.GetTotalRam()))
}
