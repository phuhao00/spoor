package main

import (
	"fmt"
	"time"

	"github.com/phuhao00/spoor"
)

func main() {
	fmt.Println("📝 Spoor Simple Usage Examples")
	fmt.Println("===============================")

	// Example 1: Quick Start
	fmt.Println("\n1. Quick Start")
	logger := spoor.Quick()
	logger.Info("Hello, Spoor!")
	logger.Debug("This is a debug message")
	logger.Warn("This is a warning")
	logger.Error("This is an error")
	logger.Close()

	// Example 2: File Logging
	fmt.Println("\n2. File Logging")
	fileLogger, err := spoor.QuickFile("logs/app.log")
	if err != nil {
		fmt.Printf("Error creating file logger: %v\n", err)
	} else {
		fileLogger.Info("This message will be written to logs/app.log")
		fileLogger.Close()
	}

	// Example 3: JSON Logging
	fmt.Println("\n3. JSON Logging")
	jsonLogger := spoor.QuickJSON()
	jsonLogger.Info("This message will be formatted as JSON")
	jsonLogger.WithField("user_id", 123).Info("User action")
	jsonLogger.Close()

	// Example 4: Async Logging
	fmt.Println("\n4. Async Logging")
	asyncLogger := spoor.QuickAsync()
	for i := 0; i < 10; i++ {
		asyncLogger.Infof("Async message %d", i)
	}
	asyncLogger.Sync() // Wait for all messages to be processed
	asyncLogger.Close()

	// Example 5: Structured Logging
	fmt.Println("\n5. Structured Logging")
	structuredLogger := spoor.Quick()
	structuredLogger.WithField("user_id", 123).
		WithField("action", "login").
		WithField("ip", "192.168.1.1").
		Info("User logged in")

	structuredLogger.WithFields(map[string]interface{}{
		"request_id": "req-123",
		"duration":   time.Millisecond * 150,
		"status":     200,
		"method":     "GET",
		"path":       "/api/users",
	}).Info("Request completed")
	structuredLogger.Close()

	// Example 6: Error Logging
	fmt.Println("\n6. Error Logging")
	errorLogger := spoor.Quick()
	
	// Simulate an error
	err = fmt.Errorf("database connection failed")
	errorLogger.WithError(err).Error("Failed to connect to database")
	
	errorLogger.WithField("error_code", "DB_CONN_001").
		WithField("retry_count", 3).
		Error("Database connection retry failed")
	errorLogger.Close()

	// Example 7: Global Logger
	fmt.Println("\n7. Global Logger")
	// Use global functions
	spoor.Info("This uses the global logger")
	spoor.WithField("component", "auth").Info("Authentication successful")
	
	// Set a custom global logger
	customLogger := spoor.QuickJSON()
	spoor.SetGlobalSimple(customLogger)
	spoor.Info("This now uses the custom global logger")
	customLogger.Close()

	// Example 8: Performance Logging
	fmt.Println("\n8. Performance Logging")
	perfLogger := spoor.Quick()
	
	start := time.Now()
	// Simulate some work
	time.Sleep(10 * time.Millisecond)
	duration := time.Since(start)
	
	perfLogger.WithField("operation", "data_processing").
		WithField("duration_ms", duration.Milliseconds()).
		WithField("records_processed", 1000).
		Info("Data processing completed")
	perfLogger.Close()

	// Example 9: Different Log Levels
	fmt.Println("\n9. Different Log Levels")
	levelLogger := spoor.Quick()
	levelLogger.SetLevel(spoor.LevelDebug) // Set to debug level
	
	levelLogger.Debug("This debug message will be shown")
	levelLogger.Info("This info message will be shown")
	levelLogger.Warn("This warning message will be shown")
	levelLogger.Error("This error message will be shown")
	
	levelLogger.SetLevel(spoor.LevelWarn) // Change to warn level
	levelLogger.Debug("This debug message will NOT be shown")
	levelLogger.Info("This info message will NOT be shown")
	levelLogger.Warn("This warning message will be shown")
	levelLogger.Error("This error message will be shown")
	levelLogger.Close()

	// Example 10: Formatted Messages
	fmt.Println("\n10. Formatted Messages")
	formatLogger := spoor.Quick()
	
	userID := 123
	action := "purchase"
	amount := 99.99
	
	formatLogger.Infof("User %d performed %s for $%.2f", userID, action, amount)
	formatLogger.Debugf("Processing request %d with %d parameters", 456, 3)
	formatLogger.Errorf("Failed to process request %d: %s", 789, "insufficient funds")
	formatLogger.Close()

	fmt.Println("\n✅ All examples completed!")
}
