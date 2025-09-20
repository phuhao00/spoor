package main

import (
	"fmt"
	"time"

	"github.com/phuhao00/spoor"
)

func main() {
	fmt.Println("=== ClickHouse 日志示例 ===\n")

	// 示例1: 基本 ClickHouse 日志
	fmt.Println("1. 基本 ClickHouse 日志")
	fmt.Println("=======================")
	clickhouseLogger, err := spoor.NewClickHouse("tcp://localhost:9000?database=logs", "app_logs", spoor.LevelInfo)
	if err != nil {
		fmt.Printf("创建 ClickHouse 日志记录器失败: %v\n", err)
		fmt.Println("注意: 需要运行 ClickHouse 实例才能使用此功能")
		return
	}
	defer clickhouseLogger.Close()

	clickhouseLogger.Info("这是一条发送到 ClickHouse 的日志")
	clickhouseLogger.Info("ClickHouse 支持高性能的日志存储和查询")
	clickhouseLogger.Info("支持列式存储和压缩")
	fmt.Println("基本日志已发送到 ClickHouse")
	fmt.Println()

	// 示例2: 结构化日志
	fmt.Println("2. 结构化日志")
	fmt.Println("=============")
	structuredLogger := clickhouseLogger.WithFields(map[string]interface{}{
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
	fmt.Println("结构化日志已发送到 ClickHouse")
	fmt.Println()

	// 示例3: 错误日志
	fmt.Println("3. 错误日志")
	fmt.Println("===========")
	errorLogger := clickhouseLogger.WithFields(map[string]interface{}{
		"component":  "database",
		"operation":  "query",
		"error_type": "connection_timeout",
	})
	errorLogger.WithError(fmt.Errorf("连接超时")).Error("数据库连接失败")
	errorLogger.WithError(fmt.Errorf("语法错误")).Error("SQL 查询失败")
	fmt.Println("错误日志已发送到 ClickHouse")
	fmt.Println()

	// 示例4: 性能监控日志
	fmt.Println("4. 性能监控日志")
	fmt.Println("===============")
	perfLogger := clickhouseLogger.WithFields(map[string]interface{}{
		"type":        "performance",
		"service":     "api-gateway",
		"metric_type": "latency",
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

	fmt.Println("性能监控日志已发送到 ClickHouse")
	fmt.Println()

	// 示例5: 批量日志写入
	fmt.Println("5. 批量日志写入")
	fmt.Println("===============")
	batchLogger := clickhouseLogger.WithFields(map[string]interface{}{
		"batch_id":   "batch-001",
		"type":       "bulk_operation",
		"batch_size": 100,
	})

	// 模拟批量操作
	for i := 0; i < 50; i++ {
		batchLogger.WithFields(map[string]interface{}{
			"item_id":            i,
			"operation":          "process",
			"status":             "success",
			"processing_time_ms": 10 + i%50,
		}).Info("批量处理项目")
	}
	fmt.Println("批量日志已发送到 ClickHouse")
	fmt.Println()

	// 示例6: 不同日志级别
	fmt.Println("6. 不同日志级别")
	fmt.Println("===============")
	levelLogger, err := spoor.NewClickHouse("tcp://localhost:9000?database=logs", "level_logs", spoor.LevelDebug)
	if err != nil {
		fmt.Printf("创建级别日志失败: %v\n", err)
		return
	}
	defer levelLogger.Close()

	levelLogger.Debug("调试信息：处理用户请求")
	levelLogger.Info("信息：用户请求处理完成")
	levelLogger.Warn("警告：内存使用率较高")
	levelLogger.Error("错误：数据库连接失败")
	fmt.Println("不同级别的日志已发送到 ClickHouse")
	fmt.Println()

	// 示例7: 时间序列日志
	fmt.Println("7. 时间序列日志")
	fmt.Println("===============")
	timeSeriesLogger := clickhouseLogger.WithFields(map[string]interface{}{
		"type":   "time_series",
		"metric": "system_metrics",
	})

	// 模拟时间序列数据
	now := time.Now()
	for i := 0; i < 10; i++ {
		timestamp := now.Add(time.Duration(i) * time.Minute)
		timeSeriesLogger.WithFields(map[string]interface{}{
			"timestamp":    timestamp,
			"cpu_usage":    20.0 + float64(i),
			"memory_usage": 1024 + i*100,
			"disk_usage":   50.0 + float64(i)*0.5,
		}).Info("系统指标")
	}
	fmt.Println("时间序列日志已发送到 ClickHouse")
	fmt.Println()

	// 示例8: 业务事件日志
	fmt.Println("8. 业务事件日志")
	fmt.Println("===============")
	businessLogger := clickhouseLogger.WithFields(map[string]interface{}{
		"type":   "business_event",
		"domain": "ecommerce",
	})

	// 模拟电商业务事件
	businessLogger.WithFields(map[string]interface{}{
		"event":    "order_created",
		"order_id": "ORD-001",
		"user_id":  12345,
		"amount":   99.99,
		"currency": "USD",
	}).Info("订单创建")

	businessLogger.WithFields(map[string]interface{}{
		"event":          "payment_processed",
		"order_id":       "ORD-001",
		"payment_method": "credit_card",
		"transaction_id": "TXN-123",
	}).Info("支付处理")

	businessLogger.WithFields(map[string]interface{}{
		"event":           "order_shipped",
		"order_id":        "ORD-001",
		"tracking_number": "TRK-456",
		"shipping_method": "standard",
	}).Info("订单发货")

	fmt.Println("业务事件日志已发送到 ClickHouse")
	fmt.Println()

	// 等待日志发送完成
	fmt.Println("等待日志发送完成...")
	time.Sleep(3 * time.Second)
	clickhouseLogger.Sync()

	fmt.Println("=== ClickHouse 日志示例完成 ===")
	fmt.Println("注意：需要运行 ClickHouse 实例才能看到实际效果")
}
