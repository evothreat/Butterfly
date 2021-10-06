package win

import (
	"golang.org/x/sys/windows"
)

func ShowInfoDialog(title, text string) {
	textPtr, _ := windows.UTF16PtrFromString(text)
	titlePtr, _ := windows.UTF16PtrFromString(title)
	go windows.MessageBox(0, textPtr, titlePtr, windows.MB_OK|windows.MB_ICONINFORMATION)
}
