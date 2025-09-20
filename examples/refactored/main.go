package main

import (
	"fmt"
	"time"

	"github.com/phuhao00/spoor"
)

func main() {
	fmt.Println("=== Spoor 重构后的使用示例 ===\n")

	// 示例1: 基本使用
	basicUsage()

	// 示例2: 结构化日志
	structuredLogging()

	// 示例3: 不同格式化器
	differentFormatters()

	// 示例4: 文件日志
	fileLogging()

	// 示例5: ClickHouse日志
	clickhouseLogging()

	// 示例6: 配置化日志
	configBasedLogging()

	// 示例7: 并发日志
	concurrentLogging()
}

func basicUsage() {
	fmt.Println("1. 基本使用示例")
	fmt.Println("================")

	// 创建控制台日志记录器
	logger := spoor.NewConsole(spoor.LevelDebug)

	logger.Debug("这是一条调试消息")
	logger.Info("这是一条信息消息")
	logger.Warn("这是一条警告消息")
	logger.Error("这是一条错误消息")

	// 格式化消息
	logger.Debugf("用户 %s 登录成功", "张三")
	logger.Infof("处理了 %d 个请求", 100)

	fmt.Println()
}

func structuredLogging() {
	fmt.Println("2. 结构化日志示例")
	fmt.Println("==================")

	writer := spoor.NewWriterFactory().CreateConsoleWriterToStdout()
	logger := spoor.NewJSON(writer, spoor.LevelInfo)

	// 添加字段
	structuredLogger := logger.WithField("user_id", 123).WithField("action", "login")
	structuredLogger.Info("用户登录")

	// 添加多个字段
	structuredLogger = logger.WithFields(map[string]interface{}{
		"request_id": "req-123",
		"duration":   time.Millisecond * 150,
		"status":     200,
		"client_ip":  "192.168.1.100",
	})
	structuredLogger.Info("API请求处理完成")

	// 添加错误
	err := fmt.Errorf("数据库连接超时")
	structuredLogger = logger.WithError(err)
	structuredLogger.Error("数据库操作失败")

	fmt.Println()
}

func differentFormatters() {
	fmt.Println("3. 不同格式化器示例")
	fmt.Println("====================")

	writer := spoor.NewWriterFactory().CreateConsoleWriterToStdout()

	// 文本格式化器
	fmt.Println("文本格式化器:")
	textLogger := spoor.NewText(writer, spoor.LevelInfo)
	textLogger.Info("这是文本格式的日志")

	// JSON格式化器
	fmt.Println("\nJSON格式化器:")
	jsonLogger := spoor.NewJSON(writer, spoor.LevelInfo)
	jsonLogger.WithField("key", "value").Info("这是JSON格式的日志")

	// Logrus格式化器
	fmt.Println("\nLogrus格式化器:")
	logrusLogger := spoor.NewLogrus(writer, spoor.LevelInfo)
	logrusLogger.WithField("key", "value").Info("这是Logrus格式的日志")

	fmt.Println()
}

func fileLogging() {
	fmt.Println("4. 文件日志示例")
	fmt.Println("===============")

	// 创建文件日志记录器
	logger, err := spoor.NewFile("logs", spoor.LevelInfo)
	if err != nil {
		fmt.Printf("创建文件日志记录器失败: %v\n", err)
		return
	}

	logger.Info("这条消息将写入到文件")
	logger.Info("文件日志支持自动轮转")
	logger.Info("支持批量写入以提高性能")

	// 等待写入完成
	time.Sleep(100 * time.Millisecond)

	fmt.Println("文件日志已写入到 logs/ 目录")
	fmt.Println()
}

func clickhouseLogging() {
	fmt.Println("5. ClickHouse日志示例")
	fmt.Println("=====================")

	// 创建ClickHouse日志记录器
	logger, err := spoor.NewClickHouse("tcp://localhost:9000?database=logs", "app_logs", spoor.LevelInfo)
	if err != nil {
		fmt.Printf("创建ClickHouse日志记录器失败: %v\n", err)
		fmt.Println("注意: 需要运行ClickHouse实例才能使用此功能")
		return
	}

	logger.Info("这条消息将发送到ClickHouse")
	logger.Info("ClickHouse支持高性能的日志存储和查询")
	logger.Info("支持结构化日志和批量写入")

	// 结构化日志
	structuredLogger := logger.WithFields(map[string]interface{}{
		"service":    "user-service",
		"version":    "1.0.0",
		"request_id": "req-123",
		"user_id":    12345,
	})
	structuredLogger.Info("用户操作日志")

	// 等待写入完成
	time.Sleep(2 * time.Second)

	fmt.Println("ClickHouse日志已发送")
	fmt.Println()
}

func configBasedLogging() {
	fmt.Println("6. 配置化日志示例")
	fmt.Println("==================")

	// 使用默认配置
	config := spoor.DefaultConfig()
	config.Level = spoor.LevelDebug

	logger, err := spoor.NewFromConfig(config)
	if err != nil {
		fmt.Printf("从配置创建日志记录器失败: %v\n", err)
		return
	}

	logger.Info("这是基于配置的日志记录器")
	logger.Debug("调试消息也会显示")

	// 修改配置
	config.Level = spoor.LevelInfo
	logger.SetLevel(spoor.LevelInfo)
	logger.Debug("这条调试消息不会显示")
	logger.Info("这条信息消息会显示")

	fmt.Println()
}

func concurrentLogging() {
	fmt.Println("7. 并发日志示例")
	fmt.Println("===============")

	writer := spoor.NewWriterFactory().CreateConsoleWriterToStdout()
	logger := spoor.New(writer, spoor.LevelInfo)

	// 启动多个goroutine并发写入日志
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 3; j++ {
				logger.Infof("Goroutine %d, 消息 %d", id, j)
			}
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 5; i++ {
		<-done
	}

	fmt.Println("并发日志写入完成")
	fmt.Println()
}

func demonstrateBusinessLogic() {
	fmt.Println("7. 业务逻辑日志示例")
	fmt.Println("====================")

	writer := spoor.NewWriterFactory().CreateConsoleWriterToStdout()
	logger := spoor.NewJSON(writer, spoor.LevelInfo)

	// 模拟用户注册
	logger.WithFields(map[string]interface{}{
		"event":     "user_registration",
		"user_id":   12345,
		"username":  "newuser",
		"email":     "newuser@example.com",
		"timestamp": time.Now().Format(time.RFC3339),
	}).Info("新用户注册")

	// 模拟API请求
	logger.WithFields(map[string]interface{}{
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
	logger.WithFields(map[string]interface{}{
		"event":       "database_error",
		"error":       err.Error(),
		"retry_count": 3,
		"duration_ms": 5000,
	}).Error("数据库操作失败")

	// 模拟性能监控
	logger.WithFields(map[string]interface{}{
		"event":         "performance_metric",
		"metric_name":   "response_time",
		"avg_duration":  25.5,
		"max_duration":  100.0,
		"min_duration":  10.0,
		"request_count": 1000,
		"error_rate":    0.01,
	}).Info("性能指标统计")

	fmt.Println()
}
