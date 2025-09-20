package spoor

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Sampler defines the interface for log sampling
type Sampler interface {
	ShouldSample(entry LogEntry) bool
}

// RateSampler samples logs at a given rate
type RateSampler struct {
	rate    float64 // Sampling rate (0.0 to 1.0)
	rand    *rand.Rand
	mu      sync.Mutex
	counter int64
}

// NewRateSampler creates a new rate sampler
func NewRateSampler(rate float64) *RateSampler {
	return &RateSampler{
		rate: rate,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ShouldSample determines if a log entry should be sampled
func (rs *RateSampler) ShouldSample(entry LogEntry) bool {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	
	atomic.AddInt64(&rs.counter, 1)
	return rs.rand.Float64() < rs.rate
}

// LevelSampler samples logs based on level
type LevelSampler struct {
	levelRates map[LogLevel]float64
	rand       *rand.Rand
	mu         sync.Mutex
}

// NewLevelSampler creates a new level-based sampler
func NewLevelSampler(levelRates map[LogLevel]float64) *LevelSampler {
	return &LevelSampler{
		levelRates: levelRates,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ShouldSample determines if a log entry should be sampled based on its level
func (ls *LevelSampler) ShouldSample(entry LogEntry) bool {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	
	rate, exists := ls.levelRates[entry.Level]
	if !exists {
		return true // Default to sampling if level not specified
	}
	
	return ls.rand.Float64() < rate
}

// Filter defines the interface for log filtering
type Filter interface {
	ShouldLog(entry LogEntry) bool
}

// LevelFilter filters logs based on level
type LevelFilter struct {
	minLevel LogLevel
}

// NewLevelFilter creates a new level filter
func NewLevelFilter(minLevel LogLevel) *LevelFilter {
	return &LevelFilter{minLevel: minLevel}
}

// ShouldLog determines if a log entry should be logged
func (lf *LevelFilter) ShouldLog(entry LogEntry) bool {
	return entry.Level >= lf.minLevel
}

// FieldFilter filters logs based on field values
type FieldFilter struct {
	field    string
	value    interface{}
	operator string // "eq", "ne", "gt", "lt", "contains"
}

// NewFieldFilter creates a new field filter
func NewFieldFilter(field string, value interface{}, operator string) *FieldFilter {
	return &FieldFilter{
		field:    field,
		value:    value,
		operator: operator,
	}
}

// ShouldLog determines if a log entry should be logged based on field values
func (ff *FieldFilter) ShouldLog(entry LogEntry) bool {
	fieldValue, exists := entry.Fields[ff.field]
	if !exists {
		return false
	}

	switch ff.operator {
	case "eq":
		return fieldValue == ff.value
	case "ne":
		return fieldValue != ff.value
	case "contains":
		if str, ok := fieldValue.(string); ok {
			if targetStr, ok := ff.value.(string); ok {
				return contains(str, targetStr)
			}
		}
		return false
	default:
		return true
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}

// CompositeFilter combines multiple filters
type CompositeFilter struct {
	filters []Filter
	mode    string // "and", "or"
}

// NewCompositeFilter creates a new composite filter
func NewCompositeFilter(filters []Filter, mode string) *CompositeFilter {
	return &CompositeFilter{
		filters: filters,
		mode:    mode,
	}
}

// ShouldLog determines if a log entry should be logged based on all filters
func (cf *CompositeFilter) ShouldLog(entry LogEntry) bool {
	if len(cf.filters) == 0 {
		return true
	}

	if cf.mode == "and" {
		for _, filter := range cf.filters {
			if !filter.ShouldLog(entry) {
				return false
			}
		}
		return true
	} else { // "or"
		for _, filter := range cf.filters {
			if filter.ShouldLog(entry) {
				return true
			}
		}
		return false
	}
}

// MetricsCollector collects logging metrics
type MetricsCollector struct {
	totalLogs     int64
	logsByLevel   map[LogLevel]int64
	droppedLogs   int64
	errorCount    int64
	lastLogTime   time.Time
	mu            sync.RWMutex
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		logsByLevel: make(map[LogLevel]int64),
	}
}

// RecordLog records a log entry
func (mc *MetricsCollector) RecordLog(entry LogEntry) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	atomic.AddInt64(&mc.totalLogs, 1)
	mc.logsByLevel[entry.Level]++
	mc.lastLogTime = time.Now()
}

// RecordDropped records a dropped log
func (mc *MetricsCollector) RecordDropped() {
	atomic.AddInt64(&mc.droppedLogs, 1)
}

// RecordError records an error
func (mc *MetricsCollector) RecordError() {
	atomic.AddInt64(&mc.errorCount, 1)
}

// GetMetrics returns current metrics
func (mc *MetricsCollector) GetMetrics() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	metrics := map[string]interface{}{
		"total_logs":    atomic.LoadInt64(&mc.totalLogs),
		"dropped_logs":  atomic.LoadInt64(&mc.droppedLogs),
		"error_count":   atomic.LoadInt64(&mc.errorCount),
		"last_log_time": mc.lastLogTime,
		"logs_by_level": make(map[string]int64),
	}
	
	// Convert logs by level to string keys
	logsByLevel := make(map[string]int64)
	for level, count := range mc.logsByLevel {
		logsByLevel[level.String()] = count
	}
	metrics["logs_by_level"] = logsByLevel
	
	return metrics
}

// Reset resets all metrics
func (mc *MetricsCollector) Reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	atomic.StoreInt64(&mc.totalLogs, 0)
	atomic.StoreInt64(&mc.droppedLogs, 0)
	atomic.StoreInt64(&mc.errorCount, 0)
	
	for level := range mc.logsByLevel {
		mc.logsByLevel[level] = 0
	}
	
	mc.lastLogTime = time.Time{}
}

// AdvancedLogger provides advanced logging features
type AdvancedLogger struct {
	*CoreLogger
	sampler  Sampler
	filter   Filter
	metrics  *MetricsCollector
}

// AdvancedConfig configures the advanced logger
type AdvancedConfig struct {
	Sampler Sampler
	Filter  Filter
	Metrics bool
}

// NewAdvancedLogger creates a new advanced logger
func NewAdvancedLogger(writer Writer, level LogLevel, config AdvancedConfig, options ...Option) *AdvancedLogger {
	coreLogger := NewCoreLogger(writer, level, options...)
	
	advancedLogger := &AdvancedLogger{
		CoreLogger: coreLogger,
		sampler:    config.Sampler,
		filter:     config.Filter,
	}
	
	if config.Metrics {
		advancedLogger.metrics = NewMetricsCollector()
	}
	
	return advancedLogger
}

// log overrides the core logger's log method to add advanced features
func (al *AdvancedLogger) log(level LogLevel, msg string, fields map[string]interface{}) {
	if level < al.level {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    make(map[string]interface{}),
	}

	// Add logger fields
	al.mu.RLock()
	for k, v := range al.fields {
		entry.Fields[k] = v
	}
	al.mu.RUnlock()

	// Add method fields
	for k, v := range fields {
		entry.Fields[k] = v
	}

	// Add caller information if enabled
	if al.caller {
		if caller := getCaller(); caller != "" {
			entry.Caller = caller
		}
	}

	// Apply filter
	if al.filter != nil && !al.filter.ShouldLog(entry) {
		if al.metrics != nil {
			al.metrics.RecordDropped()
		}
		return
	}

	// Apply sampler
	if al.sampler != nil && !al.sampler.ShouldSample(entry) {
		if al.metrics != nil {
			al.metrics.RecordDropped()
		}
		return
	}

	// Record metrics
	if al.metrics != nil {
		al.metrics.RecordLog(entry)
	}

	// Fire hooks
	for _, hook := range al.hooks {
		if al.shouldFireHook(hook, level) {
			if err := hook.Fire(entry); err != nil && al.metrics != nil {
				al.metrics.RecordError()
			}
		}
	}

	// Write the log entry
	if structuredWriter, ok := al.writer.(StructuredWriter); ok {
		if err := structuredWriter.WriteStructured(entry); err != nil && al.metrics != nil {
			al.metrics.RecordError()
		}
	} else {
		// Fallback to text format
		if data, err := al.formatter.Format(entry); err == nil {
			if _, err := al.writer.Write(data); err != nil && al.metrics != nil {
				al.metrics.RecordError()
			}
		} else if al.metrics != nil {
			al.metrics.RecordError()
		}
	}
}

// GetMetrics returns current metrics
func (al *AdvancedLogger) GetMetrics() map[string]interface{} {
	if al.metrics == nil {
		return nil
	}
	return al.metrics.GetMetrics()
}

// ResetMetrics resets all metrics
func (al *AdvancedLogger) ResetMetrics() {
	if al.metrics != nil {
		al.metrics.Reset()
	}
}

// SetSampler sets the sampler
func (al *AdvancedLogger) SetSampler(sampler Sampler) {
	al.mu.Lock()
	defer al.mu.Unlock()
	al.sampler = sampler
}

// SetFilter sets the filter
func (al *AdvancedLogger) SetFilter(filter Filter) {
	al.mu.Lock()
	defer al.mu.Unlock()
	al.filter = filter
}
