package logs

import (
	"fmt"
	"log"
	"os"
)

// Logger ...
type logger struct {
	k *log.Logger
}

var c *logger

func init() {
	c = &logger{
		k: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func print(s string) {
	_ = c.k.Output(3, s)
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
