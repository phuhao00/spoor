package spoor

import (
	"errors"
	"os"
	"time"
)

// Common errors
var (
	ErrInvalidOutputType = errors.New("invalid output type")
)

// SimpleLogger provides a simplified API for common logging needs
type SimpleLogger struct {
	logger Logger
}

// SimpleConfig provides simple configuration options
type SimpleConfig struct {
	Level      LogLevel
	Output     string // "console", "file", "json"
	FilePath   string
	Async      bool
	BatchSize  int
	FlushEvery time.Duration
}

// DefaultSimpleConfig returns default simple configuration
func DefaultSimpleConfig() SimpleConfig {
	return SimpleConfig{
		Level:      LevelInfo,
		Output:     "console",
		FilePath:   "logs/app.log",
		Async:      true,
		BatchSize:  1000,
		FlushEvery: 100 * time.Millisecond,
	}
}

// NewSimple creates a new simple logger
func NewSimple(config SimpleConfig) (*SimpleLogger, error) {
	var writer Writer

	// Create writer based on output type
	switch config.Output {
	case "console":
		writer = NewConsoleWriter(ConsoleWriterConfig{
			Output: os.Stdout,
		})
	case "file":
		fileWriter, err := NewFileWriter(FileWriterConfig{
			LogDir: config.FilePath,
		})
		if err != nil {
			return nil, err
		}
		writer = fileWriter
	case "json":
		writer = NewConsoleWriter(ConsoleWriterConfig{
			Output: os.Stdout,
		})
	default:
		return nil, ErrInvalidOutputType
	}

	// Wrap with batch writer if async is enabled
	if config.Async {
		batchConfig := DefaultBatchConfig()
		batchConfig.BatchSize = config.BatchSize
		batchConfig.FlushInterval = config.FlushEvery
		writer = NewBatchWriter(writer, batchConfig)
	}

	// Create formatter based on output type
	var formatter Formatter
	if config.Output == "json" {
		formatter = NewJSONFormatter()
	} else {
		formatter = NewTextFormatter()
	}

	// Create logger
	var logger Logger
	if config.Async {
		asyncConfig := DefaultAsyncConfig()
		asyncConfig.BufferSize = config.BatchSize
		asyncConfig.FlushInterval = config.FlushEvery
		logger = NewAsyncLogger(writer, config.Level, asyncConfig, WithFormatter(formatter))
	} else {
		logger = NewCoreLogger(writer, config.Level, WithFormatter(formatter))
	}

	return &SimpleLogger{logger: logger}, nil
}

// Quick creates a logger with sensible defaults
func Quick() *SimpleLogger {
	config := DefaultSimpleConfig()
	logger, _ := NewSimple(config)
	return logger
}

// QuickFile creates a file logger with sensible defaults
func QuickFile(filePath string) (*SimpleLogger, error) {
	config := DefaultSimpleConfig()
	config.Output = "file"
	config.FilePath = filePath
	return NewSimple(config)
}

// QuickJSON creates a JSON logger with sensible defaults
func QuickJSON() *SimpleLogger {
	config := DefaultSimpleConfig()
	config.Output = "json"
	logger, _ := NewSimple(config)
	return logger
}

// QuickAsync creates an async logger with sensible defaults
func QuickAsync() *SimpleLogger {
	config := DefaultSimpleConfig()
	config.Async = true
	logger, _ := NewSimple(config)
	return logger
}

// Log methods - simple and clean API

// Debug logs a debug message
func (sl *SimpleLogger) Debug(msg string) {
	sl.logger.Debug(msg)
}

// Info logs an info message
func (sl *SimpleLogger) Info(msg string) {
	sl.logger.Info(msg)
}

// Warn logs a warning message
func (sl *SimpleLogger) Warn(msg string) {
	sl.logger.Warn(msg)
}

// Error logs an error message
func (sl *SimpleLogger) Error(msg string) {
	sl.logger.Error(msg)
}

// Fatal logs a fatal message
func (sl *SimpleLogger) Fatal(msg string) {
	sl.logger.Fatal(msg)
}

// Debugf logs a formatted debug message
func (sl *SimpleLogger) Debugf(format string, args ...interface{}) {
	sl.logger.Debugf(format, args...)
}

// Infof logs a formatted info message
func (sl *SimpleLogger) Infof(format string, args ...interface{}) {
	sl.logger.Infof(format, args...)
}

