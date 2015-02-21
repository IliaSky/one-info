package main

// import (
// 	"fmt"
// 	"github.com/CzarekTomczak/cef2go/src/cef"
// 	"github.com/CzarekTomczak/cef2go/src/wingui"
// 	"os"
// )

// func main() {
// 	// os.Chdir(".")
// 	fmt.Print("SASA")
// 	os.Chdir("./../../cef2go/Release")

// 	_ = cef
// 	_ = wingui

// }

import (
	"cef"
	"wingui"
	// "github.com/CzarekTomczak/cef2go/src/cef"
	// "github.com/CzarekTomczak/cef2go/src/wingui"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var Logger *log.Logger = log.New(os.Stdout, "[main] ", log.Lshortfile)

func main() {
	os.Chdir("./../../cef2go/Release")
	hInstance, e := wingui.GetModuleHandle(nil)
	if e != nil {
		wingui.AbortErrNo("GetModuleHandle", e)
	}

	cef.ExecuteProcess(unsafe.Pointer(hInstance))

	settings := cef.Settings{}
	settings.CachePath = "webcache"                // Set to empty to disable
	settings.LogSeverity = cef.LOGSEVERITY_DEFAULT // LOGSEVERITY_VERBOSE
	cef.Initialize(settings)

	wndproc := syscall.NewCallback(WndProc)
	Logger.Println("CreateWindow")
	hwnd := wingui.CreateWindow("cef2go example", wndproc)

	browserSettings := cef.BrowserSettings{}
	// TODO: It should be executable's directory used
	// rather than working directory.
	url, _ := os.Getwd()
	url = "file://" + url + "/example.html"
	cef.CreateBrowser(unsafe.Pointer(hwnd), browserSettings, url)

	// It should be enough to call WindowResized after 10ms,
	// though to be sure let's extend it to 100ms.
	time.AfterFunc(time.Millisecond*100, func() {
		cef.WindowResized(unsafe.Pointer(hwnd))
	})

	cef.RunMessageLoop()
	cef.Shutdown()
	os.Exit(0)
}

func WndProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) (rc uintptr) {
	switch msg {
	case wingui.WM_CREATE:
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	case wingui.WM_SIZE:
		cef.WindowResized(unsafe.Pointer(hwnd))
	case wingui.WM_CLOSE:
		wingui.DestroyWindow(hwnd)
	case wingui.WM_DESTROY:
		cef.QuitMessageLoop()
	default:
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	}
	return
}
