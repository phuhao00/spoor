package main

import (
	"fmt"
	"os"

	"github.com/phuhao00/spoor"
)

func main() {
	fmt.Println("=== 控制台日志示例 ===\n")

	// 示例1: 基本控制台日志
	fmt.Println("1. 基本控制台日志")
	fmt.Println("==================")
	logger := spoor.NewConsole(spoor.LevelDebug)
	logger.Debug("这是一条调试消息")
	logger.Info("这是一条信息消息")
	logger.Warn("这是一条警告消息")
	logger.Error("这是一条错误消息")
	fmt.Println()

	// 示例2: 格式化消息
	fmt.Println("2. 格式化消息")
	fmt.Println("=============")
	logger.Debugf("用户 %s 登录成功", "张三")
	logger.Infof("处理了 %d 个请求", 100)
	logger.Warnf("内存使用率: %.2f%%", 85.5)
	logger.Errorf("连接失败，重试第 %d 次", 3)
	fmt.Println()

	// 示例3: 结构化日志
	fmt.Println("3. 结构化日志")
	fmt.Println("=============")
	structuredLogger := logger.WithField("user_id", 123).WithField("action", "login")
	structuredLogger.Info("用户登录")

	structuredLogger = logger.WithFields(map[string]interface{}{
		"request_id": "req-123",
		"duration":   "150ms",
		"status":     200,
		"ip":         "192.168.1.100",
	})
	structuredLogger.Info("请求完成")

	// 错误日志
	err := fmt.Errorf("数据库连接失败")
	structuredLogger = logger.WithError(err)
	structuredLogger.Error("数据库操作失败")
	fmt.Println()

	// 示例4: 不同日志级别
	fmt.Println("4. 日志级别控制")
	fmt.Println("===============")
	levelLogger := spoor.NewConsole(spoor.LevelInfo)
	levelLogger.Debug("这条调试消息不会显示")
	levelLogger.Info("这条信息消息会显示")
	levelLogger.Warn("这条警告消息会显示")
	levelLogger.Error("这条错误消息会显示")

	// 动态调整日志级别
	levelLogger.SetLevel(spoor.LevelError)
	levelLogger.Info("这条信息消息现在不会显示")
	levelLogger.Error("这条错误消息仍然会显示")
	fmt.Println()

	// 示例5: 输出到 stderr
	fmt.Println("5. 输出到 stderr")
	fmt.Println("================")
	stderrLogger := spoor.New(spoor.NewConsoleWriter(spoor.ConsoleWriterConfig{
		Output: os.Stderr,
	}), spoor.LevelInfo)
	stderrLogger.Error("这条错误消息输出到 stderr")
	fmt.Println()

	// 示例6: 不同格式化器
	fmt.Println("6. 不同格式化器")
	fmt.Println("===============")

	// JSON 格式化器
	jsonLogger := spoor.NewJSON(
		spoor.NewConsoleWriter(spoor.ConsoleWriterConfig{
			Output:    os.Stdout,
			Formatter: spoor.NewJSONFormatter(),
		}),
		spoor.LevelInfo,
	)
	jsonLogger.WithField("component", "auth").Info("JSON 格式的日志")

	// Logrus 格式化器
	logrusLogger := spoor.NewLogrus(
		spoor.NewConsoleWriter(spoor.ConsoleWriterConfig{
			Output:    os.Stdout,
			Formatter: spoor.NewLogrusFormatter(),
		}),
		spoor.LevelInfo,
	)
	logrusLogger.WithField("component", "auth").Info("Logrus 格式的日志")

	fmt.Println("\n=== 控制台日志示例完成 ===")
}
