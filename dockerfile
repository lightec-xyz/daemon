# 使用官方的golang镜像作为基础镜像
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 将go.mod和go.sum文件复制到工作目录
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 编译Go应用程序
RUN go build -o main .

# 使用一个更小的基础镜像
FROM alpine:latest

# 安装必要的依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从builder阶段复制编译后的二进制文件到当前目录
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行Go应用程序
CMD ["./main"]
