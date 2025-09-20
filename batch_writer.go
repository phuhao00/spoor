package spoor

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Common errors
var (
	ErrWriterClosed = errors.New("writer is closed")
)

// BatchWriterConfig configures the batch writer
type BatchWriterConfig struct {
	BatchSize     int           // Number of entries to batch before flushing
	FlushInterval time.Duration // Maximum time between flushes
	MaxRetries    int           // Maximum number of retry attempts
	RetryDelay    time.Duration // Delay between retries
}

// DefaultBatchConfig returns default batch writer configuration
func DefaultBatchConfig() BatchWriterConfig {
	return BatchWriterConfig{
		BatchSize:     1000,
		FlushInterval: 100 * time.Millisecond,
		MaxRetries:    3,
		RetryDelay:    10 * time.Millisecond,
	}
}

// BatchWriter wraps a writer with batching capabilities
type BatchWriter struct {
	writer   Writer
	config   BatchWriterConfig
	batch    []LogEntry
	mu       sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	closed   int32
	metrics  *BatchMetrics
}

// BatchMetrics tracks batch writer performance
type BatchMetrics struct {
	TotalBatches   int64
	TotalEntries   int64
	FailedBatches  int64
	RetryCount     int64
	LastFlushTime  time.Time
	AverageBatchSize float64
}

// NewBatchWriter creates a new batch writer
func NewBatchWriter(writer Writer, config BatchWriterConfig) *BatchWriter {
	ctx, cancel := context.WithCancel(context.Background())
	
	bw := &BatchWriter{
		writer:  writer,
		config:  config,
		batch:   make([]LogEntry, 0, config.BatchSize),
		ctx:     ctx,
		cancel:  cancel,
		metrics: &BatchMetrics{},
	}

	// Start flush goroutine
	bw.wg.Add(1)
	go bw.flushLoop()

	return bw
}

// WriteEntry adds an entry to the batch
func (bw *BatchWriter) WriteEntry(entry LogEntry) error {
	if atomic.LoadInt32(&bw.closed) == 1 {
		return ErrWriterClosed
	}

	bw.mu.Lock()
	defer bw.mu.Unlock()

	bw.batch = append(bw.batch, entry)
	atomic.AddInt64(&bw.metrics.TotalEntries, 1)

	// Flush if batch is full
	if len(bw.batch) >= bw.config.BatchSize {
		return bw.flushUnsafe()
	}

	return nil
}

// WriteStructured implements StructuredWriter interface
func (bw *BatchWriter) WriteStructured(entry LogEntry) error {
	return bw.WriteEntry(entry)
}

// Write implements io.Writer interface
func (bw *BatchWriter) Write(p []byte) (n int, err error) {
	// For simple text output, write directly
	return bw.writer.Write(p)
}

// Flush flushes the current batch
func (bw *BatchWriter) Flush() error {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	return bw.flushUnsafe()
}

// flushUnsafe flushes the batch without locking
func (bw *BatchWriter) flushUnsafe() error {
	if len(bw.batch) == 0 {
		return nil
	}

	// Create a copy of the batch to avoid holding the lock
	batch := make([]LogEntry, len(bw.batch))
	copy(batch, bw.batch)
	
	// Clear the batch
	bw.batch = bw.batch[:0]

	// Flush the batch
	return bw.flushBatch(batch)
}

// flushBatch flushes a batch of entries
func (bw *BatchWriter) flushBatch(batch []LogEntry) error {
	if len(batch) == 0 {
		return nil
	}

	var lastErr error
	retries := 0

	for retries <= bw.config.MaxRetries {
		if err := bw.writeBatch(batch); err != nil {
			lastErr = err
			atomic.AddInt64(&bw.metrics.RetryCount, 1)
			retries++
			
			if retries <= bw.config.MaxRetries {
				time.Sleep(bw.config.RetryDelay)
				continue
			}
		} else {
			// Success
			atomic.AddInt64(&bw.metrics.TotalBatches, 1)
			bw.metrics.LastFlushTime = time.Now()
			
	// Update average batch size
	avgSize := float64(atomic.LoadInt64(&bw.metrics.TotalEntries)) / 
			  float64(atomic.LoadInt64(&bw.metrics.TotalBatches))
	bw.metrics.AverageBatchSize = avgSize
			
			return nil
		}
	}

	// All retries failed
	atomic.AddInt64(&bw.metrics.FailedBatches, 1)
	return lastErr
}

// writeBatch writes a batch of entries to the underlying writer
func (bw *BatchWriter) writeBatch(batch []LogEntry) error {
	// Try structured writer first
	if structuredWriter, ok := bw.writer.(StructuredWriter); ok {
		for _, entry := range batch {
			if err := structuredWriter.WriteStructured(entry); err != nil {
				return err
			}
		}
		return nil
	}

	// Fallback to regular writer
	for _, entry := range batch {
		// This is a simplified version - in practice, you'd need a formatter
		// For now, we'll just write the message
		if _, err := bw.writer.Write([]byte(entry.Message + "\n")); err != nil {
			return err
		}
	}

	return nil
}

// flushLoop handles periodic flushing
func (bw *BatchWriter) flushLoop() {
	defer bw.wg.Done()
	
	ticker := time.NewTicker(bw.config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bw.Flush()
		case <-bw.ctx.Done():
			// Final flush before shutdown
			bw.Flush()
			return
		}
	}
}

// Close closes the batch writer
func (bw *BatchWriter) Close() error {
	if !atomic.CompareAndSwapInt32(&bw.closed, 0, 1) {
		return nil
	}

	// Cancel context to stop flush loop
	bw.cancel()
	
	// Wait for flush loop to finish
	bw.wg.Wait()
	
	// Final flush
	if err := bw.Flush(); err != nil {
		return err
	}

	// Close underlying writer
	return bw.writer.Close()
}

// GetMetrics returns current batch writer metrics
func (bw *BatchWriter) GetMetrics() BatchMetrics {
	return BatchMetrics{
		TotalBatches:     atomic.LoadInt64(&bw.metrics.TotalBatches),
		TotalEntries:     atomic.LoadInt64(&bw.metrics.TotalEntries),
		FailedBatches:    atomic.LoadInt64(&bw.metrics.FailedBatches),
		RetryCount:       atomic.LoadInt64(&bw.metrics.RetryCount),
		LastFlushTime:    bw.metrics.LastFlushTime,
		AverageBatchSize: bw.metrics.AverageBatchSize,
	}
}

// SetBatchSize updates the batch size
func (bw *BatchWriter) SetBatchSize(size int) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	bw.config.BatchSize = size
}

// SetFlushInterval updates the flush interval
func (bw *BatchWriter) SetFlushInterval(interval time.Duration) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	bw.config.FlushInterval = interval
}

// GetBatchSize returns the current batch size
func (bw *BatchWriter) GetBatchSize() int {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	return len(bw.batch)
}

// IsClosed returns whether the writer is closed
func (bw *BatchWriter) IsClosed() bool {
	return atomic.LoadInt32(&bw.closed) == 1
}
