# 使用 Golang 官方镜像作为构建阶段的基础镜像
FROM crpi-9h0ljiysae0hwd8o.cn-shenzhen.personal.cr.aliyuncs.com/convoo/golang_1.21:1.21 AS builder

COPY . /src
WORKDIR /src

# 下载 Go 依赖并构建项目
#RUN go mod tidy
RUN GOPROXY=https://goproxy.cn make build

# 使用轻量级的 Alpine 镜像作为运行时镜像
FROM crpi-9h0ljiysae0hwd8o.cn-shenzhen.personal.cr.aliyuncs.com/convoo/alpine_3.16:v3.1.6

# 安装必要的依赖（如果有，例如：glibc）
#RUN apk update && apk add --no-cache libc6-compat

# 设置工作目录为 /app
COPY --from=builder /src/bin /app

WORKDIR /app

# 暴露容器的端口
EXPOSE 8000

VOLUME /data/conf

# 运行 Kratos 应用
CMD ["./server", "-conf", "/data/conf"]