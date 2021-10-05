package win

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

type MsgBoxType int

const (
	MB_OK MsgBoxType = iota
	MB_OKCANCEL
	MB_YESNOCANCEL
	MB_YESNO
)

// blocks the current goroutine, therefore start in new one!

func ShowMessageBox(hwnd uintptr, title, msg string, mbType MsgBoxType) {
	titlePtr, _ := windows.UTF16PtrFromString(title)
	msgPtr, _ := windows.UTF16PtrFromString(msg)
	user32dll := windows.MustLoadDLL("user32.dll")
	defer user32dll.Release()
	user32dll.MustFindProc("MessageBoxW").Call(hwnd,
		uintptr(unsafe.Pointer(msgPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(mbType))
}
