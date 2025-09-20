# Spoor 版本使用指南

## 问题解决

之前您遇到的问题是Go模块版本管理的问题。现在已经完全解决了！

## 可用版本

### v1.x 系列（稳定版本）
- **最新版本**: v1.0.8
- **模块路径**: `github.com/phuhao00/spoor`
- **特点**: 经过充分测试，稳定可靠

### v2.x 系列（最新版本）
- **最新版本**: v2.0.1
- **模块路径**: `github.com/phuhao00/spoor/v2`
- **特点**: 包含最新功能，重构了代码结构

## 使用方法

### 1. 使用 v1.0.8（推荐用于生产环境）

```bash
# 下载 v1.0.8
go get github.com/phuhao00/spoor@v1.0.8

# 在代码中使用
import "github.com/phuhao00/spoor"
```

### 2. 使用 v2.0.1（推荐用于新项目）

```bash
# 下载 v2.0.1
go get github.com/phuhao00/spoor/v2@v2.0.1

# 在代码中使用
import "github.com/phuhao00/spoor/v2"
```

### 3. 解决网络问题

如果遇到网络连接问题，可以使用以下设置：

```bash
# 设置环境变量
$env:GOPROXY="direct"
$env:GOSUMDB="off"

# 然后下载
go get github.com/phuhao00/spoor@v1.0.8
# 或
go get github.com/phuhao00/spoor/v2@v2.0.1
```

## 版本选择建议

- **生产环境**: 使用 v1.0.8，稳定可靠
- **新项目**: 使用 v2.0.1，功能更丰富
- **学习测试**: 两个版本都可以，v2.0.1有更多示例

## 迁移指南

### 从 v1.x 迁移到 v2.x

1. 更新导入路径：
   ```go
   // 旧版本
   import "github.com/phuhao00/spoor"
   
   // 新版本
   import "github.com/phuhao00/spoor/v2"
   ```

2. 更新 go.mod：
   ```bash
   go get github.com/phuhao00/spoor/v2@v2.0.1
   ```

## 验证安装

```bash
# 验证 v1.0.8
go list -m github.com/phuhao00/spoor@v1.0.8

# 验证 v2.0.1
go list -m github.com/phuhao00/spoor/v2@v2.0.1
```

## 常见问题

### Q: 为什么下载的不是最新版本？
A: 因为Go模块系统遵循语义版本控制规则。v2.x版本需要使用不同的模块路径。

### Q: 如何选择使用哪个版本？
A: 生产环境推荐v1.0.8，新项目推荐v2.0.1。

### Q: 网络连接问题怎么办？
A: 使用 `GOPROXY="direct"` 和 `GOSUMDB="off"` 环境变量。

现在您可以正常使用最新版本了！
