# 使用 Golang 官方镜像作为构建阶段的基础镜像
FROM crpi-9h0ljiysae0hwd8o.cn-shenzhen.personal.cr.aliyuncs.com/convoo/golang_1.21:1.21 AS builder

COPY . /src
WORKDIR /src

# 下载 Go 依赖并构建项目
#RUN go mod tidy
RUN GOPROXY=https://goproxy.cn make build

FROM crpi-9h0ljiysae0hwd8o.cn-shenzhen.personal.cr.aliyuncs.com/convoo/debian-stable-slim:1.0.0

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

# 设置工作目录为 /app
COPY --from=builder /src/bin /app

COPY ./configs /data/configs
RUN ["chmod", "-R", "777", "/data"]

WORKDIR /app

# 暴露容器的端口
EXPOSE 8000
EXPOSE 9000

VOLUME /data/conf

# 运行 Kratos 应用
CMD ["./convoo-accounts", "-conf", "/data/configs"]