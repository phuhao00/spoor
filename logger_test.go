package spoor

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	// Test console logger
	writer := NewWriterFactory().CreateConsoleWriterToStdout()
	logger := New(writer, LevelDebug)

	if logger == nil {
		t.Fatal("Expected logger to be created")
	}

	// Test logging
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")
}

func TestNewWithDefaults(t *testing.T) {
	logger := NewWithDefaults()

	if logger == nil {
		t.Fatal("Expected default logger to be created")
	}

	logger.Info("Default logger test")
}

func TestNewConsole(t *testing.T) {
	logger := NewConsole(LevelInfo)

	if logger == nil {
		t.Fatal("Expected console logger to be created")
	}

	logger.Info("Console logger test")
}

func TestNewFile(t *testing.T) {
	logger, err := NewFile("test_logs", LevelInfo)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected file logger to be created")
	}

	logger.Info("File logger test")

	// Clean up
	os.RemoveAll("test_logs")
}

func TestNewElastic(t *testing.T) {
	logger := NewElastic("http://localhost:9200", "test-logs", LevelInfo)

	if logger == nil {
		t.Fatal("Expected elastic logger to be created")
	}

	logger.Info("Elastic logger test")
}

func TestNewClickHouse(t *testing.T) {
	logger, err := NewClickHouse("tcp://localhost:9000?database=test", "test_logs", LevelInfo)
	if err != nil {
		t.Skipf("ClickHouse not available: %v", err)
		return
	}

	if logger == nil {
		t.Fatal("Expected clickhouse logger to be created")
	}

	logger.Info("ClickHouse logger test")
}

func TestNewJSON(t *testing.T) {
	writer := NewWriterFactory().CreateConsoleWriterToStdout()
	logger := NewJSON(writer, LevelInfo)

	if logger == nil {
		t.Fatal("Expected JSON logger to be created")
	}

	logger.Info("JSON logger test")
}

func TestNewText(t *testing.T) {
	writer := NewWriterFactory().CreateConsoleWriterToStdout()
	logger := NewText(writer, LevelInfo)

	if logger == nil {
		t.Fatal("Expected text logger to be created")
	}

	logger.Info("Text logger test")
}

func TestNewLogrus(t *testing.T) {
	writer := NewWriterFactory().CreateConsoleWriterToStdout()
	logger := NewLogrus(writer, LevelInfo)

	if logger == nil {
		t.Fatal("Expected logrus logger to be created")
	}

	logger.Info("Logrus logger test")
}

func TestStructuredLogging(t *testing.T) {
	writer := NewWriterFactory().CreateConsoleWriterToStdout()
	logger := New(writer, LevelDebug)

	// Test WithField
	structuredLogger := logger.WithField("user_id", 123)
	structuredLogger.Info("User action")

	// Test WithFields
	structuredLogger = logger.WithFields(map[string]interface{}{
		"request_id": "req-123",
		"duration":   time.Millisecond * 150,
		"status":     200,
	})
	structuredLogger.Info("Request completed")

	// Test WithError
	err := fmt.Errorf("test error")
	structuredLogger = logger.WithError(err)
	structuredLogger.Error("Operation failed")
}

func TestLogLevels(t *testing.T) {
	writer := NewWriterFactory().CreateConsoleWriterToStdout()
	logger := New(writer, LevelInfo)

	// Test level filtering
	logger.Debug("This should not appear")
	logger.Info("This should appear")
	logger.Warn("This should appear")
	logger.Error("This should appear")

	// Change level
	logger.SetLevel(LevelDebug)
	logger.Debug("This should now appear")
}

func TestFormattedLogging(t *testing.T) {
	writer := NewWriterFactory().CreateConsoleWriterToStdout()
	logger := New(writer, LevelInfo)

	logger.Debugf("Debug message: %s", "test")
	logger.Infof("Info message: %d", 123)
	logger.Warnf("Warn message: %v", time.Now())
	logger.Errorf("Error message: %s", "error details")
}

func TestDefaultLogger(t *testing.T) {
	// Test default logger functions
	Debug("Default debug message")
	Info("Default info message")
	Warn("Default warn message")
	Error("Default error message")

	Debugf("Default debug: %s", "test")
	Infof("Default info: %d", 123)
	Warnf("Default warn: %v", time.Now())
	Errorf("Default error: %s", "error details")

	// Test structured logging with default logger
	WithField("key", "value").Info("Default structured message")
	WithFields(map[string]interface{}{
		"field1": "value1",
		"field2": 123,
	}).Info("Default structured message with fields")

	// Test error logging
	err := fmt.Errorf("test error")
	WithError(err).Error("Default error with context")
}

func TestConfig(t *testing.T) {
	// Test default config
	config := DefaultConfig()
	if config == nil {
		t.Fatal("Expected default config to be created")
	}

	// Test logger creation from config
	logger, err := NewFromConfig(config)
	if err != nil {
		t.Fatalf("Failed to create logger from config: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be created from config")
	}

	logger.Info("Config-based logger test")
}

func TestFormatters(t *testing.T) {
	writer := NewWriterFactory().CreateConsoleWriterToStdout()

	// Test text formatter
	textLogger := NewText(writer, LevelInfo)
	textLogger.Info("Text formatter test")

	// Test JSON formatter
	jsonLogger := NewJSON(writer, LevelInfo)
	jsonLogger.Info("JSON formatter test")

	// Test logrus formatter
	logrusLogger := NewLogrus(writer, LevelInfo)
	logrusLogger.Info("Logrus formatter test")
}

func TestConcurrentLogging(t *testing.T) {
	writer := NewWriterFactory().CreateConsoleWriterToStdout()
	logger := New(writer, LevelInfo)

	// Test concurrent logging
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				logger.Infof("Goroutine %d, message %d", id, j)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
