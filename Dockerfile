# FROM 指定了基础镜像，这里使用了一个定制的 Golang 1.21 镜像。AS builder 给这个阶段命名为 builder，这样后续可以通过 --from=builder 来引用它
FROM crpi-9h0ljiysae0hwd8o.cn-shenzhen.personal.cr.aliyuncs.com/convoo/golang_1.21:1.21 AS builder

# 将当前目录（宿主机上的目录）下的所有文件复制到镜像内的 /src 目录
COPY . /src

# 设置工作目录为 /src。这意味着后续的命令都将在这个目录下执行
WORKDIR /src

# 使用 RUN 指令执行命令。在这里，它设置了 Go 的代理服务器为 https://goproxy.cn，然后执行 make build，用来构建 Go 项目（假设你有一个 Makefile 文件，make build 会触发构建流程）
RUN GOPROXY=https://goproxy.cn GOARCH=amd64 GOOS=linux make build

# 这行指定了第二个阶段的基础镜像，这里是一个定制的 Debian 镜像。与第一个镜像不同，这个镜像通常用于运行应用程序，而不是构建应用程序。
FROM crpi-9h0ljiysae0hwd8o.cn-shenzhen.personal.cr.aliyuncs.com/convoo/debian-stable-slim:1.0.0

# 更新 apt 的包列表，并安装 ca-certificates 和 netbase 包（这些是最基本的依赖）。--no-install-recommends 选项表示只安装必需的包，不安装推荐包。
#安装后，清理缓存以减小镜像大小（rm -rf /var/lib/apt/lists/ 以及执行 apt-get autoremove 和 apt-get autoclean）
RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

# 从 builder 阶段复制 /src/bin 目录下的文件到 /app 目录，这样就把构建好的 Go 程序从构建镜像复制到运行镜像中
COPY --from=builder /src/bin /app

# 将宿主机当前目录下的 configs 目录复制到镜像中的 /data/configs 目录，这通常是应用的配置文件
COPY ./configs /data/configs

# 给予 /data 目录及其下所有文件 777 的权限，这表示所有用户都能读取、写入和执行这个目录的内容
RUN ["chmod", "-R", "777", "/data"]

# 设置工作目录为 /app，这意味着后续的命令（包括启动容器时的命令）将在这个目录下执行
WORKDIR /app

# 暴露容器的端口
EXPOSE 8000
EXPOSE 9000

# 创建一个数据卷挂载点 /data/conf。这表示容器运行时可以将宿主机上的某个目录或 Docker 卷挂载到容器的这个位置，通常用来持久化存储数据或配置
VOLUME /data/conf

# 设置容器启动时执行的命令。这里是运行 /app/accounts 程序并传递 -conf /data/configs 参数，意味着应用会加载 /data/configs 下的配置文件
CMD ["./convoo-accounts", "-conf", "/data/configs"]