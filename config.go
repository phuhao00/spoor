package spoor

import (
	"fmt"
	"os"
	"time"
)

// Config represents the main configuration for the logger
type Config struct {
	Level     LogLevel        `json:"level" yaml:"level"`
	Writer    WriterConfig    `json:"writer" yaml:"writer"`
	Formatter FormatterConfig `json:"formatter" yaml:"formatter"`
	Hooks     []HookConfig    `json:"hooks" yaml:"hooks"`
	Caller    bool            `json:"caller" yaml:"caller"`
}

// WriterConfig represents the writer configuration
type WriterConfig struct {
	Type   string                 `json:"type" yaml:"type"`
	Config map[string]interface{} `json:"config" yaml:"config"`
}

// FormatterConfig represents the formatter configuration
type FormatterConfig struct {
	Type   string                 `json:"type" yaml:"type"`
	Config map[string]interface{} `json:"config" yaml:"config"`
}

// HookConfig represents the hook configuration
type HookConfig struct {
	Type   string                 `json:"type" yaml:"type"`
	Config map[string]interface{} `json:"config" yaml:"config"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Level: LevelInfo,
		Writer: WriterConfig{
			Type: "console",
			Config: map[string]interface{}{
				"output": "stdout",
			},
		},
		Formatter: FormatterConfig{
			Type: "text",
			Config: map[string]interface{}{
				"timestamp_format": time.RFC3339,
				"disable_colors":   false,
				"disable_caller":   false,
			},
		},
		Hooks:  []HookConfig{},
		Caller: true,
	}
}

// NewFromConfig creates a new logger from configuration
func NewFromConfig(config *Config) (Logger, error) {
	// Create writer
	writer, err := createWriterFromConfig(config.Writer)
	if err != nil {
		return nil, fmt.Errorf("failed to create writer: %w", err)
	}

	// Create formatter
	formatter, err := createFormatterFromConfig(config.Formatter)
	if err != nil {
		return nil, fmt.Errorf("failed to create formatter: %w", err)
	}

	// Create hooks
	hooks, err := createHooksFromConfig(config.Hooks)
	if err != nil {
		return nil, fmt.Errorf("failed to create hooks: %w", err)
	}

	// Create logger
	options := []Option{
		WithFormatter(formatter),
		WithCaller(config.Caller),
	}

	if len(hooks) > 0 {
		options = append(options, WithHooks(hooks...))
	}

	return NewCoreLogger(writer, config.Level, options...), nil
}

// createWriterFromConfig creates a writer from configuration
func createWriterFromConfig(config WriterConfig) (Writer, error) {
	factory := NewWriterFactory()

	switch config.Type {
	case "console":
		cfg := ConsoleWriterConfig{}
		if output, ok := config.Config["output"].(string); ok {
			switch output {
			case "stdout":
				cfg.Output = os.Stdout
			case "stderr":
				cfg.Output = os.Stderr
			}
		}
		if batchSize, ok := config.Config["batch_size"].(int); ok {
			cfg.BatchSize = batchSize
		}
		if flushInterval, ok := config.Config["flush_interval"].(int); ok {
			cfg.FlushInterval = flushInterval
		}
		return factory.CreateWriter(WriterTypeConsole, cfg)

	case "file":
		cfg := FileWriterConfig{}
		if logDir, ok := config.Config["log_dir"].(string); ok {
			cfg.LogDir = logDir
		}
		if maxSize, ok := config.Config["max_size"].(int64); ok {
			cfg.MaxSize = maxSize
		}
		if batchSize, ok := config.Config["batch_size"].(int); ok {
			cfg.BatchSize = batchSize
		}
		if flushInterval, ok := config.Config["flush_interval"].(int); ok {
			cfg.FlushInterval = flushInterval
		}
		return factory.CreateWriter(WriterTypeFile, cfg)

	case "elastic":
		cfg := ElasticWriterConfig{}
		if url, ok := config.Config["url"].(string); ok {
			cfg.URL = url
		}
		if index, ok := config.Config["index"].(string); ok {
			cfg.Index = index
		}
		if batchSize, ok := config.Config["batch_size"].(int); ok {
			cfg.BatchSize = batchSize
		}
		if flushInterval, ok := config.Config["flush_interval"].(int); ok {
			cfg.FlushInterval = flushInterval
		}
		if httpTimeout, ok := config.Config["http_timeout"].(int); ok {
			cfg.HTTPTimeout = httpTimeout
		}
		return factory.CreateWriter(WriterTypeElastic, cfg)

	case "clickhouse":
		cfg := ClickHouseWriterConfig{}
		if dsn, ok := config.Config["dsn"].(string); ok {
			cfg.DSN = dsn
		}
		if tableName, ok := config.Config["table_name"].(string); ok {
			cfg.TableName = tableName
		}
		if batchSize, ok := config.Config["batch_size"].(int); ok {
			cfg.BatchSize = batchSize
		}
		if flushTime, ok := config.Config["flush_time"].(int); ok {
			cfg.FlushTime = flushTime
		}
		if httpTimeout, ok := config.Config["http_timeout"].(int); ok {
			cfg.HTTPTimeout = httpTimeout
		}
		return factory.CreateWriter(WriterTypeClickHouse, cfg)

	default:
		return nil, fmt.Errorf("unsupported writer type: %s", config.Type)
	}
}

// createFormatterFromConfig creates a formatter from configuration
func createFormatterFromConfig(config FormatterConfig) (Formatter, error) {
	switch config.Type {
	case "text":
		formatter := NewTextFormatter()
		if timestampFormat, ok := config.Config["timestamp_format"].(string); ok {
			formatter.TimestampFormat = timestampFormat
		}
		if disableColors, ok := config.Config["disable_colors"].(bool); ok {
			formatter.DisableColors = disableColors
		}
		if disableCaller, ok := config.Config["disable_caller"].(bool); ok {
			formatter.DisableCaller = disableCaller
		}
		return formatter, nil

	case "json":
		formatter := NewJSONFormatter()
		if timestampFormat, ok := config.Config["timestamp_format"].(string); ok {
			formatter.TimestampFormat = timestampFormat
		}
		if prettyPrint, ok := config.Config["pretty_print"].(bool); ok {
			formatter.PrettyPrint = prettyPrint
		}
		return formatter, nil

	case "logrus":
		formatter := NewLogrusFormatter()
		if timestampFormat, ok := config.Config["timestamp_format"].(string); ok {
			formatter.TimestampFormat = timestampFormat
		}
		if disableCaller, ok := config.Config["disable_caller"].(bool); ok {
			formatter.DisableCaller = disableCaller
		}
		return formatter, nil

	default:
		return nil, fmt.Errorf("unsupported formatter type: %s", config.Type)
	}
}

// createHooksFromConfig creates hooks from configuration
func createHooksFromConfig(configs []HookConfig) ([]Hook, error) {
	hooks := make([]Hook, 0, len(configs))

	for _, config := range configs {
		// TODO: Implement hook creation based on type
		// For now, we'll skip hooks
		_ = config
	}

	return hooks, nil
}

// LoadConfigFromFile loads configuration from a file
func LoadConfigFromFile(filename string) (*Config, error) {
	// TODO: Implement file loading (YAML/JSON)
	// For now, return default config
	return DefaultConfig(), nil
}

// SaveConfigToFile saves configuration to a file
func SaveConfigToFile(config *Config, filename string) error {
	// TODO: Implement file saving (YAML/JSON)
	return fmt.Errorf("config saving not implemented yet")
}
