#!/usr/bin/env bash

# Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/clin211/miniblog-v3.git.

set -euo pipefail

# ----------------------------------------------
# 基于 deploy/dev/docker-compose.env.yml 的 MySQL 初始化脚本
# - 通过 Docker 容器 miniblog-mysql 注入 SQL：deploy/sql/user.sql
# - 自动等待 MySQL 服务可用
# ----------------------------------------------

# MySQL 连接配置（来自 docker-compose.env.yml）
MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"
MYSQL_ROOT_USER="root"
MYSQL_ROOT_PASSWORD="j478EaZGDNPUbnXb"
MYSQL_CONTAINER_NAME="miniblog-mysql"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
SQL_FILE="$REPO_ROOT/sql/user.sql"

abort() {
  echo "[ERROR] $*" >&2
  exit 1
}

info() {
  echo "[INFO] $*"
}

if [[ ! -f "$SQL_FILE" ]]; then
  abort "未找到 SQL 文件：$SQL_FILE"
fi

# 检查 Docker 与容器状态
if ! command -v docker >/dev/null 2>&1; then
  abort "未检测到 docker，请先安装 Docker Desktop 并启动。"
fi

if ! docker ps --format '{{.Names}}' | grep -q "^${MYSQL_CONTAINER_NAME}$"; then
  abort "未发现运行中的容器 ${MYSQL_CONTAINER_NAME}。请先启动 dev 环境（例如：docker compose -f deploy/dev/docker-compose.env.yml up -d mysql）。"
fi

# 等待 MySQL 就绪
wait_for_mysql() {
  local retries=60
  local count=0
  info "等待 MySQL (${MYSQL_CONTAINER_NAME}) 就绪..."
  while (( count < retries )); do
    if docker exec "${MYSQL_CONTAINER_NAME}" sh -c "mysqladmin ping -u${MYSQL_ROOT_USER} -p${MYSQL_ROOT_PASSWORD} --silent" >/dev/null 2>&1; then
      info "MySQL 已就绪。"
      return 0
    fi
    count=$((count + 1))
    sleep 1
  done
  return 1
}

wait_for_mysql || abort "MySQL 在预期时间内未就绪。"

# 执行导入
info "开始导入 SQL：$SQL_FILE"
docker exec -i "${MYSQL_CONTAINER_NAME}" sh -c "mysql -u${MYSQL_ROOT_USER} -p${MYSQL_ROOT_PASSWORD}" < "$SQL_FILE" || abort "SQL 导入失败。"

info "SQL 初始化完成。"


