package spoor

import (
	"fmt"
	"io"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	LevelDebug LogLevel = iota + 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLogLevel parses a string to LogLevel
func ParseLogLevel(level string) (LogLevel, error) {
	switch level {
	case "debug", "DEBUG":
		return LevelDebug, nil
	case "info", "INFO":
		return LevelInfo, nil
	case "warn", "WARN", "warning", "WARNING":
		return LevelWarn, nil
	case "error", "ERROR":
		return LevelError, nil
	case "fatal", "FATAL":
		return LevelFatal, nil
	default:
		return 0, fmt.Errorf("invalid log level: %s", level)
	}
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Caller    string                 `json:"caller,omitempty"`
}

// Writer defines the interface for log writers
type Writer interface {
	io.Writer
	WriteEntry(entry LogEntry) error
	Flush() error
	Close() error
}

// Logger defines the core logging interface
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger

	SetLevel(level LogLevel)
	GetLevel() LogLevel
	SetFormatter(formatter Formatter)
	SetWriter(writer Writer)
	Sync() error  // Flushes all buffered log entries
	Close() error // Closes the logger and flushes any pending logs
}

// Configurable defines the interface for configurable components
type Configurable interface {
	Configure(config interface{}) error
}

// BatchWriterInterface defines the interface for writers that support batch operations
type BatchWriterInterface interface {
	Writer
	Flush() error
	SetBatchSize(size int)
	SetFlushInterval(interval time.Duration)
}

// StructuredWriter defines the interface for writers that support structured logging
type StructuredWriter interface {
	Writer
	WriteStructured(entry LogEntry) error
}

// Formatter defines the interface for log formatters
type Formatter interface {
	Format(entry LogEntry) ([]byte, error)
}

// Hook defines the interface for log hooks
type Hook interface {
	Fire(entry LogEntry) error
	Levels() []LogLevel
}
