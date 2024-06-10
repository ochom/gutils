package logs

import (
	"fmt"
	"io"
	"log"
	"os"
)

var c *log.Logger

func init() {
	c = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

// SetOutput ...
func SetOutput(w io.Writer) {
	c.SetOutput(w)
}

func print(s string) {
	_ = c.Output(3, s)
}

// Info ...
func Info(format string, args ...any) {
	print("INFO: " + fmt.Sprintf(format, args...))
}

// Warn ...
func Warn(format string, args ...any) {
	print("WARN: " + fmt.Sprintf(format, args...))
}

// Error ...
func Error(format string, args ...any) {
	print("ERROR: " + fmt.Sprintf(format, args...))
}

// Fatal ...
func Fatal(format string, args ...any) {
	print("FATAL: " + fmt.Sprintf(format, args...))
	os.Exit(1)
}
