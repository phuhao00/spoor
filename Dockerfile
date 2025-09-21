# 多阶段构建
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o spoor ./examples/complete

# 运行阶段
FROM alpine:latest

# 安装ca证书
RUN apk --no-cache add ca-certificates tzdata

# 创建非root用户
RUN adduser -D -s /bin/sh spoor

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/spoor .

# 创建日志目录
RUN mkdir -p logs && chown spoor:spoor logs

# 切换到非root用户
USER spoor

# 暴露端口（如果需要）
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=release
ENV TZ=UTC

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ./spoor --health-check || exit 1

# 运行应用
CMD ["./spoor"]
