package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

var c *log.Logger

func init() {
	c = log.New(os.Stdout, "", log.LstdFlags)
}

// SetOutput ...
func SetOutput(w io.Writer) {
	c.SetOutput(w)
}

func print(s string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		s = fmt.Sprintf("%s:%d %s", file, line, s)
	} else {
		s = fmt.Sprintf("%s %s", file, s)
	}

	c.Println(s)
}

// Debug ...
func Debug(format string, args ...any) {
	print(fmt.Sprintf("%sDEBUG%s: %s", ColorBlue, ColorReset, fmt.Sprintf(format, args...)))
}

// Info ...
func Info(format string, args ...any) {
	print(fmt.Sprintf("%sINFO%s: %s", ColorGreen, ColorReset, fmt.Sprintf(format, args...)))
}

// Warn ...
func Warn(format string, args ...any) {
	print(fmt.Sprintf("%sWARN%s: %s", ColorYellow, ColorReset, fmt.Sprintf(format, args...)))
}

// Error ...
func Error(format string, args ...any) {
	print(fmt.Sprintf("%sERROR%s: %s", ColorRed, ColorReset, fmt.Sprintf(format, args...)))
}

// Fatal ...
func Fatal(format string, args ...any) {
	print(fmt.Sprintf("%sFATAL%s: %s", ColorRed, ColorReset, fmt.Sprintf(format, args...)))
	os.Exit(1)
}
