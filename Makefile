# Spoor Makefile

.PHONY: help build test benchmark clean install examples

# Default target
help:
	@echo "Spoor - 高性能Go日志库"
	@echo "======================"
	@echo ""
	@echo "可用命令:"
	@echo "  build      - 编译项目"
	@echo "  test       - 运行测试"
	@echo "  benchmark  - 运行性能基准测试"
	@echo "  examples   - 运行示例程序"
	@echo "  clean      - 清理构建文件"
	@echo "  install    - 安装到系统"
	@echo "  lint       - 运行代码检查"
	@echo "  fmt        - 格式化代码"
	@echo "  docs       - 生成文档"

# Build the project
build:
	@echo "🔨 编译项目..."
	go build -v ./...

# Run tests
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# Run benchmarks
benchmark:
	@echo "⚡ 运行性能基准测试..."
	go test -bench=. -benchmem -count=3

# Run specific benchmark
benchmark-async:
	@echo "⚡ 运行异步日志器基准测试..."
	go test -bench=BenchmarkAsyncLogger -benchmem -count=3

benchmark-simple:
	@echo "⚡ 运行简单日志器基准测试..."
	go test -bench=BenchmarkSimpleLogger -benchmem -count=3

benchmark-batch:
	@echo "⚡ 运行批量写入器基准测试..."
	go test -bench=BenchmarkBatchWriter -benchmem -count=3

# Run examples
examples:
	@echo "📝 运行示例程序..."
	@echo "1. 简单使用示例:"
	go run examples/simple_usage/main.go
	@echo ""
	@echo "2. 性能示例:"
	go run examples/performance/main.go
	@echo ""
	@echo "3. 完整示例:"
	go run examples/complete/main.go

# Run specific example
example-simple:
	@echo "📝 运行简单使用示例..."
	go run examples/simple_usage/main.go

example-performance:
	@echo "📝 运行性能示例..."
	go run examples/performance/main.go

example-complete:
	@echo "📝 运行完整示例..."
	go run examples/complete/main.go

# Clean build files
clean:
	@echo "🧹 清理构建文件..."
	go clean
	rm -f spoor-config.json
	rm -rf logs/

# Install to system
install:
	@echo "📦 安装到系统..."
	go install ./...

# Run linter
lint:
	@echo "🔍 运行代码检查..."
	golangci-lint run

# Format code
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...

# Generate documentation
docs:
	@echo "📚 生成文档..."
	godoc -http=:6060

# Run all checks
check: fmt lint test benchmark

# Performance test
perf-test:
	@echo "🚀 性能测试..."
	@echo "测试环境: $(shell go version)"
	@echo "CPU核心数: $(shell nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo "unknown")"
	@echo "内存: $(shell free -h 2>/dev/null | head -2 | tail -1 || echo "unknown")"
	@echo ""
	@make benchmark

# Memory test
mem-test:
	@echo "💾 内存测试..."
	go test -bench=BenchmarkMemoryAllocation -benchmem -count=5

# Concurrent test
concurrent-test:
	@echo "🔄 并发测试..."
	go test -bench=BenchmarkConcurrentLogging -benchmem -count=3

# All tests
test-all: test benchmark examples

# Quick start
quick-start:
	@echo "🚀 快速开始..."
	@make example-simple

# Development setup
dev-setup:
	@echo "🛠️  开发环境设置..."
	go mod tidy
	go mod download
	@echo "✅ 开发环境设置完成"

# Release build
release:
	@echo "📦 发布构建..."
	@echo "构建版本: $(shell git describe --tags --always --dirty)"
	go build -ldflags "-X main.version=$(shell git describe --tags --always --dirty)" -o spoor ./cmd/spoor

# Docker build
docker-build:
	@echo "🐳 构建Docker镜像..."
	docker build -t spoor:latest .

# Docker run
docker-run:
	@echo "🐳 运行Docker容器..."
	docker run --rm -it spoor:latest

# Show project info
info:
	@echo "📊 项目信息:"
	@echo "  项目名称: Spoor"
	@echo "  版本: $(shell git describe --tags --always --dirty)"
	@echo "  Go版本: $(shell go version)"
	@echo "  模块路径: $(shell go list -m)"
	@echo "  文件数量: $(shell find . -name "*.go" | wc -l)"
	@echo "  代码行数: $(shell find . -name "*.go" -exec wc -l {} + | tail -1 | awk '{print $$1}')"
