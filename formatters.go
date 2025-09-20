package spoor

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// TextFormatter formats log entries as plain text
type TextFormatter struct {
	TimestampFormat string
	DisableColors   bool
	DisableCaller   bool
}

// NewTextFormatter creates a new text formatter
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		TimestampFormat: time.RFC3339,
		DisableColors:   false,
		DisableCaller:   false,
	}
}

// Format formats a log entry as text
func (f *TextFormatter) Format(entry LogEntry) ([]byte, error) {
	var b strings.Builder

	// Timestamp
	if f.TimestampFormat != "" {
		b.WriteString(entry.Timestamp.Format(f.TimestampFormat))
		b.WriteString(" ")
	}

	// Level
	levelStr := entry.Level.String()
	if !f.DisableColors {
		levelStr = f.colorizeLevel(entry.Level)
	}
	b.WriteString(levelStr)
	b.WriteString(" ")

	// Caller
	if !f.DisableCaller && entry.Caller != "" {
		b.WriteString(entry.Caller)
		b.WriteString(" ")
	}

	// Message
	b.WriteString(entry.Message)

	// Fields
	if len(entry.Fields) > 0 {
		b.WriteString(" ")
		f.writeFields(&b, entry.Fields)
	}

	b.WriteString("\n")
	return []byte(b.String()), nil
}

// writeFields writes fields to the builder
func (f *TextFormatter) writeFields(b *strings.Builder, fields map[string]interface{}) {
	first := true
	for k, v := range fields {
		if !first {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%s=%v", k, v))
		first = false
	}
}

// colorizeLevel adds color to the log level
func (f *TextFormatter) colorizeLevel(level LogLevel) string {
	switch level {
	case LevelDebug:
		return "\033[36mDEBUG\033[0m" // Cyan
	case LevelInfo:
		return "\033[32mINFO\033[0m" // Green
	case LevelWarn:
		return "\033[33mWARN\033[0m" // Yellow
	case LevelError:
		return "\033[31mERROR\033[0m" // Red
	case LevelFatal:
		return "\033[35mFATAL\033[0m" // Magenta
	default:
		return level.String()
	}
}

// JSONFormatter formats log entries as JSON
type JSONFormatter struct {
	TimestampFormat string
	PrettyPrint     bool
}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     false,
	}
}

// Format formats a log entry as JSON
func (f *JSONFormatter) Format(entry LogEntry) ([]byte, error) {
	// Create a copy of the entry to avoid modifying the original
	jsonEntry := struct {
		Timestamp string                 `json:"timestamp"`
		Level     string                 `json:"level"`
		Message   string                 `json:"message"`
		Caller    string                 `json:"caller,omitempty"`
		Fields    map[string]interface{} `json:"fields,omitempty"`
	}{
		Timestamp: entry.Timestamp.Format(f.TimestampFormat),
		Level:     entry.Level.String(),
		Message:   entry.Message,
		Caller:    entry.Caller,
		Fields:    entry.Fields,
	}

	if f.PrettyPrint {
		return json.MarshalIndent(jsonEntry, "", "  ")
	}
	return json.Marshal(jsonEntry)
}

// LogrusFormatter formats log entries in logrus style
type LogrusFormatter struct {
	TimestampFormat string
	DisableCaller   bool
}

// NewLogrusFormatter creates a new logrus-style formatter
func NewLogrusFormatter() *LogrusFormatter {
	return &LogrusFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DisableCaller:   false,
	}
}

// Format formats a log entry in logrus style
func (f *LogrusFormatter) Format(entry LogEntry) ([]byte, error) {
	var b strings.Builder

	// Timestamp
	b.WriteString(entry.Timestamp.Format(f.TimestampFormat))
	b.WriteString(" ")

	// Level
	b.WriteString(fmt.Sprintf("level=%s", strings.ToLower(entry.Level.String())))
	b.WriteString(" ")

	// Caller
	if !f.DisableCaller && entry.Caller != "" {
		b.WriteString(fmt.Sprintf("caller=%s", entry.Caller))
		b.WriteString(" ")
	}

	// Message
	b.WriteString(fmt.Sprintf("msg=%q", entry.Message))

	// Fields
	if len(entry.Fields) > 0 {
		b.WriteString(" ")
		f.writeFields(&b, entry.Fields)
	}

	b.WriteString("\n")
	return []byte(b.String()), nil
}

// writeFields writes fields to the builder
func (f *LogrusFormatter) writeFields(b *strings.Builder, fields map[string]interface{}) {
	first := true
	for k, v := range fields {
		if !first {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%s=%v", k, v))
		first = false
	}
}