// Warnf logs a formatted warning message
func (sl *SimpleLogger) Warnf(format string, args ...interface{}) {
	sl.logger.Warnf(format, args...)
}

// Errorf logs a formatted error message
func (sl *SimpleLogger) Errorf(format string, args ...interface{}) {
	sl.logger.Errorf(format, args...)
}

// Fatalf logs a formatted fatal message
func (sl *SimpleLogger) Fatalf(format string, args ...interface{}) {
	sl.logger.Fatalf(format, args...)
}

// WithField returns a logger with a field
func (sl *SimpleLogger) WithField(key string, value interface{}) *SimpleLogger {
	return &SimpleLogger{logger: sl.logger.WithField(key, value)}
}

// WithFields returns a logger with fields
func (sl *SimpleLogger) WithFields(fields map[string]interface{}) *SimpleLogger {
	return &SimpleLogger{logger: sl.logger.WithFields(fields)}
}

// WithError returns a logger with an error
func (sl *SimpleLogger) WithError(err error) *SimpleLogger {
	return &SimpleLogger{logger: sl.logger.WithError(err)}
}

// SetLevel sets the log level
func (sl *SimpleLogger) SetLevel(level LogLevel) {
	sl.logger.SetLevel(level)
}

// GetLevel returns the current log level
func (sl *SimpleLogger) GetLevel() LogLevel {
	return sl.logger.GetLevel()
}

// Sync flushes all buffered logs
func (sl *SimpleLogger) Sync() error {
	return sl.logger.Sync()
}

// Close closes the logger
func (sl *SimpleLogger) Close() error {
	return sl.logger.Close()
}

// GetMetrics returns performance metrics (if available)
func (sl *SimpleLogger) GetMetrics() interface{} {
	if asyncLogger, ok := sl.logger.(*AsyncLogger); ok {
		return asyncLogger.GetMetrics()
	}
	return nil
}

// Global simple logger instance
var (
	globalSimple *SimpleLogger
)

func init() {
	globalSimple = Quick()
}

// SetGlobalSimple sets the global simple logger
func SetGlobalSimple(logger *SimpleLogger) {
	globalSimple = logger
}

// GetGlobalSimple returns the global simple logger
func GetGlobalSimple() *SimpleLogger {
	return globalSimple
}

// Global simple logger convenience functions

// SimpleDebug logs a debug message using the global simple logger
func SimpleDebug(msg string) {
	globalSimple.Debug(msg)
}

// SimpleInfo logs an info message using the global simple logger
func SimpleInfo(msg string) {
	globalSimple.Info(msg)
}

// SimpleWarn logs a warning message using the global simple logger
func SimpleWarn(msg string) {
	globalSimple.Warn(msg)
}

// SimpleError logs an error message using the global simple logger
func SimpleError(msg string) {
	globalSimple.Error(msg)
}

// SimpleFatal logs a fatal message using the global simple logger
func SimpleFatal(msg string) {
	globalSimple.Fatal(msg)
}

// SimpleDebugf logs a formatted debug message using the global simple logger
func SimpleDebugf(format string, args ...interface{}) {
	globalSimple.Debugf(format, args...)
}

// SimpleInfof logs a formatted info message using the global simple logger
func SimpleInfof(format string, args ...interface{}) {
	globalSimple.Infof(format, args...)
}

// SimpleWarnf logs a formatted warning message using the global simple logger
func SimpleWarnf(format string, args ...interface{}) {
	globalSimple.Warnf(format, args...)
}

// SimpleErrorf logs a formatted error message using the global simple logger
func SimpleErrorf(format string, args ...interface{}) {
	globalSimple.Errorf(format, args...)
}

// SimpleFatalf logs a formatted fatal message using the global simple logger
func SimpleFatalf(format string, args ...interface{}) {
	globalSimple.Fatalf(format, args...)
}

// SimpleWithField returns a logger with a field using the global simple logger
func SimpleWithField(key string, value interface{}) *SimpleLogger {
	return globalSimple.WithField(key, value)
}

// SimpleWithFields returns a logger with fields using the global simple logger
func SimpleWithFields(fields map[string]interface{}) *SimpleLogger {
	return globalSimple.WithFields(fields)
}

// SimpleWithError returns a logger with an error using the global simple logger
func SimpleWithError(err error) *SimpleLogger {
	return globalSimple.WithError(err)
}
