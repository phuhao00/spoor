# spoor

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Version](https://img.shields.io/badge/version-v2.0.1-green.svg)](https://github.com/phuhao00/spoor)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

一个简单易用的Go日志库，支持多种输出方式和结构化日志记录。

> **版本说明**：本项目有两个版本系列
> - **v1.x**（稳定版）：`github.com/phuhao00/spoor` 
> - **v2.x**（最新版）：`github.com/phuhao00/spoor/v2`

## ✨ 特性

- 🚀 **简单易用** - 简洁的API设计，易于集成和使用
- 📁 **多种输出方式** - 支持文件、控制台、Elasticsearch、ClickHouse、Logbus等
- 🏗️ **结构化日志** - 支持字段和上下文的结构化日志记录
- ⚡ **高性能** - 异步写入，批量处理，性能优化
- 🔧 **灵活配置** - 支持日志级别、格式、轮转等配置
- 🛡️ **线程安全** - 支持并发安全的日志记录

## 💡 安装

### 版本说明

本项目有两个主要版本系列：

- **v1.x 系列**（稳定版）：`github.com/phuhao00/spoor`
- **v2.x 系列**（最新版）：`github.com/phuhao00/spoor/v2`

### 安装最新稳定版本 (v1.0.8)

```bash
go get github.com/phuhao00/spoor@v1.0.8
```

### 安装最新版本 (v2.0.1)

```bash
go get github.com/phuhao00/spoor/v2@v2.0.1
```

### 版本选择建议

- **生产环境**：推荐使用 v1.0.8（稳定可靠）
- **新项目**：推荐使用 v2.0.1（功能更丰富）
- **学习测试**：两个版本都可以使用

### 网络问题解决

如果遇到网络连接问题，可以使用以下设置：

```bash
# Windows PowerShell
$env:GOPROXY="direct"
$env:GOSUMDB="off"
go get github.com/phuhao00/spoor/v2@v2.0.1

# Linux/macOS
export GOPROXY="direct"
export GOSUMDB="off"
go get github.com/phuhao00/spoor/v2@v2.0.1
```

## 🚀 快速开始

### v1.x 版本使用

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

### v2.x 版本使用

```go
package main

import (
    "log"
    "os"
    "github.com/phuhao00/spoor/v2"
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

## ❓ 常见问题

### Q: 为什么 `go list -m -versions github.com/phuhao00/spoor` 只显示v1.x版本？

A: 这是因为Go模块系统的设计原理。v1.x和v2.x版本使用不同的模块路径：

- v1.x版本：`github.com/phuhao00/spoor`
- v2.x版本：`github.com/phuhao00/spoor/v2`

要查看v2.x版本，请使用：
```bash
go list -m -versions github.com/phuhao00/spoor/v2
```

### Q: 如何选择使用哪个版本？

A: 
- **生产环境**：推荐使用 v1.0.8（稳定可靠，经过充分测试）
- **新项目**：推荐使用 v2.0.1（功能更丰富，代码结构更清晰）
- **学习测试**：两个版本都可以使用

### Q: 如何从v1.x迁移到v2.x？

A: 只需要更改导入路径：

```go
// 旧版本 (v1.x)
import "github.com/phuhao00/spoor"

// 新版本 (v2.x)
import "github.com/phuhao00/spoor/v2"
```

### Q: 遇到网络连接问题怎么办？

A: 使用以下环境变量绕过代理：

```bash
# Windows PowerShell
$env:GOPROXY="direct"
$env:GOSUMDB="off"

# Linux/macOS
export GOPROXY="direct"
export GOSUMDB="off"
```

### Q: 如何验证安装的版本？

A: 使用以下命令验证：

```bash
# 验证 v1.x 版本
go list -m github.com/phuhao00/spoor@v1.0.8

# 验证 v2.x 版本
go list -m github.com/phuhao00/spoor/v2@v2.0.1
```

## 📄 许可证

MIT License