package utils

import (
	"syscall"
	"unsafe"
)

const (
	MB_ICONWARNING = 0x00000030
	MB_ICONINFO    = 0x00000020
	MB_ICONERROR   = 0x00000010
)

func IntPtr(n int) uintptr {
	return uintptr(n)
}

func StrPtr(s string) uintptr {
	ptr, _ := syscall.UTF16PtrFromString(s)
	return uintptr(unsafe.Pointer(ptr))
}

func ShowMessage(title string, text string, icon string) {
	user32Dll, _ := syscall.LoadLibrary("user32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")
	MessageBox := user32.NewProc("MessageBoxW")
	var icontype int
	if icon == "info" {
		icontype = MB_ICONINFO
	} else if icon == "warning" {
		icontype = MB_ICONWARNING
	} else if icon == "error" {
		icontype = MB_ICONERROR
	}
	MessageBox.Call(IntPtr(0), StrPtr(text), StrPtr(title), IntPtr(icontype))
	defer syscall.FreeLibrary(user32Dll)
}
