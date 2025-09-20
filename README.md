# spoor

一个简单易用的Go日志库，支持多种输出方式和结构化日志记录。

## ✨ 特性

- 🚀 **简单易用** - 简洁的API设计，易于集成和使用
- 📁 **多种输出方式** - 支持文件、控制台、Elasticsearch、ClickHouse、Logbus等
- 🏗️ **结构化日志** - 支持字段和上下文的结构化日志记录
- ⚡ **高性能** - 异步写入，批量处理，性能优化
- 🔧 **灵活配置** - 支持日志级别、格式、轮转等配置
- 🛡️ **线程安全** - 支持并发安全的日志记录

## 💡 安装

```bash
go get github.com/phuhao00/spoor
```

## 🚀 快速开始

### 基本使用

```go
package main

import (
    "log"
    "os"
    "github.com/phuhao00/spoor"
)

func main() {
    // 创建控制台日志记录器
    logger := spoor.NewSpoor(
        spoor.DEBUG, 
        "", 
        log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, 
        spoor.WithConsoleWriter(os.Stdout),
    )
    
    logger.Debug("这是一条调试消息")
    logger.Info("这是一条信息消息")
    logger.Warn("这是一条警告消息")
    logger.Error("这是一条错误消息")
    
    // 格式化消息
    logger.DebugF("用户 %s 登录成功", "张三")
    logger.InfoF("处理了 %d 个请求", 100)
}
```

## 📁 输出方式

### 1. 文件输出 (FileWriter)

```go
// 创建文件写入器
fileWriter := spoor.NewFileWriter("logs", 0, 0, 0)
logger := spoor.NewSpoor(
    spoor.DEBUG, 
    "", 
    log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, 
    spoor.WithFileWriter(fileWriter),
)

logger.Info("这条消息将写入到文件")
```

### 2. 控制台输出 (ConsoleWriter)

```go
// 控制台输出
logger := spoor.NewSpoor(
    spoor.DEBUG, 
    "", 
    log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, 
    spoor.WithConsoleWriter(os.Stdout),
)

logger.Info("这条消息将输出到控制台")
```

### 3. Elasticsearch输出 (ElasticWriter)

```go
// Elasticsearch配置
config := spoor.ElasticConfig{
    URL:       "http://localhost:9200",
    Index:     "app-logs",
    BatchSize: 100,
    FlushTime: 5 * time.Second,
}

logger := spoor.NewSpoor(
    spoor.DEBUG, 
    "", 
    log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, 
    spoor.WithElasticWriter(config),
)

logger.Info("这条消息将发送到Elasticsearch")
```

### 4. ClickHouse输出 (ClickHouseWriter)

```go
// 创建ClickHouse日志记录器
logger, err := spoor.NewClickHouse("tcp://localhost:9000?database=logs", "app_logs", spoor.LevelInfo)
if err != nil {
    log.Fatal(err)
}

logger.Info("这条消息将发送到ClickHouse")
logger.Info("支持高性能的日志存储和查询")
```

### 5. Logbus输出 (LogbusWriter)

```go
// Logbus配置
config := spoor.LogbusConfig{
    URL:       "https://api.logbus.com/logs",
    Token:     "your-token",
    BatchSize: 100,
    FlushTime: 5 * time.Second,
}

logger := spoor.NewSpoor(
    spoor.DEBUG, 
    "", 
    log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, 
    spoor.WithLogbusWriter(config),
)

logger.Info("这条消息将发送到Logbus")
```

## 🏗️ 结构化日志

```go
// 创建结构化日志记录器
logger := spoor.NewSpoor(
    spoor.DEBUG, 
    "", 
    log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, 
    spoor.WithConsoleWriter(os.Stdout),
)

structuredLogger := spoor.NewStructuredLogger(logger)

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
```

## 🔧 日志级别

```go
const (
    DEBUG = Level(1)  // 调试级别
    INFO  = Level(2)  // 信息级别
    WARN  = Level(3)  // 警告级别
    ERROR = Level(4)  // 错误级别
    FATAL = Level(5)  // 致命级别
)
```

## 📝 全局日志记录器

```go
package main

import (
    "log"
    "github.com/phuhao00/spoor/logger"
)

func main() {
    // 设置全局日志记录器
    setting := &logger.LoggingSetting{
        Dir:    "logs",
        Level:  int(spoor.DEBUG),
        Prefix: "",
    }
    logger.SetLogging(setting)
    
    // 使用全局日志记录器
    logger.Debug("全局调试消息")
    logger.Info("全局信息消息")
    logger.Warn("全局警告消息")
    logger.Error("全局错误消息")
}
```

## 🧪 测试

```bash
go test -v
```

## 📄 许可证

MIT License