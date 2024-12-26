# 使用 Golang 官方镜像作为构建阶段的基础镜像
FROM crpi-9h0ljiysae0hwd8o.cn-shenzhen.personal.cr.aliyuncs.com/convoo/golang_1.20:v1.2.0 AS builder

# 设置工作目录为 /app
WORKDIR /app

# 将本地代码复制到容器中的 /app 目录
COPY . .

# 下载 Go 依赖并构建项目
#RUN go mod tidy
RUN go build -o kratos-server cmd/convoo-accounts/main.go cmd/convoo-accounts/wire_gen.go

# 使用轻量级的 Alpine 镜像作为运行时镜像
FROM crpi-9h0ljiysae0hwd8o.cn-shenzhen.personal.cr.aliyuncs.com/convoo/alpine_3.16:v3.1.6

# 安装必要的依赖（如果有，例如：glibc）
#RUN apk update && apk add --no-cache libc6-compat

# 设置工作目录为 /app
WORKDIR /app

# 从构建镜像中复制编译好的可执行文件
COPY --from=builder /app/kratos-server .

# 暴露容器的端口
EXPOSE 8000

# 运行 Kratos 应用
CMD ["./kratos-server"]