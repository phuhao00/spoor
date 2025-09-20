package spoor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

// ClickHouseWriter writes logs to ClickHouse
type ClickHouseWriter struct {
	*BaseWriter
	mu        sync.RWMutex
	db        *sql.DB
	tableName string
	batchSize int
	flushTime time.Duration
	stopChan  chan struct{}
}

// ClickHouseWriterConfig holds configuration for ClickHouse writer
type ClickHouseWriterConfig struct {
	DSN         string
	TableName   string
	Formatter   Formatter
	BatchSize   int
	FlushTime   int // in seconds
	HTTPTimeout int // in seconds
}

// NewClickHouseWriter creates a new ClickHouse writer
func NewClickHouseWriter(config ClickHouseWriterConfig) (*ClickHouseWriter, error) {
	if config.Formatter == nil {
		config.Formatter = NewJSONFormatter()
	}

	// Connect to ClickHouse
	db, err := sql.Open("clickhouse", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	// Create table if not exists
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			timestamp DateTime,
			level String,
			message String,
			fields String,
			caller String
		) ENGINE = MergeTree()
		ORDER BY timestamp
	`, config.TableName)

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	baseWriter := NewBaseWriter(nil, config.Formatter)
	if config.BatchSize > 0 {
		baseWriter.SetBatchSize(config.BatchSize)
	}
	if config.FlushTime > 0 {
		baseWriter.SetFlushInterval(time.Duration(config.FlushTime) * time.Second)
	}

	writer := &ClickHouseWriter{
		BaseWriter: baseWriter,
		db:         db,
		tableName:  config.TableName,
		batchSize:  config.BatchSize,
		flushTime:  time.Duration(config.FlushTime) * time.Second,
		stopChan:   make(chan struct{}),
	}

	// Start the flush loop
	writer.StartFlushLoop()

	return writer, nil
}

// NewClickHouseWriterWithDefaults creates a ClickHouse writer with default settings
func NewClickHouseWriterWithDefaults(dsn, tableName string) (*ClickHouseWriter, error) {
	return NewClickHouseWriter(ClickHouseWriterConfig{
		DSN:       dsn,
		TableName: tableName,
		BatchSize: 100,
		FlushTime: 5,
	})
}

// Write implements io.Writer interface
func (w *ClickHouseWriter) Write(p []byte) (n int, err error) {
	// For simple text output, create a log entry
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		Message:   string(p),
	}
	return len(p), w.WriteEntry(entry)
}

// WriteEntry writes a structured log entry
func (w *ClickHouseWriter) WriteEntry(entry LogEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Add to buffer
	w.buffer = append(w.buffer, entry)

	// Flush if buffer is full
	if len(w.buffer) >= w.batchSize {
		return w.flushToClickHouseUnsafe()
	}

	return nil
}

// WriteStructured writes a structured log entry
func (w *ClickHouseWriter) WriteStructured(entry LogEntry) error {
	return w.WriteEntry(entry)
}

// Flush flushes the buffer to ClickHouse
func (w *ClickHouseWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.flushToClickHouseUnsafe()
}

// flushToClickHouseUnsafe flushes the buffer to ClickHouse without locking
func (w *ClickHouseWriter) flushToClickHouseUnsafe() error {
	if len(w.buffer) == 0 {
		return nil
	}

	// Prepare batch insert
	stmt, err := w.db.Prepare(fmt.Sprintf(`
		INSERT INTO %s (timestamp, level, message, fields, caller) 
		VALUES (?, ?, ?, ?, ?)
	`, w.tableName))
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert each entry
	for _, entry := range w.buffer {
		// Serialize fields to JSON
		fieldsJSON := ""
		if len(entry.Fields) > 0 {
			if data, err := json.Marshal(entry.Fields); err == nil {
				fieldsJSON = string(data)
			}
		}

		_, err := stmt.Exec(
			entry.Timestamp,
			entry.Level.String(),
			entry.Message,
			fieldsJSON,
			entry.Caller,
		)
		if err != nil {
			// Log error but continue with other entries
			fmt.Printf("ClickHouse insert error: %v\n", err)
			continue
		}
	}

	// Clear buffer
	w.buffer = w.buffer[:0]
	return nil
}

// Close closes the ClickHouse writer
func (w *ClickHouseWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Flush remaining data
	if err := w.flushToClickHouseUnsafe(); err != nil {
		return err
	}

	// Close database connection
	if w.db != nil {
		return w.db.Close()
	}

	return nil
}

// SetTableName changes the ClickHouse table name
func (w *ClickHouseWriter) SetTableName(tableName string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.tableName = tableName
}

// GetTableName returns the current ClickHouse table name
func (w *ClickHouseWriter) GetTableName() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.tableName
}

// GetBufferSize returns the current buffer size
func (w *ClickHouseWriter) GetBufferSize() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.buffer)
}

// IsConnected checks if the ClickHouse connection is alive
func (w *ClickHouseWriter) IsConnected() bool {
	if w.db == nil {
		return false
	}
	return w.db.Ping() == nil
}
