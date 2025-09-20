# Elasticsearch 日志示例

这个示例展示了如何使用 spoor 库将日志发送到 Elasticsearch。

## 前置要求

需要运行 Elasticsearch 实例：

```bash
# 使用 Docker 运行 Elasticsearch
docker run -d --name elasticsearch -p 9200:9200 -e "discovery.type=single-node" elasticsearch:7.15.0

# 或者使用 Docker Compose
docker-compose up -d elasticsearch
```

## 运行示例

```bash
go run main.go
```

## 示例内容

1. **基本 Elasticsearch 日志** - 发送简单日志到 Elasticsearch
2. **结构化日志** - 使用字段记录结构化日志
3. **错误日志** - 记录错误和异常信息
4. **性能监控日志** - 记录性能指标和监控数据
5. **批量日志写入** - 批量发送日志提高性能
6. **不同日志级别** - 不同级别的日志记录
7. **应用生命周期日志** - 记录应用启动和运行状态

## 特性

- 支持批量写入提高性能
- 支持重试机制
- 支持结构化日志记录
- 支持多种日志级别
- 自动错误处理
- 支持健康检查

## 配置选项

```go
config := spoor.ElasticWriterConfig{
    URL:           "http://localhost:9200",
    Index:         "app-logs",
    Username:      "elastic",        // 可选
    Password:      "password",       // 可选
    APIKey:        "api-key",        // 可选
    BatchSize:     100,              // 批量大小
    FlushInterval: 5,                // 刷新间隔（秒）
    HTTPTimeout:   30,               // HTTP 超时（秒）
    RetryCount:    3,                // 重试次数
    RetryDelay:    1,                // 重试延迟（秒）
}
```

## 查看日志

可以通过 Elasticsearch 的 REST API 查看日志：

```bash
# 查看索引列表
curl -X GET "localhost:9200/_cat/indices?v"

# 搜索日志
curl -X GET "localhost:9200/app-logs/_search?pretty"

# 查看特定服务的日志
curl -X GET "localhost:9200/app-logs/_search?q=service:user-service&pretty"
```
