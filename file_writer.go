package spoor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileWriter writes logs to files with rotation support
type FileWriter struct {
	*BaseWriter
	mu            sync.RWMutex
	file          *os.File
	writer        *bufio.Writer
	logDir        string
	maxSize       int64
	currentSize   int64
	rotationCount int
}

// FileWriterConfig holds configuration for file writer
type FileWriterConfig struct {
	LogDir        string
	MaxSize       int64
	Formatter     Formatter
	BatchSize     int
	FlushInterval int // in seconds
}

// NewFileWriter creates a new file writer
func NewFileWriter(config FileWriterConfig) (*FileWriter, error) {
	if config.Formatter == nil {
		config.Formatter = NewTextFormatter()
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	writer := &FileWriter{
		logDir:  config.LogDir,
		maxSize: config.MaxSize,
	}

	baseWriter := NewBaseWriter(writer, config.Formatter)
	if config.BatchSize > 0 {
		baseWriter.SetBatchSize(config.BatchSize)
	}
	if config.FlushInterval > 0 {
		baseWriter.SetFlushInterval(time.Duration(config.FlushInterval) * time.Second)
	}
	writer.BaseWriter = baseWriter

	// Initialize the first log file
	if err := writer.rotateFile(); err != nil {
		return nil, fmt.Errorf("failed to create initial log file: %w", err)
	}

	// Start the flush loop
	writer.StartFlushLoop()

	return writer, nil
}

// NewFileWriterWithDefaults creates a file writer with default settings
func NewFileWriterWithDefaults(logDir string) (*FileWriter, error) {
	return NewFileWriter(FileWriterConfig{
		LogDir:        logDir,
		MaxSize:       100 * 1024 * 1024, // 100MB
		BatchSize:     100,
		FlushInterval: 5,
	})
}

// Write implements io.Writer interface
func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return 0, fmt.Errorf("file writer is closed")
	}

	// Check if we need to rotate
	if w.maxSize > 0 && w.currentSize+int64(len(p)) > w.maxSize {
		if err := w.rotateFileUnsafe(); err != nil {
			return 0, err
		}
	}

	// Write to file
	n, err = w.writer.Write(p)
	if err != nil {
		return n, err
	}

	w.currentSize += int64(n)
	return n, nil
}

// WriteEntry writes a structured log entry
func (w *FileWriter) WriteEntry(entry LogEntry) error {
	// Format the entry
	data, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}

	// Write the formatted data
	_, err = w.Write(data)
	return err
}

// rotateFile rotates the log file
func (w *FileWriter) rotateFile() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.rotateFileUnsafe()
}

// rotateFileUnsafe rotates the log file without locking
func (w *FileWriter) rotateFileUnsafe() error {
	// Close current file
	if w.file != nil {
		w.writer.Flush()
		w.file.Close()
	}

	// Create new file
	filename := w.generateFilename()
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	w.file = file
	w.writer = bufio.NewWriter(file)
	w.currentSize = 0
	w.rotationCount++

	// Write header
	header := fmt.Sprintf("Log file created at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	w.writer.WriteString(header)
	w.currentSize += int64(len(header))

	return nil
}

// generateFilename generates a unique filename for the log file
func (w *FileWriter) generateFilename() string {
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	return filepath.Join(w.logDir, fmt.Sprintf("app-%s-%d.log", timestamp, w.rotationCount))
}

// Flush flushes the writer buffer
func (w *FileWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.writer != nil {
		return w.writer.Flush()
	}
	return nil
}

// Close closes the file writer
func (w *FileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return nil
	}

	// Flush and close
	if w.writer != nil {
		w.writer.Flush()
	}

	err := w.file.Close()
	w.file = nil
	w.writer = nil

	return err
}

// GetCurrentFile returns the current log file path
func (w *FileWriter) GetCurrentFile() string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.file != nil {
		return w.file.Name()
	}
	return ""
}

// GetCurrentSize returns the current file size
func (w *FileWriter) GetCurrentSize() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.currentSize
}

// GetRotationCount returns the number of rotations
func (w *FileWriter) GetRotationCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.rotationCount
}
