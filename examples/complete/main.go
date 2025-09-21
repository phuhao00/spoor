package main

import (
	"fmt"
	"os"
	"time"

	"github.com/phuhao00/spoor"
)

func main() {
	fmt.Println("🎯 Spoor Complete Example")
	fmt.Println("=========================")

	// 1. 超简单使用
	fmt.Println("\n1. 超简单使用")
	simpleLogger := spoor.Quick()
	simpleLogger.Info("Hello, Spoor!")
	simpleLogger.Close()

	// 2. 高性能异步日志
	fmt.Println("\n2. 高性能异步日志")
	asyncLogger := spoor.QuickAsync()
	
	// 记录开始时间
	start := time.Now()
	
	// 大量日志写入
	for i := 0; i < 10000; i++ {
		asyncLogger.Infof("Async message %d", i)
	}
	
	// 等待所有消息处理完成
	asyncLogger.Sync()
	duration := time.Since(start)
	
	fmt.Printf("写入10000条消息耗时: %v\n", duration)
	fmt.Printf("吞吐量: %.0f 消息/秒\n", float64(10000)/duration.Seconds())
	
	asyncLogger.Close()

	// 3. 文件日志
	fmt.Println("\n3. 文件日志")
	fileLogger, err := spoor.QuickFile("logs/complete.log")
	if err != nil {
		fmt.Printf("创建文件日志器失败: %v\n", err)
	} else {
		fileLogger.Info("这条消息将写入到文件")
		fileLogger.WithField("user_id", 123).Info("用户操作")
		fileLogger.Close()
	}

	// 4. JSON格式日志
	fmt.Println("\n4. JSON格式日志")
	jsonLogger := spoor.QuickJSON()
	jsonLogger.WithFields(map[string]interface{}{
		"user_id":    123,
		"action":     "login",
		"ip":         "192.168.1.1",
		"user_agent": "Mozilla/5.0",
		"timestamp":  time.Now().Unix(),
	}).Info("用户登录")
	jsonLogger.Close()

	// 5. 高级功能 - 采样
	fmt.Println("\n5. 高级功能 - 采样")
	sampler := spoor.NewRateSampler(0.1) // 10%采样
	advancedConfig := spoor.AdvancedConfig{
		Sampler: sampler,
		Metrics: true,
	}
	
	writer := spoor.NewConsoleWriter(spoor.ConsoleWriterConfig{Output: os.Stdout})
	advancedLogger := spoor.NewAdvancedLogger(writer, spoor.LevelInfo, advancedConfig)
	
	// 发送100条消息，只有约10条会被记录
	for i := 0; i < 100; i++ {
		advancedLogger.Infof("Sampled message %d", i)
	}
	
	// 获取指标
	metrics := advancedLogger.GetMetrics()
	fmt.Printf("采样指标: %+v\n", metrics)
	advancedLogger.Close()

	// 6. 高级功能 - 过滤
	fmt.Println("\n6. 高级功能 - 过滤")
	levelFilter := spoor.NewLevelFilter(spoor.LevelWarn)
	filterConfig := spoor.AdvancedConfig{
		Filter:  levelFilter,
		Metrics: true,
	}
	
	filterLogger := spoor.NewAdvancedLogger(writer, spoor.LevelInfo, filterConfig)
	
	// 这些消息会被过滤掉
	filterLogger.Debug("这条调试消息会被过滤")
	filterLogger.Info("这条信息消息会被过滤")
	
	// 这些消息会被记录
	filterLogger.Warn("这条警告消息会被记录")
	filterLogger.Error("这条错误消息会被记录")
	
	filterLogger.Close()

	// 7. 性能监控
	fmt.Println("\n7. 性能监控")
	monitor := spoor.NewPerformanceMonitor()
	
	// 创建带监控的日志器
	monitoredLogger := spoor.QuickAsync()
	
	// 记录一些日志
	for i := 0; i < 1000; i++ {
		start := time.Now()
		monitoredLogger.Infof("Monitored message %d", i)
		monitor.RecordLog()
		monitor.RecordLatency(time.Since(start))
	}
	
	monitoredLogger.Sync()
	monitoredLogger.Close()
	
	// 打印性能统计
	monitor.PrintStats()
	monitor.Close()

	// 8. 批量写入
	fmt.Println("\n8. 批量写入")
	batchConfig := spoor.DefaultBatchConfig()
	batchConfig.BatchSize = 100
	batchConfig.FlushInterval = 50 * time.Millisecond
	
	batchWriter := spoor.NewBatchWriter(writer, batchConfig)
	batchLogger := spoor.NewCoreLogger(batchWriter, spoor.LevelInfo)
	
	start = time.Now()
	for i := 0; i < 1000; i++ {
		batchLogger.Infof("Batch message %d", i)
	}
	batchLogger.Sync()
	duration = time.Since(start)
	
	fmt.Printf("批量写入1000条消息耗时: %v\n", duration)
	fmt.Printf("吞吐量: %.0f 消息/秒\n", float64(1000)/duration.Seconds())
	
	// 获取批量写入指标
	batchMetrics := batchWriter.GetMetrics()
	fmt.Printf("批量写入指标: %+v\n", batchMetrics)
	
	batchLogger.Close()

	// 9. 内存池使用
	fmt.Println("\n9. 内存池使用")
	poolLogger := spoor.Quick()
	
	start = time.Now()
	for i := 0; i < 1000; i++ {
		// 使用内存池优化
		poolLogger.WithFields(map[string]interface{}{
			"request_id": fmt.Sprintf("req-%d", i),
			"user_id":    i,
			"duration":   time.Millisecond * time.Duration(i%100),
			"status":     200,
		}).Info("Pooled message")
	}
	poolLogger.Sync()
	duration = time.Since(start)
	
	fmt.Printf("内存池优化写入1000条消息耗时: %v\n", duration)
	poolLogger.Close()

	// 10. 配置管理
	fmt.Println("\n10. 配置管理")
	config := spoor.DefaultConfig()
	
	// 保存配置到文件
	if err := spoor.SaveConfig(config, "spoor-config.json"); err != nil {
		fmt.Printf("保存配置失败: %v\n", err)
	} else {
		fmt.Println("配置已保存到 spoor-config.json")
	}
	
	// 从配置创建日志器
	if configLogger, err := spoor.CreateLoggerFromConfig(&config.Loggers["default"]); err != nil {
		fmt.Printf("从配置创建日志器失败: %v\n", err)
	} else {
		configLogger.Info("从配置创建的日志器")
		configLogger.Close()
	}

	// 11. 全局日志器
	fmt.Println("\n11. 全局日志器")
	// 使用全局日志器
	spoor.Info("使用全局日志器")
	spoor.WithField("component", "main").Info("主组件日志")
	
	// 设置自定义全局日志器
	customLogger := spoor.QuickJSON()
	spoor.SetGlobalSimple(customLogger)
	spoor.Info("现在使用自定义全局日志器")
	customLogger.Close()

	// 12. 性能对比
	fmt.Println("\n12. 性能对比")
	comparePerformance()

	fmt.Println("\n✅ 完整示例运行完成!")
}

