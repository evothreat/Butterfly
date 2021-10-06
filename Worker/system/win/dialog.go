package win

import (
	"golang.org/x/sys/windows"
)

// blocks the current goroutine, therefore start in new one!

func ShowInfoDialog(title, text string) {
	textPtr, _ := windows.UTF16PtrFromString(text)
	titlePtr, _ := windows.UTF16PtrFromString(title)
	windows.MessageBox(0, textPtr, titlePtr, windows.MB_OK|windows.MB_ICONINFORMATION)
}
