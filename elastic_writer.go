package spoor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ElasticWriter writes logs to Elasticsearch
type ElasticWriter struct {
	*BaseWriter
	mu         sync.RWMutex
	url        string
	index      string
	httpClient *http.Client
	bulkBuffer []ElasticBulkItem
}

// ElasticBulkItem represents a single item in Elasticsearch bulk API
type ElasticBulkItem struct {
	Index map[string]string `json:"index"`
	Data  LogEntry          `json:"-"`
}

// ElasticWriterConfig holds configuration for Elasticsearch writer
type ElasticWriterConfig struct {
	URL           string
	Index         string
	Username      string
	Password      string
	APIKey        string
	Formatter     Formatter
	BatchSize     int
	FlushInterval int // in seconds
	HTTPTimeout   int // in seconds
	RetryCount    int
	RetryDelay    int // in seconds
}

// NewElasticWriter creates a new Elasticsearch writer
func NewElasticWriter(config ElasticWriterConfig) *ElasticWriter {
	if config.Formatter == nil {
		config.Formatter = NewJSONFormatter()
	}

	// Set defaults
	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}
	if config.FlushInterval <= 0 {
		config.FlushInterval = 5
	}
	if config.HTTPTimeout <= 0 {
		config.HTTPTimeout = 30
	}
	if config.RetryCount <= 0 {
		config.RetryCount = 3
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = 1
	}

	baseWriter := NewBaseWriter(nil, config.Formatter)
	baseWriter.SetBatchSize(config.BatchSize)
	baseWriter.SetFlushInterval(time.Duration(config.FlushInterval) * time.Second)

	httpTimeout := time.Duration(config.HTTPTimeout) * time.Second

	writer := &ElasticWriter{
		BaseWriter: baseWriter,
		url:        strings.TrimSuffix(config.URL, "/"),
		index:      config.Index,
		httpClient: &http.Client{Timeout: httpTimeout},
		bulkBuffer: make([]ElasticBulkItem, 0, config.BatchSize),
	}

	// Start the flush loop
	writer.StartFlushLoop()

	return writer
}

// Write implements io.Writer interface
func (w *ElasticWriter) Write(p []byte) (n int, err error) {
	// For simple text output, create a log entry
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		Message:   string(p),
	}
	return len(p), w.WriteEntry(entry)
}

// WriteEntry writes a structured log entry
func (w *ElasticWriter) WriteEntry(entry LogEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Add to bulk buffer
	bulkItem := ElasticBulkItem{
		Index: map[string]string{
			"_index": w.index,
		},
		Data: entry,
	}

	w.bulkBuffer = append(w.bulkBuffer, bulkItem)

	// Flush if buffer is full
	if len(w.bulkBuffer) >= w.batchSize {
		return w.flushBulkUnsafe()
	}

	return nil
}

// WriteStructured writes a structured log entry
func (w *ElasticWriter) WriteStructured(entry LogEntry) error {
	return w.WriteEntry(entry)
}

// Flush flushes the bulk buffer to Elasticsearch
func (w *ElasticWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.flushBulkUnsafe()
}

// flushBulkUnsafe flushes the bulk buffer without locking
func (w *ElasticWriter) flushBulkUnsafe() error {
	if len(w.bulkBuffer) == 0 {
		return nil
	}

	// Prepare bulk request
	var bulkBody bytes.Buffer
	encoder := json.NewEncoder(&bulkBody)

	for _, item := range w.bulkBuffer {
		// Write index action
		if err := encoder.Encode(item.Index); err != nil {
			continue
		}

		// Write document
		if err := encoder.Encode(item.Data); err != nil {
			continue
		}
	}

	// Send to Elasticsearch with retry
	return w.sendBulkRequestWithRetry(&bulkBody)
}

// sendBulkRequestWithRetry sends bulk request with retry mechanism
func (w *ElasticWriter) sendBulkRequestWithRetry(bulkBody *bytes.Buffer) error {
	var lastErr error

	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("POST", w.url+"/_bulk", bulkBody)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		req.Header.Set("Content-Type", "application/x-ndjson")
		req.Header.Set("Accept", "application/json")

		resp, err := w.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to send to Elasticsearch: %w", err)
			if attempt < 2 {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return lastErr
		}
		defer resp.Body.Close()

		// Check response
		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("Elasticsearch error (status %d): %s", resp.StatusCode, string(body))
			if resp.StatusCode >= 500 && attempt < 2 {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return lastErr
		}

		// Success - clear buffer
		w.bulkBuffer = w.bulkBuffer[:0]
		return nil
	}

	return lastErr
}

// Close closes the Elasticsearch writer
func (w *ElasticWriter) Close() error {
	// Flush remaining data
	if err := w.Flush(); err != nil {
		return err
	}

	// Close base writer
	return w.BaseWriter.Close()
}

// SetIndex changes the Elasticsearch index
func (w *ElasticWriter) SetIndex(index string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.index = index
}

// GetIndex returns the current Elasticsearch index
func (w *ElasticWriter) GetIndex() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.index
}

// GetBulkBufferSize returns the current bulk buffer size
func (w *ElasticWriter) GetBulkBufferSize() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.bulkBuffer)
}

// NewElasticWriterWithDefaults creates an Elasticsearch writer with default settings
func NewElasticWriterWithDefaults(url, index string) *ElasticWriter {
	return NewElasticWriter(ElasticWriterConfig{
		URL:           url,
		Index:         index,
		BatchSize:     100,
		FlushInterval: 5,
		HTTPTimeout:   30,
		RetryCount:    3,
		RetryDelay:    1,
	})
}

// HealthCheck checks if Elasticsearch is accessible
func (w *ElasticWriter) HealthCheck() error {
	req, err := http.NewRequest("GET", w.url+"/_cluster/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check Elasticsearch health: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Elasticsearch health check failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetURL returns the Elasticsearch URL
func (w *ElasticWriter) GetURL() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.url
}

// SetURL changes the Elasticsearch URL
func (w *ElasticWriter) SetURL(url string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.url = strings.TrimSuffix(url, "/")
}
