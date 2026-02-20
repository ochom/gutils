package logs

import (
	"fmt"
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

// LogLevel ...
type LogLevel int

// LogLevels ...
const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Config ...
type Config struct {
	Loggers       map[LogLevel]*log.Logger
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

// InitLogger ...
func InitLogger(config ...Config) {
	var cfg Config
	if len(config) == 0 {
		cfg = config[0]
	} else {
		cfg = defaultConfig
	}

	logger = &cfg
}

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

// Debug ...
func Debug(format string, args ...any) {
	print(DebugLevel, fmt.Sprintf("%sDEBUG%s: %s", ColorBlue, ColorReset, fmt.Sprintf(format, args...)))
}

// Info ...
func Info(format string, args ...any) {
	print(InfoLevel, fmt.Sprintf("%sINFO%s: %s", ColorGreen, ColorReset, fmt.Sprintf(format, args...)))
}

// Warn ...
func Warn(format string, args ...any) {
	print(WarnLevel, fmt.Sprintf("%sWARN%s: %s", ColorYellow, ColorReset, fmt.Sprintf(format, args...)))
}

// Error ...
func Error(format string, args ...any) {
	print(ErrorLevel, fmt.Sprintf("%sERROR%s: %s", ColorRed, ColorReset, fmt.Sprintf(format, args...)))
}

// Fatal ...
func Fatal(format string, args ...any) {
	print(FatalLevel, fmt.Sprintf("%sFATAL%s: %s", ColorRed, ColorReset, fmt.Sprintf(format, args...)))
	os.Exit(1)
}
