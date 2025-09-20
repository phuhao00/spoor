package spoor

import (
	"io"
	"os"
	"time"
)

// ConsoleWriter writes logs to console
type ConsoleWriter struct {
	*BaseWriter
	output io.Writer
}

// ConsoleWriterConfig holds configuration for console writer
type ConsoleWriterConfig struct {
	Output        io.Writer
	Formatter     Formatter
	BatchSize     int
	FlushInterval int // in seconds
}

// NewConsoleWriter creates a new console writer
func NewConsoleWriter(config ConsoleWriterConfig) *ConsoleWriter {
	if config.Output == nil {
		config.Output = os.Stdout
	}
	if config.Formatter == nil {
		config.Formatter = NewTextFormatter()
	}

	baseWriter := NewBaseWriter(config.Output, config.Formatter)
	if config.BatchSize > 0 {
		baseWriter.SetBatchSize(config.BatchSize)
	}
	if config.FlushInterval > 0 {
		baseWriter.SetFlushInterval(time.Duration(config.FlushInterval) * time.Second)
	}

	writer := &ConsoleWriter{
		BaseWriter: baseWriter,
		output:     config.Output,
	}

	// Start the flush loop
	writer.StartFlushLoop()

	return writer
}

// NewConsoleWriterWithDefaults creates a console writer with default settings
func NewConsoleWriterWithDefaults() *ConsoleWriter {
	return NewConsoleWriter(ConsoleWriterConfig{})
}

// NewConsoleWriterToStderr creates a console writer that writes to stderr
func NewConsoleWriterToStderr() *ConsoleWriter {
	return NewConsoleWriter(ConsoleWriterConfig{
		Output: os.Stderr,
	})
}

// Write implements io.Writer interface
func (w *ConsoleWriter) Write(p []byte) (n int, err error) {
	return w.output.Write(p)
}

// WriteEntry writes a structured log entry
func (w *ConsoleWriter) WriteEntry(entry LogEntry) error {
	// For console output, we can write immediately or use batching
	return w.BaseWriter.WriteEntry(entry)
}

// SetOutput changes the output writer
func (w *ConsoleWriter) SetOutput(output io.Writer) {
	w.BaseWriter.mu.Lock()
	defer w.BaseWriter.mu.Unlock()
	w.output = output
	w.BaseWriter.output = output
}
