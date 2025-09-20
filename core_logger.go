package spoor

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// CoreLogger is the core implementation of the Logger interface
type CoreLogger struct {
	mu        sync.RWMutex
	writer    Writer
	level     LogLevel
	formatter Formatter
	hooks     []Hook
	fields    map[string]interface{}
	caller    bool
}

// NewCoreLogger creates a new core logger instance
func NewCoreLogger(writer Writer, level LogLevel, options ...Option) *CoreLogger {
	logger := &CoreLogger{
		writer:    writer,
		level:     level,
		formatter: &TextFormatter{},
		hooks:     make([]Hook, 0),
		fields:    make(map[string]interface{}),
		caller:    true,
	}

	for _, opt := range options {
		opt(logger)
	}

	return logger
}

// Option defines a function that configures a logger
type Option func(*CoreLogger)

// WithFormatter sets the formatter for the logger
func WithFormatter(formatter Formatter) Option {
	return func(l *CoreLogger) {
		l.formatter = formatter
	}
}

// WithHooks sets the hooks for the logger
func WithHooks(hooks ...Hook) Option {
	return func(l *CoreLogger) {
		l.hooks = append(l.hooks, hooks...)
	}
}

// WithCaller enables or disables caller information
func WithCaller(enable bool) Option {
	return func(l *CoreLogger) {
		l.caller = enable
	}
}

// Debug logs a debug message
func (l *CoreLogger) Debug(msg string) {
	l.log(LevelDebug, msg, nil)
}

// Info logs an info message
func (l *CoreLogger) Info(msg string) {
	l.log(LevelInfo, msg, nil)
}

// Warn logs a warning message
func (l *CoreLogger) Warn(msg string) {
	l.log(LevelWarn, msg, nil)
}

// Error logs an error message
func (l *CoreLogger) Error(msg string) {
	l.log(LevelError, msg, nil)
}

// Fatal logs a fatal message
func (l *CoreLogger) Fatal(msg string) {
	l.log(LevelFatal, msg, nil)
}

// Debugf logs a formatted debug message
func (l *CoreLogger) Debugf(format string, args ...interface{}) {
	l.log(LevelDebug, fmt.Sprintf(format, args...), nil)
}

// Infof logs a formatted info message
func (l *CoreLogger) Infof(format string, args ...interface{}) {
	l.log(LevelInfo, fmt.Sprintf(format, args...), nil)
}

// Warnf logs a formatted warning message
func (l *CoreLogger) Warnf(format string, args ...interface{}) {
	l.log(LevelWarn, fmt.Sprintf(format, args...), nil)
}

// Errorf logs a formatted error message
func (l *CoreLogger) Errorf(format string, args ...interface{}) {
	l.log(LevelError, fmt.Sprintf(format, args...), nil)
}

// Fatalf logs a formatted fatal message
func (l *CoreLogger) Fatalf(format string, args ...interface{}) {
	l.log(LevelFatal, fmt.Sprintf(format, args...), nil)
}

// WithField returns a new logger with the specified field
func (l *CoreLogger) WithField(key string, value interface{}) Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &CoreLogger{
		writer:    l.writer,
		level:     l.level,
		formatter: l.formatter,
		hooks:     l.hooks,
		fields:    newFields,
		caller:    l.caller,
	}
}

// WithFields returns a new logger with the specified fields
func (l *CoreLogger) WithFields(fields map[string]interface{}) Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &CoreLogger{
		writer:    l.writer,
		level:     l.level,
		formatter: l.formatter,
		hooks:     l.hooks,
		fields:    newFields,
		caller:    l.caller,
	}
}

// WithError returns a new logger with the specified error
func (l *CoreLogger) WithError(err error) Logger {
	return l.WithField("error", err.Error())
}

// SetLevel sets the log level
func (l *CoreLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current log level
func (l *CoreLogger) GetLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// log is the internal logging method
func (l *CoreLogger) log(level LogLevel, msg string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    make(map[string]interface{}),
	}

	// Add logger fields
	l.mu.RLock()
	for k, v := range l.fields {
		entry.Fields[k] = v
	}
	l.mu.RUnlock()

	// Add method fields
	for k, v := range fields {
		entry.Fields[k] = v
	}

	// Add caller information if enabled
	if l.caller {
		if caller := getCaller(); caller != "" {
			entry.Caller = caller
		}
	}

	// Fire hooks
	for _, hook := range l.hooks {
		if l.shouldFireHook(hook, level) {
			hook.Fire(entry)
		}
	}

	// Write the log entry
	if structuredWriter, ok := l.writer.(StructuredWriter); ok {
		structuredWriter.WriteStructured(entry)
	} else {
		// Fallback to text format
		if data, err := l.formatter.Format(entry); err == nil {
			l.writer.Write(data)
		}
	}
}

// shouldFireHook checks if a hook should be fired for the given level
func (l *CoreLogger) shouldFireHook(hook Hook, level LogLevel) bool {
	levels := hook.Levels()
	if len(levels) == 0 {
		return true
	}

	for _, hookLevel := range levels {
		if hookLevel == level {
			return true
		}
	}
	return false
}

// getCaller returns the caller information
func getCaller() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return ""
	}

	// Get just the filename, not the full path
	parts := strings.Split(file, "/")
	filename := parts[len(parts)-1]

	return fmt.Sprintf("%s:%d", filename, line)
}

// Sync flushes all buffered log entries
func (l *CoreLogger) Sync() error {
	return l.writer.Flush()
}

// Close closes the logger and flushes any pending logs
func (l *CoreLogger) Close() error {
	return l.writer.Close()
}

// SetFormatter sets the formatter for the logger's writer
func (l *CoreLogger) SetFormatter(formatter Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if baseWriter, ok := l.writer.(*BaseWriter); ok {
		baseWriter.SetFormatter(formatter)
	}
}

// SetWriter sets the output writer for the logger
func (l *CoreLogger) SetWriter(writer Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer = writer
}
