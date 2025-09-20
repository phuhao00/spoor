package spoor

import (
	"strings"
	"sync"
	"time"
)

// ObjectPool provides a generic object pool for memory reuse
type ObjectPool[T any] struct {
	pool sync.Pool
}

// NewObjectPool creates a new object pool
func NewObjectPool[T any](newFunc func() T) *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return newFunc()
			},
		},
	}
}

// Get gets an object from the pool
func (p *ObjectPool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put returns an object to the pool
func (p *ObjectPool[T]) Put(obj T) {
	p.pool.Put(obj)
}

// LogEntryPool provides a pool for LogEntry objects
var LogEntryPool = NewObjectPool(func() LogEntry {
	return LogEntry{
		Fields: make(map[string]interface{}),
	}
})

// GetLogEntry gets a LogEntry from the pool
func GetLogEntry() LogEntry {
	entry := LogEntryPool.Get()
	// Reset the entry
	entry.Timestamp = time.Time{}
	entry.Level = 0
	entry.Message = ""
	entry.Caller = ""
	// Clear fields map
	for k := range entry.Fields {
		delete(entry.Fields, k)
	}
	return entry
}

// PutLogEntry returns a LogEntry to the pool
func PutLogEntry(entry LogEntry) {
	LogEntryPool.Put(entry)
}

// BufferPool provides a pool for byte buffers
var BufferPool = NewObjectPool(func() []byte {
	return make([]byte, 0, 1024)
})

// GetBuffer gets a buffer from the pool
func GetBuffer() []byte {
	return BufferPool.Get()
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf []byte) {
	// Reset buffer length but keep capacity
	BufferPool.Put(buf[:0])
}

// StringBuilderPool provides a pool for string builders
var StringBuilderPool = NewObjectPool(func() *strings.Builder {
	return &strings.Builder{}
})

// GetStringBuilder gets a string builder from the pool
func GetStringBuilder() *strings.Builder {
	sb := StringBuilderPool.Get()
	sb.Reset()
	return sb
}

// PutStringBuilder returns a string builder to the pool
func PutStringBuilder(sb *strings.Builder) {
	StringBuilderPool.Put(sb)
}
