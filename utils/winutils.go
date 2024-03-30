package utils

import (
	"fmt"
	"syscall"
	"unsafe"
)

// c&p von https://gist.github.com/NaniteFactory/0bd94e84bbe939cda7201374a0c261fd
// MessageBox of Win32 API.
func MessageBox(hwnd uintptr, caption, title string, flags uint) int {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(flags))

	return int(ret)
}

func ThrowErrorAndQuit(context string, err error) {
	const (
		MB_ICONERROR = 0x00000010
		MB_OK        = 0x00000000
	)
	MessageBox(0, "Scriptprox", fmt.Sprintf("%s\n\n%s", context, err), MB_ICONERROR|MB_OK)
}
