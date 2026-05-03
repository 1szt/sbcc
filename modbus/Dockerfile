# 构建阶段
FROM golang:1.26-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制源码
COPY . .

# 下载依赖
RUN go mod download

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o sbcc ./modbus/main/run.go

# 运行阶段
FROM alpine:3.19

WORKDIR /app

# 安装 CA 证书（用于 HTTPS 请求）
RUN apk add --no-cache ca-certificates tzdata

# 从构建阶段复制可执行文件
COPY --from=builder /app/sbcc .
COPY --from=builder /app/modbus/static ./static
COPY --from=builder /app/modbus/data ./data

# 创建必要目录
RUN mkdir -p /app/data

# 暴露端口
EXPOSE 9081

# 设置时区
ENV TZ=Asia/Shanghai

# 启动命令
CMD ["./sbcc"]