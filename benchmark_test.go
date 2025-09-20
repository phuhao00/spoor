package spoor

import (
	"os"
	"testing"
	"time"
)

// BenchmarkSimpleLogger benchmarks the simple logger
func BenchmarkSimpleLogger(b *testing.B) {
	logger := Quick()
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message")
		}
	})
}

// BenchmarkAsyncLogger benchmarks the async logger
func BenchmarkAsyncLogger(b *testing.B) {
	config := DefaultSimpleConfig()
	config.Async = true
	logger, _ := NewSimple(config)
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message")
		}
	})
}

// BenchmarkCoreLogger benchmarks the core logger
func BenchmarkCoreLogger(b *testing.B) {
	writer := NewConsoleWriter(ConsoleWriterConfig{Output: os.Stdout})
	logger := NewCoreLogger(writer, LevelInfo)
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message")
		}
	})
}

// BenchmarkWithFields benchmarks logging with fields
func BenchmarkWithFields(b *testing.B) {
	logger := Quick()
	defer logger.Close()

	fields := map[string]interface{}{
		"user_id":    123,
		"request_id": "req-123",
		"duration":   time.Millisecond * 150,
		"status":     200,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.WithFields(fields).Info("benchmark message with fields")
		}
	})
}

// BenchmarkJSONFormatter benchmarks JSON formatting
func BenchmarkJSONFormatter(b *testing.B) {
	writer := NewConsoleWriter(ConsoleWriterConfig{Output: os.Stdout})
	formatter := NewJSONFormatter()
	logger := NewCoreLogger(writer, LevelInfo, WithFormatter(formatter))
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message")
		}
	})
}

// BenchmarkTextFormatter benchmarks text formatting
func BenchmarkTextFormatter(b *testing.B) {
	writer := NewConsoleWriter(ConsoleWriterConfig{Output: os.Stdout})
	formatter := NewTextFormatter()
	logger := NewCoreLogger(writer, LevelInfo, WithFormatter(formatter))
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message")
		}
	})
}

// BenchmarkBatchWriter benchmarks batch writing
func BenchmarkBatchWriter(b *testing.B) {
	writer := NewConsoleWriter(ConsoleWriterConfig{Output: os.Stdout})
	batchConfig := DefaultBatchConfig()
	batchConfig.BatchSize = 1000
	batchWriter := NewBatchWriter(writer, batchConfig)
	logger := NewCoreLogger(batchWriter, LevelInfo)
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message")
		}
	})
}

// BenchmarkSampling benchmarks log sampling
func BenchmarkSampling(b *testing.B) {
	writer := NewConsoleWriter(ConsoleWriterConfig{Output: os.Stdout})
	sampler := NewRateSampler(0.1) // 10% sampling
	config := AdvancedConfig{
		Sampler: sampler,
		Metrics: true,
	}
	logger := NewAdvancedLogger(writer, LevelInfo, config)
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message")
		}
	})
}

// BenchmarkFiltering benchmarks log filtering
func BenchmarkFiltering(b *testing.B) {
	writer := NewConsoleWriter(ConsoleWriterConfig{Output: os.Stdout})
	filter := NewLevelFilter(LevelWarn)
	config := AdvancedConfig{
		Filter:  filter,
		Metrics: true,
	}
	logger := NewAdvancedLogger(writer, LevelInfo, config)
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message")
		}
	})
}

// BenchmarkMemoryAllocation benchmarks memory allocation
func BenchmarkMemoryAllocation(b *testing.B) {
	logger := Quick()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message")
		}
	})
}

// BenchmarkConcurrentLogging benchmarks concurrent logging
func BenchmarkConcurrentLogging(b *testing.B) {
	logger := Quick()
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("concurrent benchmark message")
		}
	})
}

// BenchmarkDifferentLevels benchmarks different log levels
func BenchmarkDifferentLevels(b *testing.B) {
	logger := Quick()
	defer logger.Close()

	levels := []LogLevel{LevelDebug, LevelInfo, LevelWarn, LevelError}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			level := levels[i%len(levels)]
			switch level {
			case LevelDebug:
				logger.Debug("debug message")
			case LevelInfo:
				logger.Info("info message")
			case LevelWarn:
				logger.Warn("warn message")
			case LevelError:
				logger.Error("error message")
			}
			i++
		}
	})
}

// BenchmarkFormattedLogging benchmarks formatted logging
func BenchmarkFormattedLogging(b *testing.B) {
	logger := Quick()
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Infof("benchmark message with %d and %s", 123, "string")
		}
	})
}

// BenchmarkStructuredLogging benchmarks structured logging
func BenchmarkStructuredLogging(b *testing.B) {
	logger := Quick()
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.WithField("user_id", 123).
				WithField("action", "login").
				WithField("duration", time.Millisecond*150).
				Info("user login")
		}
	})
}
