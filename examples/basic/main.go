package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/phuhao00/spoor"
)

func main() {
	// 示例1: 控制台输出
	consoleLogger := spoor.NewSpoor(
		spoor.DEBUG,
		"",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile,
		spoor.WithConsoleWriter(os.Stdout),
	)

	consoleLogger.Debug("这是一条调试消息")
	consoleLogger.Info("这是一条信息消息")
	consoleLogger.Warn("这是一条警告消息")
	consoleLogger.Error("这是一条错误消息")

	// 格式化消息
	consoleLogger.DebugF("用户 %s 登录成功", "张三")
	consoleLogger.InfoF("处理了 %d 个请求", 100)

	// 示例2: 文件输出
	fileWriter := spoor.NewFileWriter("logs", 0, 0, 0)
	fileLogger := spoor.NewSpoor(
		spoor.DEBUG,
		"",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile,
		spoor.WithFileWriter(fileWriter),
	)

	fileLogger.Info("这条消息将写入到文件")
	fileLogger.InfoF("当前时间: %v", time.Now())

	// 示例3: 结构化日志
	structuredLogger := spoor.NewStructuredLogger(consoleLogger)

	// 添加字段
	structuredLogger.WithField("user_id", 123).WithField("action", "login").Info("用户登录")

	// 添加多个字段
	structuredLogger.WithFields(spoor.Fields{
		"request_id": "req-123",
		"duration":   time.Millisecond * 150,
		"status":     200,
	}).Info("请求完成")

	// 添加错误
	err := fmt.Errorf("数据库连接失败")
	structuredLogger.WithError(err).Error("数据库操作失败")

	// 示例4: 日志级别过滤
	infoLogger := spoor.NewSpoor(
		spoor.INFO, // 只记录INFO及以上级别
		"",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile,
		spoor.WithConsoleWriter(os.Stdout),
	)

	infoLogger.Debug("这条调试消息不会显示") // 级别太低
	infoLogger.Info("这条信息消息会显示")
	infoLogger.Warn("这条警告消息会显示")
	infoLogger.Error("这条错误消息会显示")
}
