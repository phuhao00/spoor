# Spoor 版本迁移指南

## 概述

本指南将帮助您从 Spoor v1.x 迁移到 v2.x 版本。

## 主要变化

### 1. 模块路径变化

**v1.x 版本：**
```go
import "github.com/phuhao00/spoor"
```

**v2.x 版本：**
```go
import "github.com/phuhao00/spoor/v2"
```

### 2. 安装命令变化

**v1.x 版本：**
```bash
go get github.com/phuhao00/spoor@v1.0.8
```

**v2.x 版本：**
```bash
go get github.com/phuhao00/spoor/v2@v2.0.1
```

## 迁移步骤

### 步骤 1：更新 go.mod 文件

在您的项目中，将：
```go
require github.com/phuhao00/spoor v1.0.8
```

更改为：
```go
require github.com/phuhao00/spoor/v2 v2.0.1
```

### 步骤 2：更新导入语句

将所有文件中的导入语句从：
```go
import "github.com/phuhao00/spoor"
```

更改为：
```go
import "github.com/phuhao00/spoor/v2"
```

### 步骤 3：运行 go mod tidy

```bash
go mod tidy
```

### 步骤 4：测试应用程序

确保所有功能正常工作：
```bash
go test ./...
go run main.go
```

## 兼容性说明

### API 兼容性

v2.x 版本与 v1.x 版本在 API 层面基本兼容，主要变化是模块路径。

### 功能增强

v2.x 版本包含以下增强：

1. **重构的代码结构** - 更清晰的模块组织
2. **改进的 Elasticsearch Writer** - 更好的性能和稳定性
3. **完善的示例** - 更丰富的使用示例
4. **更好的错误处理** - 更详细的错误信息

## 回滚方案

如果需要在迁移后回滚到 v1.x 版本：

### 步骤 1：更新 go.mod

```go
require github.com/phuhao00/spoor v1.0.8
```

### 步骤 2：更新导入语句

```go
import "github.com/phuhao00/spoor"
```

### 步骤 3：清理依赖

```bash
go mod tidy
```

## 常见问题

### Q: 迁移后出现编译错误怎么办？

A: 检查以下几点：
1. 确保所有导入语句都已更新
2. 运行 `go mod tidy` 清理依赖
3. 检查是否有遗漏的文件

### Q: 可以同时使用两个版本吗？

A: 不推荐。虽然技术上可行，但会增加复杂性。建议选择一个版本并保持一致。

### Q: 迁移需要多长时间？

A: 对于大多数项目，迁移只需要几分钟：
1. 更新导入语句（1-2分钟）
2. 更新 go.mod（30秒）
3. 测试（2-5分钟）

## 获取帮助

如果您在迁移过程中遇到问题，请：

1. 查看 [常见问题解答](README.md#常见问题)
2. 检查 [示例代码](examples/)
3. 提交 [Issue](https://github.com/phuhao00/spoor/issues)

## 总结

从 v1.x 迁移到 v2.x 是一个简单的过程，主要是更新模块路径。v2.x 版本提供了更好的性能和更丰富的功能，建议新项目直接使用 v2.x 版本。
