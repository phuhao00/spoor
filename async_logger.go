package spoor

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// AsyncLogger is a high-performance asynchronous logger
type AsyncLogger struct {
	*CoreLogger
	entryChan    chan LogEntry
	workerCount  int
	bufferSize   int
	flushTicker  *time.Ticker
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	closed       int32
	metrics      *LoggerMetrics
}

// LoggerMetrics tracks logger performance metrics
type LoggerMetrics struct {
	TotalLogs     int64
	DroppedLogs   int64
	BufferSize    int64
	FlushCount    int64
	ErrorCount    int64
	LastFlushTime time.Time
}

// AsyncLoggerConfig configures the async logger
type AsyncLoggerConfig struct {
	WorkerCount int           // Number of worker goroutines
	BufferSize  int           // Channel buffer size
	FlushInterval time.Duration // Auto-flush interval
	DropOnFull  bool          // Drop logs when buffer is full
}

// DefaultAsyncConfig returns default async logger configuration
func DefaultAsyncConfig() AsyncLoggerConfig {
	return AsyncLoggerConfig{
		WorkerCount:   runtime.NumCPU(),
		BufferSize:    10000,
		FlushInterval: 100 * time.Millisecond,
		DropOnFull:    true,
	}
}

// NewAsyncLogger creates a new high-performance async logger
func NewAsyncLogger(writer Writer, level LogLevel, config AsyncLoggerConfig, options ...Option) *AsyncLogger {
	ctx, cancel := context.WithCancel(context.Background())
	
	logger := &AsyncLogger{
		CoreLogger:   NewCoreLogger(writer, level, options...),
		entryChan:    make(chan LogEntry, config.BufferSize),
		workerCount:  config.WorkerCount,
		bufferSize:   config.BufferSize,
		flushTicker:  time.NewTicker(config.FlushInterval),
		ctx:          ctx,
		cancel:       cancel,
		metrics:      &LoggerMetrics{},
	}

	// Start workers
	for i := 0; i < config.WorkerCount; i++ {
		logger.wg.Add(1)
		go logger.worker(i)
	}

	// Start flush ticker
	go logger.flushLoop()

	return logger
}

// worker processes log entries from the channel
func (l *AsyncLogger) worker(id int) {
	defer l.wg.Done()
	
	batch := make([]LogEntry, 0, 100)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case entry, ok := <-l.entryChan:
			if !ok {
				// Channel closed, flush remaining batch
				if len(batch) > 0 {
					l.flushBatch(batch)
				}
				return
			}
			
			batch = append(batch, entry)
			atomic.AddInt64(&l.metrics.TotalLogs, 1)
			
			// Flush if batch is full
			if len(batch) >= 100 {
				l.flushBatch(batch)
				batch = batch[:0]
			}
			
		case <-ticker.C:
			// Periodic flush
			if len(batch) > 0 {
				l.flushBatch(batch)
				batch = batch[:0]
			}
			
		case <-l.ctx.Done():
			// Context cancelled, flush remaining batch
			if len(batch) > 0 {
				l.flushBatch(batch)
			}
			return
		}
	}
}

// flushBatch flushes a batch of log entries
func (l *AsyncLogger) flushBatch(batch []LogEntry) {
	if len(batch) == 0 {
		return
	}

	// Use structured writer if available
	if structuredWriter, ok := l.writer.(StructuredWriter); ok {
		for _, entry := range batch {
			if err := structuredWriter.WriteStructured(entry); err != nil {
				atomic.AddInt64(&l.metrics.ErrorCount, 1)
			}
		}
	} else {
		// Fallback to formatted output
		for _, entry := range batch {
			if data, err := l.formatter.Format(entry); err == nil {
				if _, err := l.writer.Write(data); err != nil {
					atomic.AddInt64(&l.metrics.ErrorCount, 1)
				}
			}
		}
	}

	atomic.AddInt64(&l.metrics.FlushCount, 1)
	l.metrics.LastFlushTime = time.Now()
}

// flushLoop handles periodic flushing
func (l *AsyncLogger) flushLoop() {
	for {
		select {
		case <-l.flushTicker.C:
			l.Sync()
		case <-l.ctx.Done():
			return
		}
	}
}

// log sends a log entry to the async channel
func (l *AsyncLogger) log(level LogLevel, msg string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	if atomic.LoadInt32(&l.closed) == 1 {
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

	// Send to channel (non-blocking)
	select {
	case l.entryChan <- entry:
		atomic.AddInt64(&l.metrics.BufferSize, 1)
	default:
		// Channel is full
		if l.metrics != nil {
			atomic.AddInt64(&l.metrics.DroppedLogs, 1)
		}
	}
}

// Sync flushes all buffered log entries
func (l *AsyncLogger) Sync() error {
	// Force flush by sending a special entry
	select {
	case l.entryChan <- LogEntry{Level: -1}: // Special marker
	default:
	}
	return l.writer.Flush()
}

// Close closes the async logger
func (l *AsyncLogger) Close() error {
	if !atomic.CompareAndSwapInt32(&l.closed, 0, 1) {
		return nil
	}

	// Stop the flush ticker
	l.flushTicker.Stop()
	
	// Cancel context to stop workers
	l.cancel()
	
	// Close the channel
	close(l.entryChan)
	
	// Wait for workers to finish
	l.wg.Wait()
	
	// Close the underlying writer
	return l.writer.Close()
}

// GetMetrics returns current logger metrics
func (l *AsyncLogger) GetMetrics() LoggerMetrics {
	return LoggerMetrics{
		TotalLogs:     atomic.LoadInt64(&l.metrics.TotalLogs),
		DroppedLogs:   atomic.LoadInt64(&l.metrics.DroppedLogs),
		BufferSize:    int64(len(l.entryChan)),
		FlushCount:    atomic.LoadInt64(&l.metrics.FlushCount),
		ErrorCount:    atomic.LoadInt64(&l.metrics.ErrorCount),
		LastFlushTime: l.metrics.LastFlushTime,
	}
}

// ResetMetrics resets the logger metrics
func (l *AsyncLogger) ResetMetrics() {
	atomic.StoreInt64(&l.metrics.TotalLogs, 0)
	atomic.StoreInt64(&l.metrics.DroppedLogs, 0)
	atomic.StoreInt64(&l.metrics.FlushCount, 0)
	atomic.StoreInt64(&l.metrics.ErrorCount, 0)
	l.metrics.LastFlushTime = time.Time{}
}
