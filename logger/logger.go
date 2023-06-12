package logger

import (
	"fmt"
	"log"
	"os"
)

// Logger ...
type Logger struct {
	k *log.Logger
}

// NewLogger ...
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)}
}

// WithFile ...
func (c *Logger) WithFile(file string) *Logger {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		os.Exit(1)
	}

	c.k.SetOutput(f)
	return c
}

// WithPrefix ...
func (c *Logger) WithPrefix(prefix string) *Logger {
	c.k.SetPrefix(prefix)
	return c
}

func (c *Logger) print(s string) {
	_ = c.k.Output(3, s)
}

// Info ...
func (c *Logger) Info(format string, args ...any) {
	c.print("INFO: " + fmt.Sprintf(format, args...))
}

// Warn ...
func (c *Logger) Warn(format string, args ...any) {
	c.print("WARN: " + fmt.Sprintf(format, args...))
}

// Error ...
func (c *Logger) Error(format string, args ...any) {
	c.print("ERROR: " + fmt.Sprintf(format, args...))
}

// Fatal ...
func (c *Logger) Fatal(format string, args ...any) {
	c.print("FATAL: " + fmt.Sprintf(format, args...))
	os.Exit(1)
}
