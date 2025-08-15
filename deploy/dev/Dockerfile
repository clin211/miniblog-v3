# 基础镜像
FROM golang:1.24.2-alpine AS base

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git curl bash

# 安装 Air 热重载工具
RUN go install github.com/air-verse/air@latest

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# API服务镜像
FROM base AS api-service
LABEL service.type="api"
EXPOSE 8888

# RPC服务镜像
FROM base AS rpc-service
LABEL service.type="rpc"
EXPOSE 8889

# 创建tmp目录用于存放编译后的二进制文件
RUN mkdir -p /app/tmp

# 默认命令
CMD ["air", "-c", "/app/deploy/dev/air/user-api.toml"]
