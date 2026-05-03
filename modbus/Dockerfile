# 构建阶段
FROM golang:1.26 AS builder

WORKDIR /app

# 复制整个项目
COPY . .

# 设置 Go 代理
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

# 进入 modbus 目录执行（因为 go.mod 在那里）
WORKDIR /app/modbus
RUN go mod download

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/sbcc ./main/run.go

# 运行阶段
FROM alpine:3.19

WORKDIR /app

# 安装 CA 证书和时区数据
RUN apk add --no-cache ca-certificates tzdata

# 从构建阶段复制可执行文件和静态资源
COPY --from=builder /app/sbcc /app/sbcc
COPY --from=builder /app/modbus/static /app/static

# 暴露端口
EXPOSE 9081

# 设置时区
ENV TZ=Asia/Shanghai

# 启动命令
CMD ["/app/sbcc"]