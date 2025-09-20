package spoor

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// BaseWriter provides common functionality for all writers
type BaseWriter struct {
	mu            sync.RWMutex
	output        io.Writer
	formatter     Formatter
	batchSize     int
	flushInterval time.Duration
	buffer        []LogEntry
	stopChan      chan struct{}
	closed        bool
}

// NewBaseWriter creates a new base writer
func NewBaseWriter(output io.Writer, formatter Formatter) *BaseWriter {
	return &BaseWriter{
		output:        output,
		formatter:     formatter,
		batchSize:     100,
		flushInterval: 5 * time.Second,
		buffer:        make([]LogEntry, 0, 100),
		stopChan:      make(chan struct{}),
	}
}

// Write implements io.Writer interface
func (w *BaseWriter) Write(p []byte) (n int, err error) {
	// For simple text output, write directly
	return w.output.Write(p)
}

// WriteEntry writes a structured log entry
func (w *BaseWriter) WriteEntry(entry LogEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return fmt.Errorf("writer is closed")
	}

	// Add to buffer
	w.buffer = append(w.buffer, entry)

	// Flush if buffer is full
	if len(w.buffer) >= w.batchSize {
		return w.flushUnsafe()
	}

	return nil
}

// WriteStructured writes a structured log entry (implements StructuredWriter)
func (w *BaseWriter) WriteStructured(entry LogEntry) error {
	return w.WriteEntry(entry)
}

// Flush flushes the buffer
func (w *BaseWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.flushUnsafe()
}

// flushUnsafe flushes the buffer without locking (caller must hold the lock)
func (w *BaseWriter) flushUnsafe() error {
	if len(w.buffer) == 0 {
		return nil
	}

	// Format and write each entry
	for _, entry := range w.buffer {
		data, err := w.formatter.Format(entry)
		if err != nil {
			continue // Skip malformed entries
		}

		if _, err := w.output.Write(data); err != nil {
			return err
		}
	}

	// Clear buffer
	w.buffer = w.buffer[:0]
	return nil
}

// Close closes the writer and flushes remaining data
func (w *BaseWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return nil
	}

	w.closed = true
	close(w.stopChan)

	// Flush remaining data
	return w.flushUnsafe()
}

// SetBatchSize sets the batch size for flushing
func (w *BaseWriter) SetBatchSize(size int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.batchSize = size
}

// SetFlushInterval sets the flush interval
func (w *BaseWriter) SetFlushInterval(interval time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.flushInterval = interval
}

// SetFormatter sets the formatter for the writer
func (w *BaseWriter) SetFormatter(formatter Formatter) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.formatter = formatter
}

// StartFlushLoop starts the automatic flush loop
func (w *BaseWriter) StartFlushLoop() {
	go w.flushLoop()
}

// flushLoop runs the automatic flush loop
func (w *BaseWriter) flushLoop() {
	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.Flush()
		case <-w.stopChan:
			return
		}
	}
}

// GetBufferSize returns the current buffer size
func (w *BaseWriter) GetBufferSize() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.buffer)
}

// IsClosed returns whether the writer is closed
func (w *BaseWriter) IsClosed() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.closed
}
