package ptlog

import (
	"fmt"
	"runtime"
	"time"
)

var (
	red     = string([]byte{27, 91, 57, 49, 109})
	green   = string([]byte{27, 91, 57, 50, 109})
	yellow  = string([]byte{27, 91, 57, 51, 109})
	blue    = string([]byte{27, 91, 57, 52, 109})
	magenta = string([]byte{27, 91, 57, 53, 109})
	cyan    = string([]byte{27, 91, 57, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func Error(err error) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	now := time.Now().String()[0:19]
	fmt.Print(cyan, now, reset)
	fmt.Print("[", magenta, "ERROR", reset, "]")
	fmt.Print(blue, f.Name(), "(", line, "): ", reset)
	fmt.Print(red, err.Error(), reset)
	fmt.Print("\n")

}
func Debug(msg interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	now := time.Now().String()[0:19]
	fmt.Print(cyan, now, reset)
	fmt.Print("[", yellow, "INFO", reset, "]")
	fmt.Print(blue, f.Name(), "(", line, "):", reset)
	fmt.Print(" ", msg)
	fmt.Print("\n")
}
