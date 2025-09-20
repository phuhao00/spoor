package spoor

import (
	"fmt"
	"testing"
	"time"
)

func TestConsoleWriter(t *testing.T) {
	logger := NewConsole(LevelDebug)

	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	logger.Debugf("Formatted debug message: %s", "test")
	logger.Infof("Formatted info message: %d", 123)
	logger.Warnf("Formatted warning message: %v", time.Now())
	logger.Errorf("Formatted error message: %s", "error details")
}

func TestFileWriter(t *testing.T) {
	logger, err := NewFile("log", LevelDebug)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}

	logger.Debug("File debug message")
	logger.Info("File info message")
	logger.Warn("File warning message")
	logger.Error("File error message")

	// Give some time for the file writer to flush
	time.Sleep(100 * time.Millisecond)
}

func TestStructuredLogger(t *testing.T) {
	logger := NewConsole(LevelDebug)

	// Test with fields
	structuredLogger := logger.WithField("user_id", 123).WithField("action", "login")
	structuredLogger.Info("User logged in")

	structuredLogger = logger.WithFields(map[string]interface{}{
		"request_id": "req-123",
		"duration":   time.Millisecond * 150,
		"status":     200,
	})
	structuredLogger.Info("Request completed")

	// Test with error
	err := fmt.Errorf("database connection failed")
	structuredLogger = logger.WithError(err)
	structuredLogger.Error("Database operation failed")
}

func TestLevelFiltering(t *testing.T) {
	logger := NewConsole(LevelInfo)

	// These should not be printed (level too low)
	logger.Debug("This debug message should not appear")

	// These should be printed
	logger.Info("This info message should appear")
	logger.Warn("This warning message should appear")
	logger.Error("This error message should appear")
}

func TestElasticWriter(t *testing.T) {
	logger := NewElastic("http://localhost:9200", "test-logs", LevelDebug)

	logger.Info("Test log message for Elasticsearch")
	logger.Infof("Formatted message: %s", "elasticsearch test")

	// Give some time for the writer to flush
	time.Sleep(2 * time.Second)
}

func TestClickHouseWriter(t *testing.T) {
	// Note: This test requires a running ClickHouse instance
	// Skip if ClickHouse is not available
	logger, err := NewClickHouse("tcp://localhost:9000?database=test", "test_logs", LevelDebug)
	if err != nil {
		t.Skipf("ClickHouse not available: %v", err)
		return
	}

	logger.Info("Test log message for ClickHouse")
	logger.Infof("Formatted message: %s", "clickhouse test")

	// Give some time for the writer to flush
	time.Sleep(2 * time.Second)
}

func TestLogbusWriter(t *testing.T) {
	// Logbus writer not implemented yet
	t.Skip("Logbus writer not implemented yet")
}
