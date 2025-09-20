package main

import (
	"fmt"
	"os"

	"github.com/phuhao00/spoor"
)

func main() {
	fmt.Println("=== 文件日志示例 ===\n")

	// 示例1: 基本文件日志
	fmt.Println("1. 基本文件日志")
	fmt.Println("===============")
	fileLogger, err := spoor.NewFile("logs", spoor.LevelDebug)
	if err != nil {
		fmt.Printf("创建文件日志失败: %v\n", err)
		return
	}
	defer func() {
		fileLogger.Close()
		os.RemoveAll("logs") // 清理测试文件
	}()

	fileLogger.Debug("这是一条调试消息")
	fileLogger.Info("这是一条信息消息")
	fileLogger.Warn("这是一条警告消息")
	fileLogger.Error("这是一条错误消息")
	fmt.Println("文件日志已写入到 logs/ 目录")
	fmt.Println()

	// 示例2: 自定义文件配置
	fmt.Println("2. 自定义文件配置")
	fmt.Println("=================")
	customFileLogger, err := spoor.NewFile("custom_logs", spoor.LevelInfo)
	if err != nil {
		fmt.Printf("创建自定义文件日志失败: %v\n", err)
		return
	}
	defer func() {
		customFileLogger.Close()
		os.RemoveAll("custom_logs")
	}()

	// 设置自定义格式化器
	customFileLogger.SetFormatter(spoor.NewJSONFormatter())
	customFileLogger.Info("JSON 格式的文件日志")
	customFileLogger.WithField("service", "file-logger").Info("结构化文件日志")
	fmt.Println("自定义文件日志已写入到 custom_logs/ 目录")
	fmt.Println()

	// 示例3: 日志轮转
	fmt.Println("3. 日志轮转测试")
	fmt.Println("===============")
	rotationLogger, err := spoor.NewFile("rotation_logs", spoor.LevelDebug)
	if err != nil {
		fmt.Printf("创建轮转日志失败: %v\n", err)
		return
	}
	defer func() {
		rotationLogger.Close()
		os.RemoveAll("rotation_logs")
	}()

	// 写入大量日志以触发轮转
	for i := 0; i < 1000; i++ {
		rotationLogger.Infof("轮转测试日志 %d", i)
	}
	fmt.Println("轮转测试完成，检查 rotation_logs/ 目录中的文件")
	fmt.Println()

	// 示例4: 并发文件写入
	fmt.Println("4. 并发文件写入")
	fmt.Println("===============")
	concurrentLogger, err := spoor.NewFile("concurrent_logs", spoor.LevelInfo)
	if err != nil {
		fmt.Printf("创建并发日志失败: %v\n", err)
		return
	}
	defer func() {
		concurrentLogger.Close()
		os.RemoveAll("concurrent_logs")
	}()

	// 启动多个 goroutine 并发写入
	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				concurrentLogger.Infof("Goroutine %d, 消息 %d", id, j)
			}
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 5; i++ {
		<-done
	}

	// 确保所有日志都已刷新
	concurrentLogger.Sync()
	fmt.Println("并发文件写入完成")
	fmt.Println()

	// 示例5: 不同日志级别的文件
	fmt.Println("5. 不同日志级别的文件")
	fmt.Println("=====================")
	levelFileLogger, err := spoor.NewFile("level_logs", spoor.LevelWarn)
	if err != nil {
		fmt.Printf("创建级别日志失败: %v\n", err)
		return
	}
	defer func() {
		levelFileLogger.Close()
		os.RemoveAll("level_logs")
	}()

	levelFileLogger.Debug("这条调试消息不会写入文件")
	levelFileLogger.Info("这条信息消息不会写入文件")
	levelFileLogger.Warn("这条警告消息会写入文件")
	levelFileLogger.Error("这条错误消息会写入文件")
	fmt.Println("只有 WARN 和 ERROR 级别的日志被写入文件")
	fmt.Println()

	// 示例6: 结构化文件日志
	fmt.Println("6. 结构化文件日志")
	fmt.Println("=================")
	structuredFileLogger, err := spoor.NewFile("structured_logs", spoor.LevelInfo)
	if err != nil {
		fmt.Printf("创建结构化文件日志失败: %v\n", err)
		return
	}
	defer func() {
		structuredFileLogger.Close()
		os.RemoveAll("structured_logs")
	}()

	// 设置 JSON 格式化器用于结构化日志
	structuredFileLogger.SetFormatter(spoor.NewJSONFormatter())

	// 模拟应用日志
	appLogger := structuredFileLogger.WithFields(map[string]interface{}{
		"service": "user-service",
		"version": "1.0.0",
		"env":     "production",
	})

	appLogger.Info("服务启动")
	appLogger.WithField("port", 8080).Info("HTTP 服务器启动")
	appLogger.WithField("user_id", 12345).Info("用户登录")
	appLogger.WithField("request_id", "req-abc123").Warn("请求处理缓慢")
	appLogger.WithError(fmt.Errorf("数据库连接超时")).Error("数据库操作失败")

	fmt.Println("结构化文件日志已写入到 structured_logs/ 目录")
	fmt.Println()

	fmt.Println("=== 文件日志示例完成 ===")
}