func comparePerformance() {
	iterations := 10000
	
	// 测试1: 简单日志器
	start := time.Now()
	simpleLogger := spoor.Quick()
	for i := 0; i < iterations; i++ {
		simpleLogger.Info("Simple message")
	}
	simpleLogger.Sync()
	simpleDuration := time.Since(start)
	simpleLogger.Close()

	// 测试2: 异步日志器
	start = time.Now()
	asyncLogger := spoor.QuickAsync()
	for i := 0; i < iterations; i++ {
		asyncLogger.Info("Async message")
	}
	asyncLogger.Sync()
	asyncDuration := time.Since(start)
	asyncLogger.Close()

	// 测试3: 批量日志器
	start = time.Now()
	batchConfig := spoor.DefaultBatchConfig()
	batchConfig.BatchSize = 1000
	writer := spoor.NewConsoleWriter(spoor.ConsoleWriterConfig{Output: os.Stdout})
	batchWriter := spoor.NewBatchWriter(writer, batchConfig)
	batchLogger := spoor.NewCoreLogger(batchWriter, spoor.LevelInfo)
	for i := 0; i < iterations; i++ {
		batchLogger.Info("Batch message")
	}
	batchLogger.Sync()
	batchDuration := time.Since(start)
	batchLogger.Close()

	// 计算吞吐量
	simpleThroughput := float64(iterations) / simpleDuration.Seconds()
	asyncThroughput := float64(iterations) / asyncDuration.Seconds()
	batchThroughput := float64(iterations) / batchDuration.Seconds()

	fmt.Printf("简单日志器:  %v (%.0f 消息/秒)\n", simpleDuration, simpleThroughput)
	fmt.Printf("异步日志器:  %v (%.0f 消息/秒)\n", asyncDuration, asyncThroughput)
	fmt.Printf("批量日志器:  %v (%.0f 消息/秒)\n", batchDuration, batchThroughput)
	
	// 计算性能提升
	asyncImprovement := (asyncThroughput - simpleThroughput) / simpleThroughput * 100
	batchImprovement := (batchThroughput - simpleThroughput) / simpleThroughput * 100
	
	fmt.Printf("异步日志器性能提升: %.1f%%\n", asyncImprovement)
	fmt.Printf("批量日志器性能提升: %.1f%%\n", batchImprovement)
}
