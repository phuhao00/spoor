package main

import (
	"fmt"
	"os"
	"time"

	"github.com/phuhao00/spoor"
)

func main() {
	fmt.Println("🚀 Spoor Performance Examples")
	fmt.Println("==============================")

	// Example 1: Simple Logger
	fmt.Println("\n1. Simple Logger (Default)")
	simpleLogger := spoor.Quick()
	simpleLogger.Info("This is a simple info message")
	simpleLogger.WithField("user_id", 123).Info("User logged in")
	simpleLogger.Close()

	// Example 2: Async Logger
	fmt.Println("\n2. Async Logger (High Performance)")
	asyncConfig := spoor.DefaultAsyncConfig()
	asyncConfig.BufferSize = 10000
	asyncConfig.WorkerCount = 4
	
	writer := spoor.NewConsoleWriter(spoor.ConsoleWriterConfig{Output: os.Stdout})
	asyncLogger := spoor.NewAsyncLogger(writer, spoor.LevelInfo, asyncConfig)
	
	// Log many messages quickly
	start := time.Now()
	for i := 0; i < 1000; i++ {
		asyncLogger.Infof("Async message %d", i)
	}
	asyncLogger.Sync() // Wait for all messages to be processed
	duration := time.Since(start)
	fmt.Printf("Logged 1000 messages in %v\n", duration)
	asyncLogger.Close()

	// Example 3: Batch Writer
	fmt.Println("\n3. Batch Writer")
	batchConfig := spoor.DefaultBatchConfig()
	batchConfig.BatchSize = 100
	batchConfig.FlushInterval = 50 * time.Millisecond
	
	batchWriter := spoor.NewBatchWriter(writer, batchConfig)
	batchLogger := spoor.NewCoreLogger(batchWriter, spoor.LevelInfo)
	
	start = time.Now()
	for i := 0; i < 1000; i++ {
		batchLogger.Infof("Batch message %d", i)
	}
	batchLogger.Sync()
	duration = time.Since(start)
	fmt.Printf("Logged 1000 messages with batching in %v\n", duration)
	batchLogger.Close()

	// Example 4: Advanced Logger with Sampling
	fmt.Println("\n4. Advanced Logger with Sampling")
	sampler := spoor.NewRateSampler(0.1) // 10% sampling
	advancedConfig := spoor.AdvancedConfig{
		Sampler: sampler,
		Metrics: true,
	}
	advancedLogger := spoor.NewAdvancedLogger(writer, spoor.LevelInfo, advancedConfig)
	
	start = time.Now()
	for i := 0; i < 1000; i++ {
		advancedLogger.Infof("Sampled message %d", i)
	}
	advancedLogger.Sync()
	duration = time.Since(start)
	
	metrics := advancedLogger.GetMetrics()
	fmt.Printf("Logged 1000 messages with sampling in %v\n", duration)
	fmt.Printf("Metrics: %+v\n", metrics)
	advancedLogger.Close()

	// Example 5: Memory Pool Usage
	fmt.Println("\n5. Memory Pool Usage")
	poolLogger := spoor.Quick()
	
	start = time.Now()
	for i := 0; i < 1000; i++ {
		// This will use the memory pool internally
		poolLogger.WithFields(map[string]interface{}{
			"request_id": fmt.Sprintf("req-%d", i),
			"user_id":    i,
			"duration":   time.Millisecond * time.Duration(i%100),
		}).Info("Pooled message")
	}
	poolLogger.Sync()
	duration = time.Since(start)
	fmt.Printf("Logged 1000 messages with memory pooling in %v\n", duration)
	poolLogger.Close()

	// Example 6: Performance Comparison
	fmt.Println("\n6. Performance Comparison")
	comparePerformance()
}

func comparePerformance() {
	iterations := 10000
	
	// Test 1: Simple Logger
	start := time.Now()
	simpleLogger := spoor.Quick()
	for i := 0; i < iterations; i++ {
		simpleLogger.Info("Simple message")
	}
	simpleLogger.Sync()
	simpleDuration := time.Since(start)
	simpleLogger.Close()

	// Test 2: Async Logger
	start = time.Now()
	asyncConfig := spoor.DefaultAsyncConfig()
	asyncConfig.BufferSize = 10000
	writer := spoor.NewConsoleWriter(spoor.ConsoleWriterConfig{Output: os.Stdout})
	asyncLogger := spoor.NewAsyncLogger(writer, spoor.LevelInfo, asyncConfig)
	for i := 0; i < iterations; i++ {
		asyncLogger.Info("Async message")
	}
	asyncLogger.Sync()
	asyncDuration := time.Since(start)
	asyncLogger.Close()

	// Test 3: Batch Writer
	start = time.Now()
	batchConfig := spoor.DefaultBatchConfig()
	batchConfig.BatchSize = 1000
	batchWriter := spoor.NewBatchWriter(writer, batchConfig)
	batchLogger := spoor.NewCoreLogger(batchWriter, spoor.LevelInfo)
	for i := 0; i < iterations; i++ {
		batchLogger.Info("Batch message")
	}
	batchLogger.Sync()
	batchDuration := time.Since(start)
	batchLogger.Close()

	fmt.Printf("Simple Logger:  %v (%d messages/sec)\n", 
		simpleDuration, int64(time.Duration(iterations)*time.Second)/int64(simpleDuration))
	fmt.Printf("Async Logger:   %v (%d messages/sec)\n", 
		asyncDuration, int64(time.Duration(iterations)*time.Second)/int64(asyncDuration))
	fmt.Printf("Batch Logger:   %v (%d messages/sec)\n", 
		batchDuration, int64(time.Duration(iterations)*time.Second)/int64(batchDuration))
}
