# Spoor 版本管理说明

## 版本策略

本项目采用语义版本控制（Semantic Versioning），并遵循Go模块的版本管理规范。

## 版本系列

### v1.x 系列（稳定版）

- **模块路径**: `github.com/phuhao00/spoor`
- **最新版本**: v1.0.8
- **特点**: 稳定可靠，经过充分测试
- **推荐用途**: 生产环境

### v2.x 系列（最新版）

- **模块路径**: `github.com/phuhao00/spoor/v2`
- **最新版本**: v2.0.1
- **特点**: 功能更丰富，代码结构更清晰
- **推荐用途**: 新项目

## 版本历史

### v1.x 系列

| 版本 | 发布日期 | 主要特性 |
|------|----------|----------|
| v1.0.8 | 2024-01-15 | 添加日志实例管理 |
| v1.0.6 | 2024-01-10 | 逻辑优化 |
| v1.0.5 | 2024-01-08 | 逻辑优化 |
| v1.0.4 | 2024-01-05 | 逻辑优化 |
| v1.0.2 | 2024-01-01 | 修复callerSkip问题 |
| v1.0.1 | 2023-12-28 | 初始版本 |

### v2.x 系列

| 版本 | 发布日期 | 主要特性 |
|------|----------|----------|
| v2.0.1 | 2024-01-20 | 修复模块路径，完善功能 |
| v2.0.0 | 2024-01-18 | 重构代码结构，完善Elasticsearch Writer |

## 安装命令

### 安装特定版本

```bash
# 安装 v1.0.8
go get github.com/phuhao00/spoor@v1.0.8

# 安装 v2.0.1
go get github.com/phuhao00/spoor/v2@v2.0.1
```

### 安装最新版本

```bash
# 安装 v1.x 最新版本
go get github.com/phuhao00/spoor@latest

# 安装 v2.x 最新版本
go get github.com/phuhao00/spoor/v2@latest
```

## 版本检查

### 查看可用版本

```bash
# 查看 v1.x 版本
go list -m -versions github.com/phuhao00/spoor

# 查看 v2.x 版本
go list -m -versions github.com/phuhao00/spoor/v2
```

### 查看当前使用的版本

```bash
# 查看 v1.x 版本
go list -m github.com/phuhao00/spoor

# 查看 v2.x 版本
go list -m github.com/phuhao00/spoor/v2
```

## 版本兼容性

### Go 版本要求

- **最低要求**: Go 1.18+
- **推荐版本**: Go 1.24+

### 模块兼容性

- v1.x 和 v2.x 使用不同的模块路径
- 不能在同一项目中混用两个版本
- 迁移时需要更新导入路径

## 发布流程

### 1. 创建新版本

```bash
# 创建新标签
git tag v2.0.2
git push origin v2.0.2
```

### 2. 验证发布

```bash
# 验证标签
git show v2.0.2

# 验证模块
go list -m github.com/phuhao00/spoor/v2@v2.0.2
```

### 3. 更新文档

- 更新 README.md
- 更新 CHANGELOG.md
- 更新版本徽章

## 故障排除

### 网络问题

如果遇到网络连接问题：

```bash
# Windows PowerShell
$env:GOPROXY="direct"
$env:GOSUMDB="off"

# Linux/macOS
export GOPROXY="direct"
export GOSUMDB="off"
```

### 版本冲突

如果遇到版本冲突：

1. 清理模块缓存：
   ```bash
   go clean -modcache
   ```

2. 重新下载：
   ```bash
   go mod download
   ```

3. 更新依赖：
   ```bash
   go mod tidy
   ```

## 最佳实践

### 1. 版本选择

- **生产环境**: 使用 v1.0.8
- **新项目**: 使用 v2.0.1
- **学习测试**: 任选其一

### 2. 依赖管理

- 在 go.mod 中明确指定版本
- 定期更新到最新稳定版本
- 避免使用 `@latest` 在生产环境中

### 3. 迁移策略

- 逐步迁移，先测试后部署
- 保持版本一致性
- 记录迁移过程

## 支持

如果您在版本管理方面遇到问题，请：

1. 查看本文档
2. 查看 [常见问题解答](README.md#常见问题)
3. 提交 [Issue](https://github.com/phuhao00/spoor/issues)
