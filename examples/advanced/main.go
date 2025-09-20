package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/phuhao00/spoor"
)

func main() {
	// 示例1: Elasticsearch输出
	elasticConfig := spoor.ElasticConfig{
		URL:       "http://localhost:9200",
		Index:     "app-logs",
		BatchSize: 100,
		FlushTime: 5 * time.Second,
	}

	elasticLogger := spoor.NewSpoor(
		spoor.DEBUG,
		"",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile,
		spoor.WithElasticWriter(elasticConfig),
	)

	elasticLogger.Info("这条消息将发送到Elasticsearch")
	elasticLogger.InfoF("用户操作: %s", "登录")

	// 示例2: ClickHouse输出
	clickhouseConfig := spoor.ClickHouseConfig{
		DSN:       "tcp://localhost:9000?database=logs",
		TableName: "app_logs",
		BatchSize: 100,
		FlushTime: 5 * time.Second,
	}

	clickhouseLogger := spoor.NewSpoor(
		spoor.DEBUG,
		"",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile,
		spoor.WithClickHouseWriter(clickhouseConfig),
	)

	clickhouseLogger.Info("这条消息将发送到ClickHouse")
	clickhouseLogger.InfoF("数据库查询耗时: %v", time.Millisecond*50)

	// 示例3: Logbus输出
	logbusConfig := spoor.LogbusConfig{
		URL:       "https://api.logbus.com/logs",
		Token:     "your-token-here",
		BatchSize: 100,
		FlushTime: 5 * time.Second,
	}

	logbusLogger := spoor.NewSpoor(
		spoor.DEBUG,
		"",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile,
		spoor.WithLogbusWriter(logbusConfig),
	)

	logbusLogger.Info("这条消息将发送到Logbus")
	logbusLogger.InfoF("API调用: %s", "GET /api/users")

	// 示例4: 复杂的结构化日志
	consoleLogger := spoor.NewSpoor(
		spoor.DEBUG,
		"",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile,
		spoor.WithConsoleWriter(os.Stdout),
	)

	structuredLogger := spoor.NewStructuredLogger(consoleLogger)

	// 模拟用户操作日志
	userLogger := structuredLogger.WithFields(spoor.Fields{
		"service": "user-service",
		"version": "1.0.0",
		"env":     "production",
	})

	userLogger.WithField("user_id", 12345).WithField("action", "profile_update").Info("用户更新了个人资料")
	userLogger.WithField("user_id", 12345).WithField("action", "password_change").Warn("用户修改了密码")

	// 模拟API请求日志
	apiLogger := structuredLogger.WithFields(spoor.Fields{
		"service": "api-gateway",
		"version": "2.1.0",
		"env":     "production",
	})

	apiLogger.WithFields(spoor.Fields{
		"method":      "POST",
		"path":        "/api/v1/users",
		"status_code": 201,
		"duration_ms": 150,
		"client_ip":   "192.168.1.100",
	}).Info("API请求处理完成")

	// 模拟错误日志
	err := fmt.Errorf("数据库连接超时")
	apiLogger.WithFields(spoor.Fields{
		"method":      "GET",
		"path":        "/api/v1/users/123",
		"error":       err.Error(),
		"duration_ms": 5000,
	}).Error("API请求失败")

	// 示例5: 性能监控日志
	perfLogger := structuredLogger.WithFields(spoor.Fields{
		"component": "performance-monitor",
		"metric":    "response_time",
	})

	perfLogger.WithFields(spoor.Fields{
		"endpoint":      "/api/v1/health",
		"avg_duration":  25.5,
		"max_duration":  100.0,
		"min_duration":  10.0,
		"request_count": 1000,
	}).Info("性能指标统计")

	// 等待异步写入完成
	time.Sleep(6 * time.Second)
}
