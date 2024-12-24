# 使用国内的 Golang 镜像作为基础镜像
FROM registry.cn-hangzhou.aliyuncs.com/aliyun/golang:1.20 AS builder

# 设置工作目录
WORKDIR /app

# 将项目代码复制到容器中
COPY . .

# 下载依赖并构建项目
RUN go mod tidy
RUN go build -o kratos-server main.go

# 使用国内的 Alpine 镜像
FROM registry.cn-hangzhou.aliyuncs.com/aliyun/alpine:3.16
WORKDIR /app

# 从构建镜像中复制可执行文件
COPY --from=builder /app/kratos-server .

# 暴露端口
EXPOSE 8000

# 运行服务
CMD ["./kratos-server"]