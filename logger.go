package spoor

import (
	"fmt"
)

// Logger is the main logger interface
// type Logger = Logger // 这行会导致循环定义，删除

// New creates a new logger with the specified writer and level
func New(writer Writer, level LogLevel, options ...Option) Logger {
	return NewCoreLogger(writer, level, options...)
}

// NewAsync creates a new async logger with the specified writer and level
func NewAsync(writer Writer, level LogLevel, config AsyncLoggerConfig, options ...Option) Logger {
	return NewAsyncLogger(writer, level, config, options...)
}

// NewSimpleLogger creates a new simple logger with the specified configuration
func NewSimpleLogger(config SimpleConfig) (*SimpleLogger, error) {
	return NewSimple(config)
}

// NewAdvanced creates a new advanced logger with the specified configuration
func NewAdvanced(writer Writer, level LogLevel, config AdvancedConfig, options ...Option) Logger {
	return NewAdvancedLogger(writer, level, config, options...)
}

// NewWithDefaults creates a new logger with default settings
func NewWithDefaults() Logger {
	writer := NewWriterFactory().CreateConsoleWriterToStdout()
	return NewCoreLogger(writer, LevelInfo)
}

// NewConsole creates a new console logger
func NewConsole(level LogLevel, options ...Option) Logger {
	writer := NewWriterFactory().CreateConsoleWriterToStdout()
	return NewCoreLogger(writer, level, options...)
}

// NewFile creates a new file logger
func NewFile(logDir string, level LogLevel, options ...Option) (Logger, error) {
	writer, err := NewWriterFactory().CreateFileWriterWithDefaults(logDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create file writer: %w", err)
	}
	return NewCoreLogger(writer, level, options...), nil
}

// NewElastic creates a new Elasticsearch logger
func NewElastic(url, index string, level LogLevel, options ...Option) Logger {
	writer := NewWriterFactory().CreateElasticWriterWithDefaults(url, index)
	return NewCoreLogger(writer, level, options...)
}

// NewClickHouse creates a new ClickHouse logger
func NewClickHouse(dsn, tableName string, level LogLevel, options ...Option) (Logger, error) {
	writer, err := NewWriterFactory().CreateClickHouseWriterWithDefaults(dsn, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to create ClickHouse writer: %w", err)
	}
	return NewCoreLogger(writer, level, options...), nil
}

// NewJSON creates a new logger with JSON formatting
func NewJSON(writer Writer, level LogLevel, options ...Option) Logger {
	opts := append(options, WithFormatter(NewJSONFormatter()))
	return NewCoreLogger(writer, level, opts...)
}

// NewText creates a new logger with text formatting
func NewText(writer Writer, level LogLevel, options ...Option) Logger {
	opts := append(options, WithFormatter(NewTextFormatter()))
	return NewCoreLogger(writer, level, opts...)
}

// NewLogrus creates a new logger with logrus-style formatting
func NewLogrus(writer Writer, level LogLevel, options ...Option) Logger {
	opts := append(options, WithFormatter(NewLogrusFormatter()))
	return NewCoreLogger(writer, level, opts...)
}

// DefaultLogger is the default logger instance
var DefaultLogger Logger

func init() {
	DefaultLogger = NewWithDefaults()
}

// SetDefaultLogger sets the default logger
func SetDefaultLogger(logger Logger) {
	DefaultLogger = logger
}

// GetDefaultLogger returns the default logger
func GetDefaultLogger() Logger {
	return DefaultLogger
}

// Convenience functions for the default logger

// Debug logs a debug message using the default logger
func Debug(msg string) {
	DefaultLogger.Debug(msg)
}

// Info logs an info message using the default logger
func Info(msg string) {
	DefaultLogger.Info(msg)
}

// Warn logs a warning message using the default logger
func Warn(msg string) {
	DefaultLogger.Warn(msg)
}

// Error logs an error message using the default logger
func Error(msg string) {
	DefaultLogger.Error(msg)
}

// Fatal logs a fatal message using the default logger
func Fatal(msg string) {
	DefaultLogger.Fatal(msg)
}

// Debugf logs a formatted debug message using the default logger
func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

// Infof logs a formatted info message using the default logger
func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

// Warnf logs a formatted warning message using the default logger
func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

// Errorf logs a formatted error message using the default logger
func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

// Fatalf logs a formatted fatal message using the default logger
func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

// WithField returns a new logger with the specified field using the default logger
func WithField(key string, value interface{}) Logger {
	return DefaultLogger.WithField(key, value)
}

// WithFields returns a new logger with the specified fields using the default logger
func WithFields(fields map[string]interface{}) Logger {
	return DefaultLogger.WithFields(fields)
}

// WithError returns a new logger with the specified error using the default logger
func WithError(err error) Logger {
	return DefaultLogger.WithError(err)
}

// SetLevel sets the log level for the default logger
func SetLevel(level LogLevel) {
	DefaultLogger.SetLevel(level)
}

// GetLevel returns the current log level of the default logger
func GetLevel() LogLevel {
	return DefaultLogger.GetLevel()
}
