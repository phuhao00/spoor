package spoor

import (
	"fmt"
	"os"
)

// WriterType represents the type of writer
type WriterType string

const (
	WriterTypeConsole    WriterType = "console"
	WriterTypeFile       WriterType = "file"
	WriterTypeElastic    WriterType = "elastic"
	WriterTypeClickHouse WriterType = "clickhouse"
	WriterTypeLogbus     WriterType = "logbus"
)

// WriterFactory creates writers based on configuration
type WriterFactory struct{}

// NewWriterFactory creates a new writer factory
func NewWriterFactory() *WriterFactory {
	return &WriterFactory{}
}

// CreateWriter creates a writer based on the type and configuration
func (f *WriterFactory) CreateWriter(writerType WriterType, config interface{}) (Writer, error) {
	switch writerType {
	case WriterTypeConsole:
		return f.createConsoleWriter(config)
	case WriterTypeFile:
		return f.createFileWriter(config)
	case WriterTypeElastic:
		return f.createElasticWriter(config)
	case WriterTypeClickHouse:
		return f.createClickHouseWriter(config)
	case WriterTypeLogbus:
		return f.createLogbusWriter(config)
	default:
		return nil, fmt.Errorf("unsupported writer type: %s", writerType)
	}
}

// createConsoleWriter creates a console writer
func (f *WriterFactory) createConsoleWriter(config interface{}) (Writer, error) {
	if config == nil {
		return NewConsoleWriterWithDefaults(), nil
	}

	cfg, ok := config.(ConsoleWriterConfig)
	if !ok {
		return nil, fmt.Errorf("invalid console writer config")
	}

	return NewConsoleWriter(cfg), nil
}

// createFileWriter creates a file writer
func (f *WriterFactory) createFileWriter(config interface{}) (Writer, error) {
	if config == nil {
		return nil, fmt.Errorf("file writer requires configuration")
	}

	cfg, ok := config.(FileWriterConfig)
	if !ok {
		return nil, fmt.Errorf("invalid file writer config")
	}

	return NewFileWriter(cfg)
}

// createElasticWriter creates an Elasticsearch writer
func (f *WriterFactory) createElasticWriter(config interface{}) (Writer, error) {
	if config == nil {
		return nil, fmt.Errorf("elastic writer requires configuration")
	}

	cfg, ok := config.(ElasticWriterConfig)
	if !ok {
		return nil, fmt.Errorf("invalid elastic writer config")
	}

	return NewElasticWriter(cfg), nil
}

// createClickHouseWriter creates a ClickHouse writer
func (f *WriterFactory) createClickHouseWriter(config interface{}) (Writer, error) {
	if config == nil {
		return nil, fmt.Errorf("clickhouse writer requires configuration")
	}

	cfg, ok := config.(ClickHouseWriterConfig)
	if !ok {
		return nil, fmt.Errorf("invalid clickhouse writer config")
	}

	return NewClickHouseWriter(cfg)
}

// createLogbusWriter creates a Logbus writer
func (f *WriterFactory) createLogbusWriter(config interface{}) (Writer, error) {
	// TODO: Implement Logbus writer
	return nil, fmt.Errorf("Logbus writer not implemented yet")
}

// CreateConsoleWriterToStdout creates a console writer that writes to stdout
func (f *WriterFactory) CreateConsoleWriterToStdout() Writer {
	return NewConsoleWriter(ConsoleWriterConfig{
		Output: os.Stdout,
	})
}

// CreateConsoleWriterToStderr creates a console writer that writes to stderr
func (f *WriterFactory) CreateConsoleWriterToStderr() Writer {
	return NewConsoleWriter(ConsoleWriterConfig{
		Output: os.Stderr,
	})
}

// CreateFileWriterWithDefaults creates a file writer with default settings
func (f *WriterFactory) CreateFileWriterWithDefaults(logDir string) (Writer, error) {
	return NewFileWriterWithDefaults(logDir)
}

// CreateElasticWriterWithDefaults creates an Elasticsearch writer with default settings
func (f *WriterFactory) CreateElasticWriterWithDefaults(url, index string) Writer {
	return NewElasticWriterWithDefaults(url, index)
}

// CreateClickHouseWriterWithDefaults creates a ClickHouse writer with default settings
func (f *WriterFactory) CreateClickHouseWriterWithDefaults(dsn, tableName string) (Writer, error) {
	return NewClickHouseWriterWithDefaults(dsn, tableName)
}
