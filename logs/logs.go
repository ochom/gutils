// Package logs provides a simple colored logging utility with multiple log levels.
//
// The package outputs logs to stdout with color-coded prefixes for easy visual
// identification. Each log message includes the file name and line number where
// the log was called.
//
// Log levels (in order of severity):
//   - Debug: Detailed information for debugging (blue)
//   - Info: General operational information (green)
//   - Warn: Warning conditions (yellow)
//   - Error: Error conditions (red)
//   - Fatal: Critical errors that terminate the program (red, calls os.Exit(1))
//
// Example usage:
//
//	logs.Debug("Processing request ID: %s", requestID)
//	logs.Info("Server started on port %d", port)
//	logs.Warn("Deprecated API endpoint called: %s", endpoint)
//	logs.Error("Failed to connect to database: %v", err)
//	logs.Fatal("Cannot start server: %v", err) // Exits program
package logs

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// ANSI color codes for terminal output
const (
	// ColorReset resets the terminal color to default
	ColorReset = "\033[0m"
	// ColorRed for error and fatal messages
	ColorRed = "\033[31m"
	// ColorGreen for info messages
	ColorGreen = "\033[32m"
	// ColorYellow for warning messages
	ColorYellow = "\033[33m"
	// ColorBlue for debug messages
	ColorBlue = "\033[34m"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

// Log level constants in order of increasing severity.
const (
	// DebugLevel for detailed debugging information
	DebugLevel LogLevel = iota
	// InfoLevel for general operational information
	InfoLevel
	// WarnLevel for warning conditions
	WarnLevel
	// ErrorLevel for error conditions
	ErrorLevel
	// FatalLevel for critical errors (terminates program)
	FatalLevel
)

// Config holds the logging configuration.
// Custom loggers can be provided for each log level.
type Config struct {
	// Loggers maps log levels to custom log.Logger instances
	Loggers map[LogLevel]*log.Logger
	// DefaultLogger is used when no specific logger is configured for a level
	DefaultLogger *log.Logger
}

var defaultConfig = Config{
	Loggers:       map[LogLevel]*log.Logger{},
	DefaultLogger: log.New(os.Stdout, "", log.LstdFlags),
}

var logger *Config

func init() {
	logger = &defaultConfig
}

// InitLogger initializes the logging system with custom configuration.
// If no configuration is provided, uses the default configuration.
//
// Example:
//
//	// Use default configuration
//	logs.InitLogger()
//
//	// Custom configuration with file logger for errors
//	errorFile, _ := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	logs.InitLogger(logs.Config{
//		Loggers: map[logs.LogLevel]*log.Logger{
//			logs.ErrorLevel: log.New(errorFile, "", log.LstdFlags),
//		},
//	})
func InitLogger(config ...Config) {
	var cfg Config
	if len(config) == 0 {
		cfg = config[0]
	} else {
		cfg = defaultConfig
	}

	if cfg.DefaultLogger == nil {
		cfg.DefaultLogger = log.New(os.Stdout, "", log.LstdFlags)
	}

	logger = &cfg
}

// print outputs a log message with file and line number information.
func print(l LogLevel, s string) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		s = fmt.Sprintf("%s:%d %s", file, line, s)
	} else {
		s = fmt.Sprintf("%s %s", file, s)
	}

	if log, ok := logger.Loggers[l]; ok {
		log.Println(s)
	} else {
		logger.DefaultLogger.Println(s)
	}
}

// Debug logs a debug message (blue color).
// Use for detailed information useful during development and debugging.
//
// Example:
//
//	logs.Debug("User %s requested endpoint %s", userID, endpoint)
//	logs.Debug("Cache hit for key: %s", key)
//	logs.Debug("Request payload: %v", payload)
func Debug(format string, args ...any) {
	print(DebugLevel, fmt.Sprintf("%sDEBUG%s: %s", ColorBlue, ColorReset, fmt.Sprintf(format, args...)))
}

// Info logs an informational message (green color).
// Use for general operational information about the application's state.
//
// Example:
//
//	logs.Info("Server started on port %d", port)
//	logs.Info("Connected to database: %s", dbName)
//	logs.Info("Processing batch of %d items", count)
func Info(format string, args ...any) {
	print(InfoLevel, fmt.Sprintf("%sINFO%s: %s", ColorGreen, ColorReset, fmt.Sprintf(format, args...)))
}

// Warn logs a warning message (yellow color).
// Use for potentially harmful situations or deprecated functionality.
//
// Example:
//
//	logs.Warn("Cache miss for frequently accessed key: %s", key)
//	logs.Warn("Deprecated endpoint called: %s", endpoint)
//	logs.Warn("Retry attempt %d of %d", attempt, maxRetries)
func Warn(format string, args ...any) {
	print(WarnLevel, fmt.Sprintf("%sWARN%s: %s", ColorYellow, ColorReset, fmt.Sprintf(format, args...)))
}

// Error logs an error message (red color).
// Use for error conditions that should be investigated but don't require immediate shutdown.
//
// Example:
//
//	logs.Error("Failed to send email: %v", err)
//	logs.Error("Database query failed: %v", err)
//	logs.Error("Invalid request from IP %s: %v", ip, err)
func Error(format string, args ...any) {
	print(ErrorLevel, fmt.Sprintf("%sERROR%s: %s", ColorRed, ColorReset, fmt.Sprintf(format, args...)))
}

// Fatal logs a fatal error message (red color) and terminates the program with os.Exit(1).
// Use only for unrecoverable errors that require immediate program termination.
//
// WARNING: This function does not return. Any deferred functions will not be executed.
//
// Example:
//
//	if db == nil {
//		logs.Fatal("Cannot start: database connection required")
//	}
//
//	if err := config.Load(); err != nil {
//		logs.Fatal("Failed to load configuration: %v", err)
//	}
func Fatal(format string, args ...any) {
	print(FatalLevel, fmt.Sprintf("%sFATAL%s: %s", ColorRed, ColorReset, fmt.Sprintf(format, args...)))
	os.Exit(1)
}
