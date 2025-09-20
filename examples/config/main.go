package main

import (
	"fmt"

	"github.com/phuhao00/spoor"
)

func main() {
	// 从配置文件创建日志记录器
	config, err := spoor.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		return
	}

	logger, err := spoor.CreateLoggerFromConfig(config)
	if err != nil {
		fmt.Printf("创建日志记录器失败: %v\n", err)
		return
	}

	// 使用日志记录器
	logger.Debug("这是一条调试消息")
	logger.Info("这是一条信息消息")
	logger.Warn("这是一条警告消息")
	logger.Error("这是一条错误消息")

	// 格式化消息
	logger.DebugF("用户 %s 登录成功", "张三")
	logger.InfoF("处理了 %d 个请求", 100)

	// 结构化日志
	structuredLogger := spoor.NewStructuredLogger(logger)
	structuredLogger.WithField("user_id", 123).WithField("action", "login").Info("用户登录")

	// 模拟一些业务逻辑
	simulateBusinessLogic(logger, structuredLogger)
}

func simulateBusinessLogic(logger *spoor.Spoor, structuredLogger *spoor.StructuredLogger) {
	// 模拟用户注册
	structuredLogger.WithFields(spoor.Fields{
		"event":     "user_registration",
		"user_id":   12345,
		"username":  "newuser",
		"email":     "newuser@example.com",
		"timestamp": "2024-01-01T10:00:00Z",
	}).Info("新用户注册")

	// 模拟API请求
	structuredLogger.WithFields(spoor.Fields{
		"event":       "api_request",
		"method":      "POST",
		"path":        "/api/v1/users",
		"status_code": 201,
		"duration_ms": 150,
		"client_ip":   "192.168.1.100",
		"user_agent":  "Mozilla/5.0...",
	}).Info("API请求处理完成")

	// 模拟错误
	err := fmt.Errorf("数据库连接超时")
	structuredLogger.WithFields(spoor.Fields{
		"event":       "database_error",
		"error":       err.Error(),
		"retry_count": 3,
		"duration_ms": 5000,
	}).Error("数据库操作失败")

	// 模拟性能监控
	structuredLogger.WithFields(spoor.Fields{
		"event":         "performance_metric",
		"metric_name":   "response_time",
		"avg_duration":  25.5,
		"max_duration":  100.0,
		"min_duration":  10.0,
		"request_count": 1000,
		"error_rate":    0.01,
	}).Info("性能指标统计")
}
