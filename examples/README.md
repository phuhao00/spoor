# Spoor 示例集合

这个文件夹包含了 spoor 日志库的各种使用示例，每个示例都独立运行，展示了不同的功能和用法。

## 示例列表

### 基础示例

- **[basic/](basic/)** - 基本使用示例，展示最简单的日志记录方法
- **[console/](console/)** - 控制台日志示例，展示控制台输出的各种用法
- **[file/](file/)** - 文件日志示例，展示文件输出和日志轮转

### 高级示例

- **[advanced/](advanced/)** - 高级用法示例，展示复杂场景下的日志记录
- **[refactored/](refactored/)** - 重构后的 API 使用示例
- **[config/](config/)** - 配置文件示例，展示如何通过配置文件管理日志

### 输出目标示例

- **[elasticsearch/](elasticsearch/)** - Elasticsearch 输出示例
- **[clickhouse/](clickhouse/)** - ClickHouse 输出示例

## 快速开始

每个示例都可以独立运行：

```bash
# 运行控制台日志示例
cd console
go run main.go

# 运行文件日志示例
cd file
go run main.go

# 运行 Elasticsearch 示例（需要先启动 ES）
cd elasticsearch
go run main.go
```

## 示例特性

### 控制台示例
- 多种日志级别
- 格式化消息
- 结构化日志
- 不同输出流 (stdout/stderr)
- 多种格式化器

### 文件示例
- 自动日志轮转
- 并发安全写入
- 批量写入优化
- 多种格式化器
- 自动目录创建

### Elasticsearch 示例
- 批量写入优化
- 重试机制
- 健康检查
- 结构化日志
- 性能监控

### ClickHouse 示例
- 列式存储优势
- 时间序列数据
- 批量插入
- 自动表创建
- 高性能查询

## 前置要求

某些示例需要外部服务：

- **Elasticsearch 示例**: 需要运行 Elasticsearch 实例
- **ClickHouse 示例**: 需要运行 ClickHouse 实例

## 学习建议

1. 从 `basic/` 开始，了解基本用法
2. 查看 `console/` 和 `file/` 了解不同输出方式
3. 学习 `advanced/` 和 `refactored/` 了解高级特性
4. 尝试 `elasticsearch/` 和 `clickhouse/` 了解外部存储
5. 使用 `config/` 了解配置管理

## 贡献

欢迎添加新的示例！每个示例应该：
- 有独立的 `main.go` 文件
- 包含详细的 `README.md` 说明
- 展示特定的功能或用法
- 包含注释和说明
