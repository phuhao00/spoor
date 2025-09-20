package main

import (
	"fmt"
	"time"

	"github.com/phuhao00/spoor"
)

func main() {
	fmt.Println("=== Elasticsearch 日志示例 ===\n")

	// 示例1: 基本 Elasticsearch 日志
	fmt.Println("1. 基本 Elasticsearch 日志")
	fmt.Println("===========================")
	elasticLogger := spoor.NewElastic("http://localhost:9200", "app-logs", spoor.LevelInfo)
	defer elasticLogger.Close()

	elasticLogger.Info("这是一条发送到 Elasticsearch 的日志")
	elasticLogger.Info("Elasticsearch 支持高性能的日志存储和查询")
	elasticLogger.Info("支持实时搜索和分析")
	fmt.Println("基本日志已发送到 Elasticsearch")
	fmt.Println()

	// 示例2: 结构化日志
	fmt.Println("2. 结构化日志")
	fmt.Println("=============")
	structuredLogger := elasticLogger.WithFields(map[string]interface{}{
		"service":    "user-service",
		"version":    "1.0.0",
		"request_id": "req-123",
		"user_id":    12345,
		"ip":         "192.168.1.100",
	})
	structuredLogger.Info("用户操作日志")

	// 模拟不同类型的日志
	structuredLogger.WithField("action", "login").Info("用户登录")
	structuredLogger.WithField("action", "logout").Info("用户登出")
	structuredLogger.WithField("action", "profile_update").Info("用户更新资料")
	fmt.Println("结构化日志已发送到 Elasticsearch")
	fmt.Println()

	// 示例3: 错误日志
	fmt.Println("3. 错误日志")
	fmt.Println("===========")
	errorLogger := elasticLogger.WithFields(map[string]interface{}{
		"component": "database",
		"operation": "query",
	})
	errorLogger.WithError(fmt.Errorf("连接超时")).Error("数据库连接失败")
	errorLogger.WithError(fmt.Errorf("语法错误")).Error("SQL 查询失败")
	fmt.Println("错误日志已发送到 Elasticsearch")
	fmt.Println()

	// 示例4: 性能监控日志
	fmt.Println("4. 性能监控日志")
	fmt.Println("===============")
	perfLogger := elasticLogger.WithFields(map[string]interface{}{
		"type":    "performance",
		"service": "api-gateway",
	})

	// 模拟性能指标
	perfLogger.WithFields(map[string]interface{}{
		"endpoint":      "/api/users",
		"method":        "GET",
		"duration_ms":   150,
		"status_code":   200,
		"response_size": 1024,
	}).Info("API 请求完成")

	perfLogger.WithFields(map[string]interface{}{
		"endpoint":      "/api/orders",
		"method":        "POST",
		"duration_ms":   2500,
		"status_code":   201,
		"response_size": 512,
	}).Warn("API 请求较慢")

	fmt.Println("性能监控日志已发送到 Elasticsearch")
	fmt.Println()

	// 示例5: 批量日志写入
	fmt.Println("5. 批量日志写入")
	fmt.Println("===============")
	batchLogger := elasticLogger.WithFields(map[string]interface{}{
		"batch_id": "batch-001",
		"type":     "bulk_operation",
	})

	// 模拟批量操作
	for i := 0; i < 50; i++ {
		batchLogger.WithFields(map[string]interface{}{
			"item_id":   i,
			"operation": "process",
			"status":    "success",
		}).Info("批量处理项目")
	}
	fmt.Println("批量日志已发送到 Elasticsearch")
	fmt.Println()

	// 示例6: 不同日志级别
	fmt.Println("6. 不同日志级别")
	fmt.Println("===============")
	levelLogger := spoor.NewElastic("http://localhost:9200", "level-logs", spoor.LevelDebug)
	defer levelLogger.Close()

	levelLogger.Debug("调试信息：处理用户请求")
	levelLogger.Info("信息：用户请求处理完成")
	levelLogger.Warn("警告：内存使用率较高")
	levelLogger.Error("错误：数据库连接失败")
	fmt.Println("不同级别的日志已发送到 Elasticsearch")
	fmt.Println()

	// 示例7: 应用生命周期日志
	fmt.Println("7. 应用生命周期日志")
	fmt.Println("===================")
	lifecycleLogger := elasticLogger.WithFields(map[string]interface{}{
		"component": "application",
		"phase":     "startup",
	})

	lifecycleLogger.Info("应用启动")
	lifecycleLogger.WithField("phase", "initialization").Info("初始化配置")
	lifecycleLogger.WithField("phase", "database").Info("连接数据库")
	lifecycleLogger.WithField("phase", "cache").Info("初始化缓存")
	lifecycleLogger.WithField("phase", "server").Info("启动 HTTP 服务器")
	lifecycleLogger.WithField("phase", "ready").Info("应用就绪")

	// 模拟运行时的日志
	runtimeLogger := elasticLogger.WithFields(map[string]interface{}{
		"component": "application",
		"phase":     "runtime",
	})

	runtimeLogger.WithFields(map[string]interface{}{
		"metric": "cpu_usage",
		"value":  45.2,
		"unit":   "percent",
	}).Info("系统指标")

	runtimeLogger.WithFields(map[string]interface{}{
		"metric": "memory_usage",
		"value":  1024,
		"unit":   "MB",
	}).Info("系统指标")

	fmt.Println("应用生命周期日志已发送到 Elasticsearch")
	fmt.Println()

	// 等待日志发送完成
	fmt.Println("等待日志发送完成...")
	time.Sleep(3 * time.Second)
	elasticLogger.Sync()

	fmt.Println("=== Elasticsearch 日志示例完成 ===")
	fmt.Println("注意：需要运行 Elasticsearch 实例才能看到实际效果")
}
