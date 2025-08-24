# Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/clin211/miniblog-v3.git.

# syntax=docker/dockerfile:1

# ============ 1) 安装基础工具（一次缓存，多次复用） ============
FROM golang:1.24.2-alpine AS base-tools
WORKDIR /app

# 系统工具与证书、时区
RUN apk add --no-cache git curl bash ca-certificates tzdata
ENV TZ=Asia/Shanghai

# 可选：加速国内 Go 依赖下载（按需调整或删除）
ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.org

# 固定 Air 版本，避免每次 latest 造成缓存失效
ENV AIR_VERSION=v1.62.0
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/air-verse/air@${AIR_VERSION}

# ============ 2) 仅下载依赖（层缓存 + BuildKit 缓存） ============
FROM base-tools AS deps
WORKDIR /app

# 仅复制 go.mod / go.sum，最大化层缓存命中率
COPY go.mod go.sum ./

# 使用 BuildKit 缓存挂载，显著提升重复构建速度
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

# ============ 3) 开发基础镜像（包含源码 + Air 已安装） ============
FROM deps AS dev-base
WORKDIR /app

# 创建必要目录
RUN mkdir -p /app/tmp /app/logs

# 复制源码（尽量通过 .dockerignore 减少上下文）
COPY . .

# ============ 4) API 服务镜像 ============
FROM dev-base AS api-service
LABEL service.type="api" \
      service.name="user-api"
EXPOSE 8888

# 这里如果有健康检查的 HTTP 端点可启用（示例）：
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#   CMD curl -sf http://localhost:8888/health || exit 1

# 启动 Air（使用已有的 user-api.toml）
CMD ["air", "-c", "/app/deploy/dev/air/user-api.toml"]

# ============ 5) RPC 服务镜像 ============
FROM dev-base AS rpc-service
LABEL service.type="rpc" \
      service.name="user-rpc"
EXPOSE 8889

# 如果没有 HTTP 端点，简单检查端口（可按需替换/删除）
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#   CMD nc -z localhost 8889 || exit 1

CMD ["air", "-c", "/app/deploy/dev/air/user-rpc.toml"]
