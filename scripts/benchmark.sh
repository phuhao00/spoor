#!/bin/bash

echo "🚀 Spoor Performance Benchmark"
echo "=============================="

# 设置环境变量
export GO111MODULE=on
export CGO_ENABLED=0

echo "📊 Running benchmark tests..."

# 运行基准测试
echo "1. Basic Logger Performance:"
go test -bench=BenchmarkCoreLogger -benchmem -count=3

echo "2. Async Logger Performance:"
go test -bench=BenchmarkAsyncLogger -benchmem -count=3

echo "3. Simple Logger Performance:"
go test -bench=BenchmarkSimpleLogger -benchmem -count=3

echo "4. Batch Writer Performance:"
go test -bench=BenchmarkBatchWriter -benchmem -count=3

echo "5. Memory Allocation Test:"
go test -bench=BenchmarkMemoryAllocation -benchmem -count=3

echo "6. Concurrent Logging Test:"
go test -bench=BenchmarkConcurrentLogging -benchmem -count=3

echo "7. Structured Logging Test:"
go test -bench=BenchmarkStructuredLogging -benchmem -count=3

echo "8. Formatted Logging Test:"
go test -bench=BenchmarkFormattedLogging -benchmem -count=3

echo "9. JSON Formatter Test:"
go test -bench=BenchmarkJSONFormatter -benchmem -count=3

echo "10. Text Formatter Test:"
go test -bench=BenchmarkTextFormatter -benchmem -count=3

echo "11. Sampling Test:"
go test -bench=BenchmarkSampling -benchmem -count=3

echo "12. Filtering Test:"
go test -bench=BenchmarkFiltering -benchmem -count=3

echo "13. Different Levels Test:"
go test -bench=BenchmarkDifferentLevels -benchmem -count=3

echo "14. With Fields Test:"
go test -bench=BenchmarkWithFields -benchmem -count=3

echo "✅ Benchmark completed!"
