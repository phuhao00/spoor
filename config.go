package spoor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config represents the main configuration structure
type Config struct {
	Loggers map[string]LoggerConfig `json:"loggers"`
	Default string                  `json:"default"`
}

// LoggerConfig represents configuration for a specific logger
type LoggerConfig struct {
	Type       string                 `json:"type"`
	Level      string                 `json:"level"`
	Output     string                 `json:"output"`
	Format     string                 `json:"format"`
	Async      bool                   `json:"async"`
	BatchSize  int                    `json:"batch_size"`
	FlushEvery string                 `json:"flush_every"`
	FilePath   string                 `json:"file_path,omitempty"`
	MaxSize    int                    `json:"max_size,omitempty"`
	MaxBackups int                    `json:"max_backups,omitempty"`
	MaxAge     int                    `json:"max_age,omitempty"`
	Compress   bool                   `json:"compress,omitempty"`
	Elastic    *ElasticConfig         `json:"elastic,omitempty"`
	ClickHouse *ClickHouseConfig      `json:"clickhouse,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	Sampling   *SamplingConfig        `json:"sampling,omitempty"`
	Filtering  *FilteringConfig       `json:"filtering,omitempty"`
}

// ElasticConfig represents Elasticsearch configuration
type ElasticConfig struct {
	URL       string `json:"url"`
	Index     string `json:"index"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	BatchSize int    `json:"batch_size"`
	FlushTime string `json:"flush_time"`
}

// ClickHouseConfig represents ClickHouse configuration
type ClickHouseConfig struct {
	DSN      string `json:"dsn"`
	Table    string `json:"table"`
	Database string `json:"database,omitempty"`
}

// SamplingConfig represents sampling configuration
type SamplingConfig struct {
	Type  string             `json:"type"` // "rate", "level"
	Rate  float64            `json:"rate,omitempty"`
	Level map[string]float64 `json:"level,omitempty"`
}

// FilteringConfig represents filtering configuration
type FilteringConfig struct {
	MinLevel string            `json:"min_level,omitempty"`
	Fields   map[string]string `json:"fields,omitempty"`
}

// LoadConfig loads configuration from a file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(config *Config, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// CreateLoggerFromConfig creates a logger from configuration
func CreateLoggerFromConfig(config *LoggerConfig) (Logger, error) {
	// Parse log level
	level, err := ParseLogLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// Create writer based on type
	var writer Writer
	switch config.Type {
	case "console":
		writer = NewConsoleWriter(ConsoleWriterConfig{
			Output: os.Stdout,
		})
	case "file":
		fileConfig := FileWriterConfig{
			LogDir: config.FilePath,
		}
		if config.MaxSize > 0 {
			fileConfig.MaxSize = int64(config.MaxSize)
		}
		
		var err error
		writer, err = NewFileWriter(fileConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create file writer: %w", err)
		}
	case "elastic":
		if config.Elastic == nil {
			return nil, fmt.Errorf("elastic configuration is required")
		}
		flushTime, _ := time.ParseDuration(config.Elastic.FlushTime)
		elasticConfig := ElasticWriterConfig{
			URL:       config.Elastic.URL,
			Index:     config.Elastic.Index,
			Username:  config.Elastic.Username,
			Password:  config.Elastic.Password,
			BatchSize: config.Elastic.BatchSize,
			FlushInterval: int(flushTime.Milliseconds()),
		}
		writer = NewElasticWriter(elasticConfig)
	case "clickhouse":
		if config.ClickHouse == nil {
			return nil, fmt.Errorf("clickhouse configuration is required")
		}
		clickhouseConfig := ClickHouseWriterConfig{
			DSN:      config.ClickHouse.DSN,
			TableName: config.ClickHouse.Table,
		}
		var err error
		writer, err = NewClickHouseWriter(clickhouseConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create clickhouse writer: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported logger type: %s", config.Type)
	}

	// Create formatter
	var formatter Formatter
	switch config.Format {
	case "json":
		formatter = NewJSONFormatter()
	case "text":
		formatter = NewTextFormatter()
	case "logrus":
		formatter = NewLogrusFormatter()
	default:
		formatter = NewTextFormatter()
	}

	// Create logger options
	var options []Option
	options = append(options, WithFormatter(formatter))

	// Add fields if specified
	if len(config.Fields) > 0 {
		// This would need to be implemented in the logger
	}

	// Create logger based on async setting
	if config.Async {
		asyncConfig := DefaultAsyncConfig()
		if config.BatchSize > 0 {
			asyncConfig.BufferSize = config.BatchSize
		}
		if config.FlushEvery != "" {
			if flushTime, err := time.ParseDuration(config.FlushEvery); err == nil {
				asyncConfig.FlushInterval = flushTime
			}
		}
		return NewAsyncLogger(writer, level, asyncConfig, options...), nil
	}

	return NewCoreLogger(writer, level, options...), nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Loggers: map[string]LoggerConfig{
			"default": {
				Type:       "console",
				Level:      "info",
				Format:     "text",
				Async:      false,
				BatchSize:  1000,
				FlushEvery: "100ms",
			},
			"file": {
				Type:       "file",
				Level:      "info",
				Format:     "json",
				Async:      true,
				BatchSize:  1000,
				FlushEvery: "1s",
				FilePath:   "logs/app.log",
				MaxSize:    100 * 1024 * 1024, // 100MB
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   true,
			},
			"async": {
				Type:       "console",
				Level:      "debug",
				Format:     "json",
				Async:      true,
				BatchSize:  10000,
				FlushEvery: "100ms",
			},
		},
		Default: "default",
	}
}

// CreateConfigFile creates a sample configuration file
func CreateConfigFile(filename string) error {
	config := DefaultConfig()
	return SaveConfig(config, filename)
}