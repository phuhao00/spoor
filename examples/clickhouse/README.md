# ClickHouse 日志示例

这个示例展示了如何使用 spoor 库将日志发送到 ClickHouse。

## 前置要求

需要运行 ClickHouse 实例：

```bash
# 使用 Docker 运行 ClickHouse
docker run -d --name clickhouse -p 9000:9000 -p 8123:8123 clickhouse/clickhouse-server

# 或者使用 Docker Compose
docker-compose up -d clickhouse
```

## 运行示例

```bash
go run main.go
```

## 示例内容

1. **基本 ClickHouse 日志** - 发送简单日志到 ClickHouse
2. **结构化日志** - 使用字段记录结构化日志
3. **错误日志** - 记录错误和异常信息
4. **性能监控日志** - 记录性能指标和监控数据
5. **批量日志写入** - 批量发送日志提高性能
6. **不同日志级别** - 不同级别的日志记录
7. **时间序列日志** - 记录时间序列数据
8. **业务事件日志** - 记录业务事件和操作

## 特性

- 支持列式存储，查询性能优异
- 支持批量写入提高性能
- 支持时间序列数据
- 支持结构化日志记录
- 自动创建表结构
- 支持压缩存储

## 数据库表结构

ClickHouse Writer 会自动创建以下表结构：

```sql
CREATE TABLE IF NOT EXISTS app_logs (
    timestamp DateTime,
    level String,
    message String,
    fields String,  -- JSON格式的字段
    caller String
) ENGINE = MergeTree()
ORDER BY timestamp
```

## 配置选项

```go
config := spoor.ClickHouseWriterConfig{
    DSN:         "tcp://localhost:9000?database=logs",
    TableName:   "app_logs",
    BatchSize:   100,              // 批量大小
    FlushTime:   5,                // 刷新间隔（秒）
    HTTPTimeout: 30,               // HTTP 超时（秒）
}
```

## 查看日志

可以通过 ClickHouse 客户端查看日志：

```bash
# 连接 ClickHouse
clickhouse-client

# 查看表结构
DESCRIBE app_logs;

# 查询日志
SELECT * FROM app_logs ORDER BY timestamp DESC LIMIT 10;

# 按级别查询
SELECT * FROM app_logs WHERE level = 'ERROR' ORDER BY timestamp DESC;

# 按时间范围查询
SELECT * FROM app_logs WHERE timestamp >= now() - INTERVAL 1 HOUR;

# 统计日志级别
SELECT level, count() FROM app_logs GROUP BY level;
```

## 性能优势

- **列式存储**: 只读取需要的列，提高查询性能
- **压缩**: 自动压缩存储，节省空间
- **批量插入**: 支持批量插入，提高写入性能
- **索引**: 基于时间戳的索引，快速时间范围查询
